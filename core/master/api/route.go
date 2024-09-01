package api

import (
	"login/core/master/api/account"
	"login/core/master/api/auth"
	"login/core/master/api/apps"
	"login/core/master/api/functions"
	"login/core/master/api/websockets"
	"login/core/models/server"
)

var (
	Route *server.Route = server.NewSubRouter("/api")
)

func init() {
	Route.NewSubs(
		authentication.Route,
		functions.Route,
		websockets.Route,
		account.Route,
		apps.Route,
	)
}
