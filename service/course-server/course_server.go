package course_server

import (
	"chose-course/common/utils"
	"chose-course/module"
	"chose-course/service"
	"chose-course/service/course-server/api"
	"go.uber.org/zap"
)

type Service struct {
	svc *service.BaseService
	log *zap.Logger
}

func NewService(svc *service.BaseService, log *zap.Logger) *Service {
	s := &Service{
		svc: svc,
		log: log,
	}
	return s
}

func (this_ *Service) Router() {
	this_.routerModel()
	this_.routerFunc()
}

func (this_ *Service) routerModel() {
	utils.Must(this_.svc.SQlDb().AutoMigrate(&module.Student{}))
}

func (this_ *Service) routerFunc() {
	this_.svc.Mux().HandleFunc("/findStudent", api.FindStudent)
}
