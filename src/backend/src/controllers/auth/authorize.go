package auth

import (
	"errors"
	"fmt"
	"nak-auth/services"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

type AuthController struct {
	clientService services.IClientService
}

func NewAuthController(uService services.IClientService) *AuthController {
	return &AuthController{clientService: uService}
}

func (*AuthController) Path() string {
	return "/oauth/authorize"
}

type Authorize struct {
	ResponseType string   `json:"response_type"`
	ClientId     string   `json:"client_id"`
	Scope        []string `json:"scope"`
	RedirectUri  string   `json:"redirect_uri"`
	GrantType    string   `json:"grant_type"`
}

func (l *AuthController) WriteResponse(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		auth, err := checkAuthorizeParams(r.URL)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		if shouldRedirect(auth, l.clientService) {
			oauth_state := r.URL.Query().Encode()
			http.Redirect(w, r, fmt.Sprintf("/login?%s", oauth_state), http.StatusTemporaryRedirect)
		} else {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	}
}

// type AuthCodeController struct {
// 	clientService services.IClientService
// }

// func (*AuthCodeController) Path() string {
// 	return "/login/oauth/authorize/code"
// }

// func (l *AuthCodeController) WriteResponse(w http.ResponseWriter, r *http.Request) {
// 	switch r.Method {
// 	case "GET":
// 		redirect_uri := r.URL.Query().Get("redirect_uri")
// 		auth, err := checkAuthorizeParams(r.URL)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusBadRequest)
// 			return
// 		}
// 		if shouldRedirect(auth, l.clientService) {
// 			oauth_state := r.URL.Query().Encode()
// 			http.Redirect(w, r, redirect_uri, http.StatusTemporaryRedirect)
// 		} else {
// 			http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 			return
// 		}
// 	}
// }

func checkAuthorizeParams(url *url.URL) (Authorize, error) {
	query := url.Query()

	response_type := query.Get("response_type")
	client_id := query.Get("client_id")
	scope := strings.Split(query.Get("scope"), " ")
	redirect_uri := query.Get("redirect_uri")
	grant_type := query.Get("grant_type")

	if response_type == "" || client_id == "" || scope[0] == "" || redirect_uri == "" || grant_type == "" {
		return Authorize{}, errors.New("Missing query params. response_type, client_id, scope, redirect_uri, grant_type are all required.")
	}

	return Authorize{
		ResponseType: response_type,
		ClientId:     client_id,
		Scope:        scope,
		RedirectUri:  redirect_uri,
		GrantType:    grant_type,
	}, nil
}

func shouldRedirect(auth Authorize, clientService services.IClientService) bool {

	client, err := clientService.GetByID(auth.ClientId)
	if err != nil {
		return false
	}
	var clientScopes = []string{}
	for i := 0; i < len(client.Scopes); i++ {
		scope := client.Scopes[i].From()
		clientScopes = append(clientScopes, scope)
	}
	return client.GrantType == auth.GrantType && client.RedirectURI == auth.RedirectUri && reflect.DeepEqual(clientScopes, auth.Scope)
}
