package controllers
import "net/http"

type Route struct {
	Path 		string
	Method 	string
}

type ApiHandle interface {
	Path() 	string
	WriteResponse(
		w http.ResponseWriter,
		r *http.Request,
	)
}