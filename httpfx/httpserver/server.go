package httpserver

import (
	"context"
	"net"
	"net/http"

	abczap "go.smartfba.io/abcs/log/zap"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"httpserver",
	fx.Provide(NewServer),
)

type ServerParams struct {
	fx.In

	Config     ServerConfig
	Handler    http.Handler
	Middleware []Middleware `group:"server"`

	Logger abczap.Logger `optional:"true"`
}

type ServerResult struct {
	fx.Out

	Server *http.Server
}

type ServerConfig struct {
	Addr string
}

type Middleware func(w http.ResponseWriter, r *http.Request, next http.Handler)

func chainMiddlewares(h http.Handler, ms []Middleware) http.Handler {
	if len(ms) == 0 {
		return h
	}

	m := ms[0]
	n := chainMiddlewares(h, ms[1:])

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m(w, r, n)
	})
}

func NewServer(lc fx.Lifecycle, p ServerParams) (ServerResult, error) {
	server := &http.Server{
		Addr: p.Config.Addr,
	}

	handler := p.Handler
	if len(p.Middleware) > 0 {
		handler = chainMiddlewares(p.Handler, p.Middleware)
	}

	handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(abczap.WithLogger(r.Context(), p.Logger))

		p.Handler.ServeHTTP(w, r)
	})

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			lis, err := net.Listen("tcp", p.Config.Addr)
			if err != nil {
				return err
			}

			server.Handler = handler

			go func() {
				if err := server.Serve(lis); err != http.ErrServerClosed {
					// TODO: the server failed with an error, what do we do?
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			return server.Shutdown(ctx)
		},
	})

	return ServerResult{
		Server: server,
	}, nil
}
