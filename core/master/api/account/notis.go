package account

import (
    "login/core/models/server"
    "login/core/sessions"
    "net/http"
	"encoding/json"
    "login/core/database"
)

func init() {
    Route.NewSub(server.NewRoute("/notifications", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodPost {
            // Check if user is logged in
            ok, user := sessions.IsLoggedIn(w, r)
            if !ok {
                http.Error(w, "Unauthorized", http.StatusUnauthorized)
                return
            }

			contactRequests, err := database.Container.GetPending(user.ID)
			if err != nil {
				http.Error(w, "Failed to retrieve contact requests", http.StatusInternalServerError)
				return
			}
            // Respond with a success message
            w.Write([]byte("Friend request sent successfully"))
			json.NewEncoder(w).Encode(contactRequests)
        } else {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
    }))
}


