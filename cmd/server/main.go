package main

import (
	"flag"
	"os"

	"github.com/ofavor/kratos-layout/internal/conf"
	"github.com/ofavor/kratos-layout/internal/infrastructure"
	"github.com/ofavor/kratos-layout/internal/interfaces"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/env"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"

	_ "go.uber.org/automaxprocs"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf string

	id, _ = os.Hostname()
)

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
}

func newApp(
	logger log.Logger,
	gs *grpc.Server,
	hs *http.Server,
	rr registry.Registrar,
	eh *interfaces.EventHandler,
	infra *infrastructure.Infra,
) *kratos.App {
	app := kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			gs,
			hs,
		),
	)
	kratos.Registrar(rr)
	infra.Initialize()
	eh.Initialize()
	return app
}

func main() {
	flag.Parse()

	c := config.New(
		config.WithSource(
			env.NewSource(""),
			file.NewSource(flagconf),
		),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	var rc conf.Registry
	if err := c.Scan(&rc); err != nil {
		panic(err)
	}

	// init logger
	logger := log.NewFilter(
		log.With(log.NewStdLogger(os.Stdout),
			"ts", log.DefaultTimestamp,
			"caller", log.DefaultCaller,
			"service.id", id,
			"service.name", Name,
			"service.version", Version,
			"trace.id", tracing.TraceID(),
			"span.id", tracing.SpanID(),
		),
		log.FilterLevel(log.ParseLevel(bc.Logging.Level)),
	)

	// init tracer
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(bc.Trace.Endpoint)))
	if err != nil {
		panic(err)
	}
	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(resource.NewSchemaless(
			semconv.ServiceNameKey.String(Name),
		)),
	)

	app, cleanup, err := wireApp(&bc, &rc, logger, tp)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
