package sessions

import (
	"net/http"
	"login/core/database"
)

func IsLoggedIn(w http.ResponseWriter, r *http.Request) (bool, *Session) {
	c, err := r.Cookie("session-token")
	if err != nil {
		if err == http.ErrNoCookie {
			if r.URL.Path != "/auth/" {
				http.Redirect(w, r, "/auth/", http.StatusTemporaryRedirect)
			}
			return false, nil
		}
		if r.URL.Path != "/auth/" {
			http.Redirect(w, r, "/auth/", http.StatusTemporaryRedirect)
		}
		return false, nil
	}
	sessionToken := c.Value

	userSession, exists := Sessions[sessionToken]
	if !exists {
		if r.URL.Path != "/auth/" {
			http.Redirect(w, r, "/auth/", http.StatusTemporaryRedirect)
		}
		return false, nil
	}

	if userSession.IsExpired() {
	database.Container.UpdateUserStatus(userSession.ID, false)
		delete(Sessions, sessionToken)
		if r.URL.Path != "/auth/" {
			http.Redirect(w, r, "/auth/", http.StatusTemporaryRedirect)
		}
		return false, nil
	}
	return true, &userSession
}
