package websockets

import (
	"github.com/gorilla/websocket"
	"net/http"
	"login/core/models/server"
)


var (
	Route *server.Route = server.NewSubRouter("/websockets")
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			// This is where you could add your origin checks for security
			return true
		},
	}
)