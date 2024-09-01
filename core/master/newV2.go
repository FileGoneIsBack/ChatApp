package master

import (
	"login/core/master/authentication"
	"login/core/master/dash"
	"login/core/master/api"
	"login/core/master/landing"
	"login/core/models"
	"login/core/models/server"
	"login/core/models/functions"
	"net/http"
	"login/core/models/validation"
)

var (
	Service *server.Server = server.NewServer(&server.Config{
		Addr:   "0.0.0.0:80",
		Secure: models.Config.Secure,
		Cert:   models.Config.Cert,
		Key:    models.Config.Key,
	})
	Route  *server.Route   = server.NewSubRouter("")
	Assets *server.Handler = server.NewHandler("/_assets/", http.StripPrefix("/_assets", http.FileServer(http.Dir("assets/branding"))))
)

func NewV2() {
    server.TemplateCache = functions.NewTemplateCache()
    err := server.TemplateCache.LoadTemplates("assets/html", "assets/html/render")
    if err != nil {
        Sanitize.LogMessage("error", "Failed To Load Templs", err)
    }
	server.TemplateCache.ListTemplates()

    // Set up routes and handlers
    Service.AddRoutes(
        dash.Route,
        authentication.Route,
        landing.Route,
        api.Route,
    )
    Service.AddHandler(Assets)
    Service.AddRoute(server.NewRoute("", nil))
    if err := Service.ListenAndServe(); err != nil {
        panic(err)
    }
}

