package main

import (
	"context"
	"fmt"
	"nak-auth/controllers"
	"nak-auth/db"
	"nak-auth/services"
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

func NewHTTPServer(lc fx.Lifecycle, mux *mux.Router) *http.Server {
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

func NewServeMux(routes []controllers.ApiHandle) *mux.Router {
	mux := mux.NewRouter()
	for _, route := range routes {
		mux.HandleFunc(route.Path(), route.WriteResponse)
	}
	return mux
}

func AsApiHandle(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(controllers.ApiHandle)),
		fx.ResultTags(`group:"routes"`),
	)
}

func main() {
	fx.New(
		fx.Provide(
			NewHTTPServer,
			db.NewPScaleClient,
			fx.Annotate(
				NewServeMux,
				fx.ParamTags(`group:"routes"`),
			),
			fx.Annotate(
				services.NewClientService,
				fx.As(new(services.IClientService)),
			),
			AsApiHandle(controllers.NewHealthController),
			AsApiHandle(controllers.NewCounterController),
			AsApiHandle(controllers.NewUserController),
			AsApiHandle(controllers.NewUserByIdController),
			AsApiHandle(controllers.NewClientController),
			AsApiHandle(controllers.NewClientByIdController),
		),
		fx.Invoke(func(*http.Server, *gorm.DB, services.IClientService) {}),
	).Run()
}
