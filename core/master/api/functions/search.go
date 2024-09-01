package functions

import (
	"log"
	"login/core/database"
	"login/core/models/server"
	"net/http"
	"encoding/json"
)

func init() {
    Route.NewSub(server.NewRoute("/search", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodPost {
            username := r.FormValue("username")
            log.Printf("search username: %s", username)

            users, err := database.Container.GetUsersByUsernamePattern(username)
            if err != nil {
                log.Printf("Error searching users: %v", err)
                http.Error(w, "Error searching users", http.StatusInternalServerError)
                return
            }

            w.Header().Set("Content-Type", "application/json; charset=utf-8")
            if err := json.NewEncoder(w).Encode(users); err != nil {
                log.Printf("Error encoding users to JSON: %v", err)
                http.Error(w, "Error encoding users to JSON", http.StatusInternalServerError)
            }
        } else {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
    }))
}
