package svc

// import (
// 	"context"
// 	"demo-prod/internal/conf"
// 	usrv1 "demo-usr/api/gen/user/v1"

// 	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
// 	"github.com/go-kratos/kratos/v2/middleware/recovery"
// 	"github.com/go-kratos/kratos/v2/middleware/tracing"
// 	"github.com/go-kratos/kratos/v2/transport/grpc"
// 	jwt2 "github.com/golang-jwt/jwt/v5"
// 	tracesdk "go.opentelemetry.io/otel/sdk/trace"
// )

// func NewUserServiceClient(conf *conf.Bootstrap, tp *tracesdk.TracerProvider) usrv1.UserClient {
// 	ac := conf.Auth
// 	ep := conf.Client.Endpoints["demo-usr"]
// 	conn, err := grpc.DialInsecure(
// 		context.Background(),
// 		grpc.WithEndpoint(ep),
// 		grpc.WithMiddleware(
// 			tracing.Client(tracing.WithTracerProvider(tp)),
// 			recovery.Recovery(),
// 			jwt.Client(func(token *jwt2.Token) (interface{}, error) {
// 				return []byte(ac.Key), nil
// 			}, jwt.WithSigningMethod(jwt2.SigningMethodHS256)),
// 		),
// 	)
// 	if err != nil {
// 		panic(err)
// 	}
// 	c := usrv1.NewUserClient(conn)
// 	return c
// }
