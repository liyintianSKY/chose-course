package service

import (
	"chose-course/common/natsclient"
	course_server "chose-course/service/course-server"
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

func (this_ *BaseService) SQlDb() *gorm.DB {
	return this_.db
}

func (this_ *BaseService) RedisClient() *redis.Client {
	return this_.redisClient
}

func (this_ *BaseService) Mux() *http.ServeMux {
	return this_.httpMux
}

func (this_ *BaseService) NatsClient() *natsclient.NatsClient {
	return this_.nc
}

func InitServer(log *zap.Logger, db *gorm.DB, redis *redis.Client, nc *natsclient.NatsClient, httpListen string) *BaseService {
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

func (this_ *BaseService) Start() {
	this_.Route()
	this_.nbServer.Handler = this_.httpMux
	err := this_.nbServer.Start()
	if err != nil {
		panic(err)
	}
}

func (this_ *BaseService) Stop() {
	this_.nbServer.Stop()
	//this_.nc.Close()
}

func httpServer(httpListen string) *nbhttp.Server {
	return nbhttp.NewServer(nbhttp.Config{
		Name:    "course-http",
		Network: "tcp",
		//TLSConfig: tlsC,
		Addrs: []string{httpListen},
	})
}

func (this_ *BaseService) Route() {
	course_server.NewService(this_, this_.logger).Router()
}
