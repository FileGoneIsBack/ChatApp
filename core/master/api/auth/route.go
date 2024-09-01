package authentication

import "login/core/models/server"

var (
	Route *server.Route = server.NewSubRouter("/auth")
)