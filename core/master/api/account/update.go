package account

import (
    "encoding/json"
    "login/core/database"
    "login/core/models/server"
    "login/core/models/validation"
    "login/core/sessions"
    "net/http"
)

func init() {
    Route.NewSub(server.NewRoute("/updateUser", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodPost {
            // Check if user is logged in
            ok, user := sessions.IsLoggedIn(w, r)
            if !ok {
                http.Error(w, "Unauthorized", http.StatusUnauthorized)
                return
            }

            // Parse JSON request body
            var updatedUser database.User
            err := json.NewDecoder(r.Body).Decode(&updatedUser)
            if err != nil {
                http.Error(w, "Invalid JSON", http.StatusBadRequest)
                return
            }

            // Update the user in the database
            updatedUser.ID = user.ID
			if user.ID != updatedUser.ID {
                http.Error(w, "Unauthorized action", http.StatusUnauthorized)
                return
            }
            err = database.Container.UpdateUser(&updatedUser)
            if err != nil {
                Sanitize.LogMessage("error", "Error updating user: %v", err)
                http.Error(w, "Error updating user", http.StatusInternalServerError)
                return
            }

            // Update the session with the new user information
            sessionToken, err := r.Cookie("session-token")
            if err != nil {
                http.Error(w, "Session token not found", http.StatusUnauthorized)
                return
            }
            err = sessions.UpdateSession(sessionToken.Value, &updatedUser)
            if err != nil {
                Sanitize.LogMessage("error", "Error updating session: %v", err)
                http.Error(w, "Error updating session", http.StatusInternalServerError)
                return
            }

            w.Header().Set("Content-Type", "application/json; charset=utf-8")
            w.WriteHeader(http.StatusOK)
            Sanitize.LogMessage("success", "Updated User: "+user.Username, nil)
        } else {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
    }))
}
