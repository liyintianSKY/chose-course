package service

import (
	"database/sql"
	"github.com/lesismal/nbio/nbhttp"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"net/http"
)

type BaseService struct {
	logger      *zap.Logger
	db          *sql.DB
	redisClient *redis.Client
	httpMux     *http.ServeMux
	nbServer    *nbhttp.Server
}

func InitServer(log *zap.Logger, db *sql.DB, redis *redis.Client, httpListen string) *BaseService {
	s := &BaseService{
		logger:      log,
		db:          db,
		redisClient: redis,
		httpMux:     &http.ServeMux{},
		nbServer:    httpServer(httpListen),
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
	this_.RouterModel()
	this_.RouterFunc()
}

func (this_ *BaseService) RouterModel() {}

func (this_ *BaseService) RouterFunc() {}
