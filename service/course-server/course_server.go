package course_server

import (
	baseServer "edu-project/base-server"
	"edu-project/common/errmsg"
	"edu-project/common/utils"
	"edu-project/consts"
	"edu-project/models"
	"edu-project/service/course-server/api"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
	"strings"
)

type Service struct {
	svc      *baseServer.BaseService
	log      *zap.Logger
	consumer map[string]func(*nats.Msg) *errmsg.ErrMsg
}

func NewService(svc *baseServer.BaseService, log *zap.Logger) *Service {
	s := &Service{
		svc:      svc,
		log:      log,
		consumer: make(map[string]func(*nats.Msg) *errmsg.ErrMsg),
	}
	return s
}

func (this_ *Service) Router() {
	this_.routerModel()
	this_.routerFunc()
	this_.registerNatsMessage()
	this_.subscribeMsg()
}

func (this_ *Service) routerModel() {
	utils.Must(this_.svc.SQlDb().AutoMigrate(&models.Student{}))
	utils.Must(this_.svc.SQlDb().AutoMigrate(&models.Course{}))
	utils.Must(this_.svc.SQlDb().AutoMigrate(&models.Enrollment{}))
}

func (this_ *Service) routerFunc() {
	this_.svc.SetMux("/queryCourse", api.QueryCourses)
	this_.svc.SetMux("/findStudentToCourse", api.FindStudentToCourse)
	this_.svc.SetMux("/findCourseToStudent", api.FindCourseToStudent)
}

func (this_ *Service) registerNatsMessage() {
	this_.consumer[consts.SaveEnrollments] = this_.SaveEnrollments
}

func (this_ *Service) subscribeMsg() {
	this_.svc.NatsClient().SubscribeBroadcast(consts.NatsMsgPrefixCourse+">", func(msg *nats.Msg) {
		msgName := strings.TrimPrefix(msg.Subject, consts.NatsMsgPrefixCourse)
		f, ok := this_.consumer[msgName]
		if ok {
			if f(msg) != nil {
				// 重试机制设计
			}
		}
	})
}
