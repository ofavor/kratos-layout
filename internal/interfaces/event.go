package interfaces

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/ofavor/ddd-go/pkg/event"
	"github.com/ofavor/kratos-layout/internal/application"
)

type EventHandler struct {
	logger     *log.Helper
	bus        event.EventBus
	gretAppSvc *application.GreeterAppService
}

func NewEventHandler(
	logger log.Logger,
	bus event.EventBus,
	gretAppSvc *application.GreeterAppService,
	// TODO: add new service here
) *EventHandler {
	return &EventHandler{
		logger: log.NewHelper(logger),
		bus:    bus,

		gretAppSvc: gretAppSvc,
	}
}

func (h *EventHandler) Initialize() error {
	h.bus.Subscribe("greeter.created", "greeter", h.gretAppSvc.OnGreeterCreated)
	return nil
}
