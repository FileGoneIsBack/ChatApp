package authentication

import (
	"login/core/models"
	"login/core/sessions"
	"login/core/models/functions"
	"login/core/models/server"
	"login/core/master/api/auth"
	"net/http"
	"strings"
	"html/template"
)

func init() {
	Route.NewSub(server.NewRoute("/", func(w http.ResponseWriter, r *http.Request) {
		if ok, _ := sessions.IsLoggedIn(w, r); ok {
			http.Redirect(w, r, "/dash/", http.StatusTemporaryRedirect)
		}
		switch strings.ToLower(r.Method) {
		case "get":
			type Page struct {
				Name   string
				Title  string
				Script template.JS
			}
			functions.Render(Page{
				Name:  models.Config.Name,
				Title: "Auth",
			}, w, "auth/auth.html", server.TemplateCache)
		case "post":
			authentication.Signin(w, r)
		}
	}))
}
