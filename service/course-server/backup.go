package course_server

import (
	"chose-course/common/errmsg"
	"github.com/nats-io/nats.go"
)

func (this_ *Service) SaveEnrollments(msg *nats.Msg) *errmsg.ErrMsg {
	this_.log.Info("测试数据" + string(msg.Data))
	return nil
}
