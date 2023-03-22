package main

import (
	"context"
	"fmt"
	"nak-auth/controllers"
	"nak-auth/db"
	"net"
	"net/http"

	"go.uber.org/fx"
	"gorm.io/gorm"
)

func NewHTTPServer(lc fx.Lifecycle, mux *http.ServeMux) *http.Server {
	srv := &http.Server{Addr: ":8080", Handler: mux}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", srv.Addr)
			if err != nil {
				return err
			}
			fmt.Println("Starting HTTP server at", srv.Addr)
			go srv.Serve(ln)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})
	return srv
}

// NewServeMux builds a ServeMux that will route requests
// to the given EchoHandler.
func NewServeMux(routes []controllers.ApiHandle) *http.ServeMux {
	mux := http.NewServeMux()
	for _, route := range routes {
		mux.HandleFunc(route.Path(), route.WriteResponse)
	}
	return mux
}

// AsRoute annotates the given constructor to state that
// it provides a route to the "routes" group.
func AsApiHandle(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(controllers.ApiHandle)),
		fx.ResultTags(`group:"routes"`),
	)
}

func main() {
	fx.New(
		fx.Provide(NewHTTPServer,
			db.NewPScaleClient,
			fx.Annotate(
				NewServeMux,
				fx.ParamTags(`group:"routes"`),
			),
			AsApiHandle(controllers.NewHealthController),
			AsApiHandle(controllers.NewCounterController),
		),
		fx.Invoke(func(*http.Server, *gorm.DB) {}),
	).Run()
}
