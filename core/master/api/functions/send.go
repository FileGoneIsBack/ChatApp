package functions

import (
	"encoding/json"
	"log"
	"login/core/database"
	"login/core/models/server"
	"login/core/models/validation" // Import your validation package
	"login/core/sessions"
	"net/http"
	"strconv"
)

func init() {
	Route.NewSub(server.NewRoute("/send", func(w http.ResponseWriter, r *http.Request) {
		// Check if the user is logged in
		ok, user := sessions.IsLoggedIn(w, r)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		cookie, err := r.Cookie("chatID")
		if err != nil {
			http.Error(w, "ChatID cookie not found", http.StatusBadRequest)
			return
		}
		chatID, _ := strconv.Atoi(cookie.Value)
		// Parse the request body
		var requestData struct {
			Message string `json:"message"`
		}
		if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
			log.Printf("Error decoding request body: %v", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Validate the message
		requestData.Message = Sanitize.SanitizeInput(requestData.Message)
		if requestData.Message == "" || !Sanitize.ValidateMessage(requestData.Message) {
			http.Error(w, "Invalid message", http.StatusBadRequest)
			return
		}

		// Insert the new message into the database
		err = database.Container.NewMessage(&database.Message{
			ChatID:  chatID,
			UserID:  user.ID,
			Content: requestData.Message,
		})
		if err != nil {
			log.Printf("Error creating new message: %v", err)
			http.Error(w, "Failed to create message", http.StatusInternalServerError)
			return
		}
		// Respond with success
		w.WriteHeader(http.StatusOK)
	}))
}
