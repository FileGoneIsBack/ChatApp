package account

import (
	"login/core/database"
	"login/core/models/server"
	"login/core/sessions"
	"net/http"
)

func init() {
    Route.NewSub(server.NewRoute("/contacts", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodPost {
            // Check if user is logged in
            ok, user := sessions.IsLoggedIn(w, r)
            if !ok {
                http.Error(w, "Unauthorized", http.StatusUnauthorized)
                return
            }

            if err := r.ParseForm(); err != nil {
                http.Error(w, "Invalid form data", http.StatusBadRequest)
                return
            }

            contactIDStr := r.FormValue("contactID")
            if contactIDStr == "" {
                http.Error(w, "ContactID is required", http.StatusBadRequest)
                return
            }

            database.Container.SendReq(user.Username, contactIDStr)

            w.Write([]byte("Friend request sent successfully"))
        } else {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
    }))
}
