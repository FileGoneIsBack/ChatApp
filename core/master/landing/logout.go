package landing

import (
	"net/http"
	"time"
	"login/core/database"
	"login/core/sessions"
	"login/core/models/server"
)

func init() {
	Route.NewSub(server.NewRoute("/logout", func(w http.ResponseWriter, r *http.Request) {
		// Check if user is logged in
		ok, user := sessions.IsLoggedIn(w, r)
		if !ok {
			// If not logged in, redirect to home
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		// Retrieve session token cookie
		c, err := r.Cookie("session-token")
		if err != nil {
			if err == http.ErrNoCookie {
				// Handle missing cookie case
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			// Handle other cookie errors
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Update the user's status to inactive
		err = database.Container.UpdateUserStatus(user.ID, false)
		if err != nil {
			// Log the error or handle it appropriately
			http.Error(w, "Failed to update user status", http.StatusInternalServerError)
			return
		}

		// Remove the session
		delete(sessions.Sessions, c.Value)

		// Clear the session token cookie
		http.SetCookie(w, &http.Cookie{
			Name:    "session-token",
			Value:   "",
			Expires: time.Now().Add(-time.Hour), // Set expiration to the past to ensure it is cleared
			Path:    "/",
			HttpOnly: true,
			Secure:   true, // Ensure the secure flag matches your HTTPS configuration
		})

		// Redirect to home after logging out
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}))
}
