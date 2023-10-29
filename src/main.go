package main

import (
	"context"
	"fmt"
	"log"
	"nak-auth/controllers"
	"nak-auth/controllers/auth"
	"nak-auth/controllers/pages"
	"nak-auth/db"
	"nak-auth/services"
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
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
	// mux.PathPrefix("/").Handler(http.FileServer((http.Dir("./static/"))))
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

func AsSingleton(implementation, specification interface{}) interface{} {
	return fx.Annotate(
		implementation,
		fx.As(specification),
	)
}

func InjectThis(*http.Server, *gorm.DB, services.IClientService, services.IUserService, services.ILoginService, services.ITokenService) {
	// This is here to force fx to inject the dependencies
	// into the functions above. Otherwise, they will not
	// be injected.
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	fx.New(
		fx.Provide(
			NewHTTPServer,
			db.NewLibSqlClient,
			fx.Annotate(
				NewServeMux,
				fx.ParamTags(`group:"routes"`),
			),
			// Adds Services
			AsSingleton(services.NewClientService, new(services.IClientService)),
			AsSingleton(services.NewUserService, new(services.IUserService)),
			AsSingleton(services.NewLoginService, new(services.ILoginService)),
			AsSingleton(services.NewTokenService, new(services.ITokenService)),

			// Adds Controllers
			AsApiHandle(controllers.NewHealthController),
			AsApiHandle(controllers.NewCounterController),
			AsApiHandle(controllers.NewUserController),
			AsApiHandle(controllers.NewUserByIdController),
			AsApiHandle(controllers.NewClientController),
			AsApiHandle(controllers.NewClientByIdController),
			AsApiHandle(auth.NewAuthController),
			AsApiHandle(auth.NewAccessController),
			AsApiHandle(auth.NewLoginController),
			AsApiHandle(pages.NewHomeController),
		),
		fx.Invoke(InjectThis),
	).Run()
}
