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

func NewServeMux(routes []controllers.ApiHandle, pages []controllers.PageHandle) *mux.Router {
	mux := mux.NewRouter()
	for _, route := range routes {
		mux.HandleFunc(fmt.Sprintf("/api%s", route.Path()), route.WriteResponse)
	}
	for _, page := range pages {
		mux.Handle(page.Path(), page)
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

func AsPageHandle(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(controllers.PageHandle)),
		fx.ResultTags(`group:"pages"`),
	)
}

func AsSingleton(implementation, specification interface{}) interface{} {
	return fx.Annotate(
		implementation,
		fx.As(specification),
	)
}

func Startup(*http.Server, *gorm.DB, services.IClientService, services.IUserService, services.ILoginService, services.ITokenService) {

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
				fx.ParamTags(`group:"routes"`, `group:"pages"`),
			),
			// Adds Services
			AsSingleton(services.NewClientService, new(services.IClientService)),
			AsSingleton(services.NewUserService, new(services.IUserService)),
			AsSingleton(services.NewLoginService, new(services.ILoginService)),
			AsSingleton(services.NewTokenService, new(services.ITokenService)),

			// Adds API Controllers
			AsApiHandle(controllers.NewHealthController),
			AsApiHandle(controllers.NewUserController),
			AsApiHandle(controllers.NewUserByIdController),
			AsApiHandle(controllers.NewClientController),
			AsApiHandle(controllers.NewClientByIdController),
			AsApiHandle(auth.NewAuthController),
			AsApiHandle(auth.NewAccessController),
			AsApiHandle(auth.NewLoginController),

			// Adds Pages
			AsPageHandle(pages.NewHomePage),
			AsPageHandle(pages.NewLoginPage),
			AsPageHandle(pages.NewClientsPage),
		),
		fx.Invoke(Startup),
	).Run()
}
