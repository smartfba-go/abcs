package httprouter

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"httprouter",
	fx.Provide(New),
)

type Params struct {
	fx.In

	Handlers []Handler `group:"router"`
}

type Result struct {
	fx.Out

	Handler http.Handler
}

func New(lc fx.Lifecycle, p Params) (Result, error) {
	router := httprouter.New()

	for _, h := range p.Handlers {
		if h.Handle != nil {
			router.Handle(h.Method, h.Path, h.Handle)
		} else if h.Handler != nil {
			router.Handler(h.Method, h.Path, h.Handler)
		} else if h.HandlerFunc != nil {
			router.HandlerFunc(h.Method, h.Path, h.HandlerFunc)
		}
	}

	return Result{
		Handler: router,
	}, nil
}

type Handler struct {
	Method      string
	Path        string
	Handle      httprouter.Handle
	Handler     http.Handler
	HandlerFunc http.HandlerFunc
}
