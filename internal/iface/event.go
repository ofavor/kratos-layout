package iface

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/ofavor/ddd-go/pkg/event"
	"github.com/ofavor/kratos-layout/internal/app"
)

type EventHandler struct {
	logger  *log.Helper
	bus     event.EventBus
	myEvent *app.MyEventAppService
}

func NewEventHandler(
	logger log.Logger,
	bus event.EventBus,
	myEvent *app.MyEventAppService,
	// TODO: add new service here
) *EventHandler {
	return &EventHandler{
		logger: log.NewHelper(logger),
		bus:    bus,

		myEvent: myEvent,
	}
}

func (h *EventHandler) Initialize() error {
	h.bus.Subscribe("greeter.created", "my-event", h.myEvent.OnGreeterCreated)
	return nil
}
