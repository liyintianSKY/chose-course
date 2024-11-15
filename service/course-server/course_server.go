package course_server

import (
	baseServer "chose-course/base-server"
	"chose-course/common/utils"
	"chose-course/models"
	"chose-course/service/course-server/api"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type Service struct {
	svc *baseServer.BaseService
	log *zap.Logger
	//consumer map[string]func()
}

func NewService(svc *baseServer.BaseService, log *zap.Logger) *Service {
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
	utils.Must(this_.svc.SQlDb().AutoMigrate(&models.Student{}))
	utils.Must(this_.svc.SQlDb().AutoMigrate(&models.Course{}))
	utils.Must(this_.svc.SQlDb().AutoMigrate(&models.Enrollment{}))
}

func (this_ *Service) routerFunc() {
	this_.svc.SetMux("/findStudent", api.FindStudent)
}

func (this_ *Service) subscribeMsg() {
	this_.svc.NatsClient().SubscribeBroadcast("course.>", func(msg *nats.Msg) {

	})
}
