package main

import (
	"log"
	"nak-auth/db"
	"nak-auth/models"
	"nak-auth/services"
	"nak-auth/templates"
	"net/http"
	"strconv"
	"strings"

	"github.com/ShounakA/roids"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// fx.New(
	// 	fx.Provide(
	// 		// Adds API Controllers
	// 		AsApiHandle(auth.NewAuthController),
	// 		AsApiHandle(auth.NewTokenController),
	// 		AsApiHandle(auth.NewLoginController),

	// 		// Adds Pages
	// 		AsPageHandle(pages.NewClientsPage),
	// 	),
	// 	fx.Invoke(Startup),
	// ).Run()

	app := echo.New()

	_ = roids.GetRoids()

	isReady := false

	err = roids.AddStaticService(new(db.ILibSqlClientFactory), db.NewLibSqlClientFactory)
	if err != nil {
		log.Fatal("Did not bind service.", err.Error())
		return
	}
	err = roids.AddTransientService(new(services.IClientService), services.NewClientService)
	if err != nil {
		log.Fatal("Did not bind service.", err.Error())
		return
	}
	err = roids.AddTransientService(new(services.IUserService), services.NewUserService)
	if err != nil {
		log.Fatal("Did not bind service.", err.Error())
		return
	}
	err = roids.AddTransientService(new(services.ILoginService), services.NewLoginService)
	if err != nil {
		log.Fatal("Did not bind service.", err.Error())
		return
	}
	err = roids.AddTransientService(new(services.ITokenService), services.NewTokenService)
	if err != nil {
		log.Fatal("Did not bind service.", err.Error())
		return
	}

	roids.Build()

	libSqlFact := roids.Inject[db.ILibSqlClientFactory]()
	dbConn, err := libSqlFact.CreateClient()
	if err != nil {
		panic("NO DB CONNECTION")
	}

	dbConn.AutoMigrate(&models.Code{}, &models.User{}, &models.Client{})

	app.Renderer = templates.TemplateFactory()

	app.GET("/api/health/live", func(c echo.Context) error {
		return c.String(http.StatusOK, "Live")
	})

	app.GET("/api/health/ready", func(c echo.Context) error {
		if !isReady {
			return echo.NewHTTPError(http.StatusServiceUnavailable, "Not Ready")
		}
		return c.String(http.StatusOK, "Ready")
	})

	app.GET("/api/users", func(c echo.Context) error {
		userService := roids.Inject[services.IUserService]()

		users, dbErr := userService.GetAll()
		if dbErr != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
		}
		return c.JSON(http.StatusOK, users)
	})

	app.POST("/api/users", func(c echo.Context) error {
		var usrBody models.User
		userService := roids.Inject[services.IUserService]()
		if err := c.Bind(&usrBody); err != nil {
			return err
		}
		dbErr := userService.Create(usrBody)
		if dbErr != nil {
			return echo.ErrInternalServerError
		}
		return c.JSON(http.StatusOK, usrBody)
	})

	app.GET("/api/users/:id", func(c echo.Context) error {
		var id string
		userService := roids.Inject[services.IUserService]()
		if err := (&echo.DefaultBinder{}).BindPathParams(c, &id); err != nil {
			return echo.ErrBadRequest
		}
		marks, err := strconv.Atoi(id)
		if err != nil {
			return echo.ErrBadGateway
		}
		user, dbErr := userService.GetByID(marks)
		if dbErr != nil {
			return echo.ErrInternalServerError
		}
		return c.JSON(http.StatusOK, user)
	})

	app.DELETE("/api/users/:id", func(c echo.Context) error {
		var id string
		userService := roids.Inject[services.IUserService]()
		if err := (&echo.DefaultBinder{}).BindPathParams(c, &id); err != nil {
			return echo.ErrBadRequest
		}
		marks, err := strconv.Atoi(id)
		if err != nil {
			return echo.ErrBadGateway
		}
		if userService.Delete(marks) != nil {
			return echo.ErrInternalServerError
		}
		return c.JSON(http.StatusOK, "Deleted")
	})

	app.GET("/api/clients", func(c echo.Context) error {
		clientService := roids.Inject[services.IClientService]()
		clients, dbErr := clientService.GetAll()
		if dbErr != nil {
			return echo.ErrInternalServerError
		}
		client_json := models.ListOfClientsToListOfClientJson(clients)

		// if c.Request().Header["Hx-Request"][0] == "true" {
		// 	templates.WriteClientsFragment(w, clients)
		// } else {
		// json.NewEncoder(w).Encode(client_json)
		return c.JSON(http.StatusOK, client_json)

		// }
	})

	app.POST("/api/clients", func(c echo.Context) error {
		var client models.ClientJson
		clientService := roids.Inject[services.IClientService]()
		if err := c.Bind(&client); err != nil {
			return err
		}
		if strings.TrimSpace(client.Name) == "" {
			return echo.ErrBadRequest
		}
		if strings.TrimSpace(client.GrantType) == "" {
			return echo.ErrBadRequest
		}

		//Create the user
		err := clientService.Create(client)
		if err != nil {
			return echo.ErrInternalServerError
		}
		return c.JSON(http.StatusCreated, client)
	})

	app.GET("/api/clients/:id", func(c echo.Context) error {
		var id string
		clientService := roids.Inject[services.IClientService]()
		if err := (&echo.DefaultBinder{}).BindPathParams(c, &id); err != nil {
			return echo.ErrBadRequest
		}
		client, err := clientService.GetByID(id)
		if err != nil {
			return echo.ErrInternalServerError
		}
		client_json := client.From()
		return c.JSON(http.StatusOK, client_json)
	})

	app.DELETE("/api/clients/:id", func(c echo.Context) error {
		var id string
		clientService := roids.Inject[services.IClientService]()
		if err := (&echo.DefaultBinder{}).BindPathParams(c, &id); err != nil {
			return echo.ErrBadRequest
		}
		err := clientService.Delete(id)
		if err != nil {
			return echo.ErrInternalServerError
		}
		return c.String(http.StatusNoContent, "Deleted")
	})

	app.GET("/client", func(c echo.Context) error {
		loginService := roids.Inject[services.ILoginService]()
		clientService := roids.Inject[services.IClientService]()
		if !loginService.ClientIsAuthenticated(c.Request()) {
			return c.Redirect(http.StatusSeeOther, "/login")
		}
		clients, err := clientService.GetAll()
		if err != nil {
			return echo.ErrInternalServerError
		}
		return c.Render(http.StatusOK, "clientLayout", clients)
	})

	app.GET("/login", func(c echo.Context) error {
		var redirect_uri, issuer string
		if err := (&echo.DefaultBinder{}).BindPathParams(c, &redirect_uri); err != nil {
			redirect_uri = "http://localhost:8080"
		}
		if err = (&echo.DefaultBinder{}).BindPathParams(c, &issuer); err != nil {
			issuer = "nak-auth"
		}

		return c.Render(http.StatusOK, "loginLayout", templates.LoginPageData{Redirect: redirect_uri, Issuer: issuer})
	})

	app.GET("/", func(c echo.Context) error {
		loginService := roids.Inject[services.ILoginService]()
		if !loginService.ClientIsAuthenticated(c.Request()) {
			return c.Redirect(http.StatusSeeOther, "/login")
		}
		return c.Render(http.StatusOK, "homeLayout", nil)
	})

	if app.Start(":8080"); err != http.ErrServerClosed {
		log.Fatal(err)
	}
	isReady = true
}
