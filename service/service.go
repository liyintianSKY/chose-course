package service

import (
	"chose-course/base-server"
	"chose-course/common/natsclient"
	courseserver "chose-course/service/course-server"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Service struct {
	baseServer *base_server.BaseService
	log        *zap.Logger
}

func InitServer(log *zap.Logger, db *gorm.DB, redis *redis.Client, nc *natsclient.NatsClient, httpListen string) *Service {
	s := &Service{
		baseServer: base_server.InitBaseServer(log, db, redis, nc, httpListen),
	}
	return s
}

func (this_ *Service) RegisterSubModule() {
	courseserver.NewService(this_.baseServer, this_.log).Router()
}

func (this_ *Service) Start() {
	this_.baseServer.StartHttp()
}

func (this_ *Service) Stop() {
	this_.baseServer.StopHttp()
	this_.baseServer.NatsClient().Close()
	this_.baseServer.NatsClient().Shutdown()
}
