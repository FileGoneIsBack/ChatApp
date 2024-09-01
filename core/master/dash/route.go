package dash

import "login/core/models/server"

var (
	Route *server.Route = server.NewSubRouter("/dash")
)