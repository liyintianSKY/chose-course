package natsclient

import (
	"chose-course/common/cmap"
	"chose-course/common/errmsg"
	"chose-course/common/utils"
	"chose-course/consts"
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"sync/atomic"
	"time"

	"github.com/nats-io/nats.go"

	"go.uber.org/zap"
)

type NatsMessage struct {
	FromServer consts.ModulesCode `json:"from_server"`
	Data       interface{}        `json:"data"`
}

type NatsClient struct {
	subs   cmap.ConcurrentMap[string, *nats.Subscription]
	name   string
	urls   string
	log    *zap.Logger
	conn   *nats.Conn
	closed int32
	f      func(*nats.Msg)
}

func NewNatsClient(name string, urls string, log *zap.Logger) *NatsClient {
	nc := &NatsClient{
		subs: cmap.New[*nats.Subscription](),
		urls: urls,
		log:  log,
		name: name,
	}
	c, err := nats.Connect(urls, nats.ReconnectWait(time.Millisecond*10), nats.MaxReconnects(math.MaxInt64),
		nats.PingInterval(time.Second*3), nats.MaxPingsOutstanding(2), nats.Timeout(time.Second),
		nats.DrainTimeout(time.Second*5), nats.Name(name),
		nats.DisconnectErrHandler(func(conn *nats.Conn, err error) {
			if atomic.LoadInt32(&nc.closed) == 0 {
				log.Error("nats disconnected", zap.Error(err), zap.String("urls", urls), zap.String("nats-server", conn.ConnectedAddr()))
			}
		}),
		nats.ReconnectHandler(func(conn *nats.Conn) {
			log.Warn("nats reconnected", zap.String("urls", urls), zap.String("nats-server", conn.ConnectedAddr()))
		}),
		nats.ClosedHandler(func(conn *nats.Conn) {
			log.Warn("nats closed", zap.String("urls", urls), zap.String("nats-server", conn.ConnectedAddr()))
		}),
	)
	if err != nil {
		panic(err)
	}
	nc.conn = c
	return nc
}

// Close 关闭nats
func (this_ *NatsClient) Close() {
	this_.subs.IterCb(func(key string, v *nats.Subscription) bool {
		if v.IsValid() {
			err := v.Drain()
			if err != nil {
				this_.log.Warn("Drain error", zap.String("subj", key), zap.Error(err))
			}
		}
		return true
	})
	this_.subs.IterCb(func(key string, v *nats.Subscription) bool {
		for v.IsValid() {
			time.Sleep(time.Millisecond * 10)
		}
		return true
	})
}

// Publish 推送数据
func (this_ *NatsClient) Publish(fromModule consts.ModulesCode, msgName string, msg interface{}) *errmsg.ErrMsg {
	message := &NatsMessage{
		FromServer: fromModule,
		Data:       msg,
	}
	msgBytes, err := json.Marshal(message)
	if err != nil {
		return errmsg.NewNormalErrorInfo("nats message marshal error", err.Error())
	}

	return errmsg.NewProtocolErrorInfo(this_.conn.Publish(msgName, msgBytes).Error())
}

func (this_ *NatsClient) Request(fromModule consts.ModulesCode, msgName string, req *NatsMessage) ([]byte, *errmsg.ErrMsg) {
	msgBytes, err := json.Marshal(req)
	if err != nil {
		return nil, errmsg.NewNormalErrorInfo("nats message marshal error", err.Error())
	}
	outMsg, e := this_.conn.Request(msgName, msgBytes, time.Second*10)
	if e != nil {
		return nil, errmsg.NewProtocolErrorInfo(e.Error())
	}

	return outMsg.Data, nil
}

// Shutdown 关闭NATS
func (this_ *NatsClient) Shutdown() {
	if atomic.CompareAndSwapInt32(&this_.closed, 0, 1) {
		_ = this_.conn.FlushTimeout(time.Second * 3)
		this_.conn.Close()
	}
}

// Subscribe 订阅主题
func (this_ *NatsClient) Subscribe(subj string, h nats.MsgHandler) {
	if _, ok := this_.subs.Get(subj); ok {
		panic(fmt.Sprintf("subj [%s] had Subscribed", subj))
	}
	this_.log.Info("Subscribe", zap.String("urls", this_.urls), zap.String("subj", subj))
	sub, err := this_.conn.Subscribe(subj, h)
	utils.Must(err)
	this_.subs.Set(subj, sub)
}

type UserSubject struct {
	GameServerId int64
	RoleId       int64
	MsgName      string
}

// SubscribeHandler 订阅Topic
func (this_ *NatsClient) SubscribeHandler(subj string, f func(*nats.Msg)) {
	if _, ok := this_.subs.Get(subj); ok {
		panic(fmt.Sprintf("subj [%s] had Subscribed", subj))
	}
	group := strings.ReplaceAll(subj, ">", "group")
	this_.log.Info("SubscribeHandler", zap.String("urls", this_.urls), zap.String("subj", subj), zap.String("group", group))
	cb := this_.f
	if f != nil {
		cb = f
	}
	sub, err := this_.conn.QueueSubscribe(subj, group, cb)
	utils.Must(err)
	this_.subs.Set(subj, sub)
}

// UnSub 解除订阅
func (this_ *NatsClient) UnSub(subj string) {
	if s, ok := this_.subs.Get(subj); ok {
		this_.log.Info("Unsubscribe", zap.String("subj", subj))
		_ = s.Unsubscribe()
		this_.subs.Remove(subj)
	}
}

// SubscribeBroadcast 订阅广播主题
func (this_ *NatsClient) SubscribeBroadcast(subj string, f func(*nats.Msg)) {
	if _, ok := this_.subs.Get(subj); ok {
		panic(fmt.Sprintf("subj [%s] had Subscribed", subj))
	}
	this_.log.Info("SubscribeBroadcast", zap.String("urls", this_.urls), zap.String("subj", subj))
	cb := this_.f
	if f != nil {
		cb = f
	}
	sub, err := this_.conn.Subscribe(subj, cb)
	utils.Must(err)
	this_.subs.Set(subj, sub)
}