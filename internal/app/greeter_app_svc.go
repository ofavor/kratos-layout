package app

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/ofavor/ddd-go/pkg/event"
	v1 "github.com/ofavor/kratos-layout/api/gen/helloworld/v1"
	"github.com/ofavor/kratos-layout/internal/domain/entity"
	"github.com/ofavor/kratos-layout/internal/domain/repository"
	"github.com/ofavor/kratos-layout/internal/domain/vo"
	"github.com/ofavor/kratos-layout/internal/infra"
)

type GreeterAppService struct {
	v1.UnimplementedGreeterServer

	repo   repository.GreeterRepo
	event  event.EventBus
	logger *log.Helper
}

func NewGreeterAppService(infra *infra.Infra, logger log.Logger) *GreeterAppService {
	return &GreeterAppService{
		repo:   infra.GetGreeterRepo(),
		event:  infra.GetEvent(),
		logger: log.NewHelper(logger),
	}
}

func (s *GreeterAppService) SayHello(ctx context.Context, in *v1.HelloRequest) (*v1.HelloReply, error) {
	gt, err := s.repo.Get(nil, in.Id)
	if err != nil {
		return nil, err
	}
	g := gt.SayHello()
	return &v1.HelloReply{Message: g}, nil
}

func (s *GreeterAppService) Create(ctx context.Context, in *v1.CreateRequest) (*v1.CreateResponse, error) {
	gt, err := entity.NewGreeter(in.Name, in.Greeting, vo.NewAddress("A", "B", "C", "D"))
	if err != nil {
		return nil, err
	}
	err = s.repo.Save(nil, gt)
	if err != nil {
		return nil, err
	}
	s.event.Publish("greeter.created", map[string]interface{}{"id": gt.GetId()})
	return &v1.CreateResponse{Id: int64(gt.GetId())}, nil
}
