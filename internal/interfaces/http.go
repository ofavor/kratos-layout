package interfaces

import (
	"context"

	"github.com/gorilla/handlers"
	v1 "github.com/ofavor/kratos-layout/api/gen/helloworld/v1"
	"github.com/ofavor/kratos-layout/internal/application"
	"github.com/ofavor/kratos-layout/internal/conf"

	nethttp "net/http"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-kratos/swagger-api/openapiv2"
	jwt2 "github.com/golang-jwt/jwt/v5"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

func newAuthWhiteListMatcher() selector.MatchFunc {
	whiteList := make(map[string]struct{})
	whiteList["/helloworld.v1.Greeter/Create"] = struct{}{}
	whiteList["/helloworld.v1.Greeter/SayHello"] = struct{}{}
	return func(ctx context.Context, operation string) bool {
		if _, ok := whiteList[operation]; ok {
			return false
		}
		return true
	}
}

func responseEncoder(
	w http.ResponseWriter,
	r *http.Request,
	i interface{},
) error {
	type response struct {
		Code     int               `json:"code"`
		Message  string            `json:"message"`
		Reason   string            `json:"reason"`
		Metadata map[string]string `json:"metadata"`
		Data     interface{}       `json:"data,omitempty"`
	}
	// reply := &response{
	// 	Code:    200,
	// 	Message: "success",
	// 	Data:    i,
	// }
	// codec, _ := http.CodecForRequest(r, "Accept")
	// data, err := codec.Marshal(reply)
	// if err != nil {
	// 	return err
	// }
	// w.Header().Set("Content-Type", "application/json")
	// w.Write(data)
	// return nil

	// 解决0值被忽略的问题
	codec, _ := http.CodecForRequest(r, "Accept")
	data, err := codec.Marshal(i)
	if err != nil {
		return err
	}
	w.WriteHeader(nethttp.StatusOK)

	reply := &response{
		Code:    200,
		Message: "success",
	}

	replyData, err := codec.Marshal(reply)
	if err != nil {
		return err
	}

	var newData = make([]byte, 0, len(replyData)+len(data)+8)
	newData = append(newData, replyData[:len(replyData)-1]...)
	newData = append(newData, []byte(`,"data":`)...)
	newData = append(newData, data...)
	newData = append(newData, '}')

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(newData)
	if err != nil {
		return err
	}
	return nil
}

// NewHTTPServer new an HTTP server.
func NewHTTPServer(
	logger log.Logger,
	tp *tracesdk.TracerProvider,
	bc *conf.Bootstrap,
	greeter *application.GreeterAppService,
	// TODO: add new service here
) *http.Server {
	c := bc.Server
	ac := bc.Auth
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			logging.Server(logger),
			tracing.Server(tracing.WithTracerProvider(tp)),
			selector.Server(
				jwt.Server(func(token *jwt2.Token) (interface{}, error) {
					return []byte(ac.Key), nil
				}, jwt.WithSigningMethod(jwt2.SigningMethodHS256), jwt.WithClaims(func() jwt2.Claims {
					return &jwt2.MapClaims{}
				})),
			).Match(newAuthWhiteListMatcher()).Build(),
		),
		http.Filter(handlers.CORS(
			handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}),
			handlers.AllowedOrigins([]string{"*"}),
		)),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	opts = append(opts, http.ResponseEncoder(responseEncoder))
	srv := http.NewServer(opts...)
	h := openapiv2.NewHandler()
	srv.HandlePrefix("/q/", h)
	v1.RegisterGreeterHTTPServer(srv, greeter)
	return srv
}
