package sessions

import (
	"login/core/database"
	"net/http"
	"time"

	"github.com/google/uuid"
)

var (
	Sessions = map[string]Session{}
)

// Session is used to store the user & expiry
type Session struct {
	*database.User
	Expiry time.Time
}

func (s Session) IsExpired() bool {
	return s.Expiry.Before(time.Now())
}

func Count() int {
    return len(Sessions)
}

func CreateSession(user *database.User) *http.Cookie {
    sessionToken := uuid.NewString()
    expiresAt := time.Now().Add(30 * time.Minute)

    Sessions[sessionToken] = Session{
        User:   user,
        Expiry: expiresAt,
    }

    err := database.Container.UpdateUserStatus(user.ID, true) 
    if err != nil {
        return nil
    }

    // Create the session cookie
    cookie := &http.Cookie{
        Name:     "session-token",
        Value:    sessionToken,
        Expires:  expiresAt,
        Path:     "/",
        HttpOnly: true,
        Secure:   true,
    }

    return cookie
}

// UpdateSession updates the user's information in the session.
func UpdateSession(sessionToken string, updatedUser *database.User) error {
    // Find the session by the session token
    session, exists := Sessions[sessionToken]
    if !exists {
        return nil
    }

    // Update the user information in the session
    session.User = updatedUser

    // Update the session expiry time if needed (optional)
    session.Expiry = time.Now().Add(30 * time.Minute)

    // Save the updated session back to the map
    Sessions[sessionToken] = session

    return nil
}