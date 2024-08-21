package iface

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/ofavor/ddd-go/pkg/event"
	"github.com/ofavor/kratos-layout/internal/app"
	"github.com/ofavor/kratos-layout/internal/infra"
)

type EventHandler struct {
	logger  *log.Helper
	bus     event.EventBus
	myEvent *app.MyEventAppService
}

func NewEventHandler(infra *infra.Infra, myEvent *app.MyEventAppService, logger log.Logger) *EventHandler {
	return &EventHandler{
		logger: log.NewHelper(logger),
		bus:    infra.GetEvent(),

		myEvent: myEvent,
	}
}

func (h *EventHandler) Initialize() error {
	h.bus.Subscribe("greeter.created", "my-event", h.myEvent.OnGreeterCreated)
	return nil
}
