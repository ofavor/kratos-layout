package app

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/ofavor/ddd-go/pkg/event"
)

type MyEventAppService struct {
	logger *log.Helper
}

func NewMyEventAppService(logger log.Logger) *MyEventAppService {
	return &MyEventAppService{
		logger: log.NewHelper(logger),
	}
}

func (s *MyEventAppService) OnGreeterCreated(e *event.Event) {
	s.logger.Debugf("on greeter created: %s", e)
}
