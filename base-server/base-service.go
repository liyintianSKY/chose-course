package base_server

import (
	"chose-course/common/natsclient"
	"chose-course/common/utils"
	"github.com/lesismal/nbio/nbhttp"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
)

type BaseService struct {
	logger      *zap.Logger
	db          *gorm.DB
	redisClient *redis.Client
	httpMux     *http.ServeMux
	nbServer    *nbhttp.Server
	nc          *natsclient.NatsClient
}

func InitBaseServer(log *zap.Logger, db *gorm.DB, redis *redis.Client, nc *natsclient.NatsClient, httpListen string) *BaseService {
	s := &BaseService{
		logger:      log,
		db:          db,
		redisClient: redis,
		httpMux:     &http.ServeMux{},
		nbServer:    httpServer(httpListen),
		nc:          nc,
	}
	return s
}

func (this_ *BaseService) SQlDb() *gorm.DB {
	return this_.db
}

func (this_ *BaseService) RedisClient() *redis.Client {
	return this_.redisClient
}

func (this_ *BaseService) SetMux(url string, f func(w http.ResponseWriter, r *http.Request)) {
	this_.httpMux.Handle(url, http.HandlerFunc(f))
}

func (this_ *BaseService) NatsClient() *natsclient.NatsClient {
	return this_.nc
}

func (this_ *BaseService) Nbhttp() *nbhttp.Server {
	return this_.nbServer
}

func httpServer(httpListen string) *nbhttp.Server {
	return nbhttp.NewServer(nbhttp.Config{
		Name:    "course-http",
		Network: "tcp",
		//TLSConfig: tlsC,
		Addrs: []string{httpListen},
	})
}

func (this_ *BaseService) StartHttp() {
	this_.nbServer.Handler = this_.httpMux
	err := this_.nbServer.Start()
	utils.Must(err)
}

func (this_ *BaseService) StopHttp() {
	this_.nbServer.Stop()
}
