package landing

import "login/core/models/server"

var (
	Route *server.Route = server.NewSubRouter("")
)
