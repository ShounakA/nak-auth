package main
import (
	"net/http"
	"net"
	"database/sql"
	"fmt"
	"io"
	"context"
	"log"
    "os"
     _ "github.com/go-sql-driver/mysql"
	"go.uber.org/fx"
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

func NewPScaleClient(lc fx.Lifecycle) (*sql.DB, error) {
	db, err := sql.Open("mysql", os.Getenv("DSN"))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
		  	defer db.Close()
			if err := db.Ping(); err != nil {
				log.Fatalf("failed to ping: %v", err)
			} else {
				log.Println("Succesfully connected to planet scale!");
			}
			return err
		},
		OnStop: func(ctx context.Context) error {
			db.Close()
			return nil
		},
	})
	return db, err
}

// EchoHandler is an http.Handler that copies its request body
// back to the response.
type EchoHandler struct{}

// NewEchoHandler builds a new EchoHandler.
func NewEchoHandler() *EchoHandler {
  return &EchoHandler{}
}

// ServeHTTP handles an HTTP request to the /echo endpoint.
func (*EchoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  if _, err := io.Copy(w, r.Body); err != nil {
    fmt.Fprintln(os.Stderr, "Failed to handle request:", err)
  }
}

// NewServeMux builds a ServeMux that will route requests
// to the given EchoHandler.
func NewServeMux(echo *EchoHandler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/echo", echo)
	return mux
}


func main() {
   //  router := gin.Default()
   //  router.GET("/albums", getAlbums)
	fx.New(
		fx.Provide(	NewHTTPServer, 
						NewPScaleClient,
						NewServeMux,
						NewEchoHandler),
		fx.Invoke(func(*http.Server, *sql.DB) {}),
	).Run()

   //  router.Run("localhost:8080")
}