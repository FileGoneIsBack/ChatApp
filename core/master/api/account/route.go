package account

import "login/core/models/server"

var (
	Route *server.Route = server.NewSubRouter("/account")
)