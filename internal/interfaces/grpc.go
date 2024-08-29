package interfaces

import (
	jwt2 "github.com/golang-jwt/jwt/v5"
	v1 "github.com/ofavor/kratos-layout/api/gen/helloworld/v1"
	"github.com/ofavor/kratos-layout/internal/application"
	"github.com/ofavor/kratos-layout/internal/conf"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(
	logger log.Logger,
	tp *tracesdk.TracerProvider,
	bc *conf.Bootstrap,
	greeter *application.GreeterAppService,
	// TODO: add new service here
) *grpc.Server {
	c := bc.Server
	ac := bc.Auth
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
			logging.Server(logger),
			tracing.Server(tracing.WithTracerProvider(tp)),
			jwt.Server(func(token *jwt2.Token) (interface{}, error) {
				return []byte(ac.Key), nil
			}, jwt.WithSigningMethod(jwt2.SigningMethodHS256), jwt.WithClaims(func() jwt2.Claims {
				return &jwt2.MapClaims{}
			})),
		),
	}
	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}
	srv := grpc.NewServer(opts...)
	v1.RegisterGreeterServer(srv, greeter)
	return srv
}
