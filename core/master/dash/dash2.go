package dash

import (
	"html/template"
	"login/core/database"
	"login/core/master/api/auth"
	"login/core/models"
	"login/core/models/functions"
	"login/core/models/server"
	"login/core/sessions"
	"net/http"
	"strings"
	"time"
)

func init() {
	Route.NewSub(server.NewRoute("/", func(w http.ResponseWriter, r *http.Request) {
		ok, user := sessions.IsLoggedIn(w, r)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
        avatarPath, _ := database.Container.GetAvatar(user.ID)
		switch strings.ToLower(r.Method) {
		case "get":
			type Page struct {
				Name   string
				Avatar string
				Title  string
				Time     string
				*sessions.Session
				Script template.JS
			}
			formattedTime := time.Now().Format("03:04PM") + " " + user.TMZ
			functions.Render(Page{
				Name: 	 models.Config.Name,
				Avatar:	 "/_assets/media/avatar/"+avatarPath,
				Title: 	"Auth",
				Time: 	formattedTime,
				Session: user,
			}, w, "dash/dash.html", server.TemplateCache)
		case "post":
			authentication.Signin(w, r)
		}
	}))
}
