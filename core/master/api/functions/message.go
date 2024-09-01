package functions

import (
	"encoding/json"
	"log"
	"login/core/database"
	"login/core/models/server"
	"login/core/sessions"
	"net/http"
	"strconv"
)

func init() {
	Route.NewSub(server.NewRoute("/messages", func(w http.ResponseWriter, r *http.Request) {
		// Check if user is logged in
		ok, user := sessions.IsLoggedIn(w, r)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Read chatID from the cookie
		chatIDCookie, err := r.Cookie("chatID")
		if err != nil {
			if err == http.ErrNoCookie {
				http.Error(w, "chatID cookie not found", http.StatusBadRequest)
				return
			}
			http.Error(w, "Invalid chat ID", http.StatusBadRequest)
			return
		}

		// Convert chatID from string to int
		chatID, err := strconv.Atoi(chatIDCookie.Value)
		if err != nil {
			http.Error(w, "Invalid chat ID", http.StatusBadRequest)
			return
		}

		// Fetch just the new message
		messages, err := database.Container.GetMessagesForChat(chatID)
		if err != nil {
			log.Printf("Error fetching new message for user %d: %v", user.ID, err)
			http.Error(w, "Error fetching new message", http.StatusInternalServerError)
			return
		}

		// Create JSON response
		response := generateMessageJSON(messages, user.ID)

		// Return the JSON response
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("Error encoding response to JSON: %v", err)
			http.Error(w, "Error encoding response to JSON", http.StatusInternalServerError)
			return
		}
	}))
}

func generateMessageJSON(messages []database.Message, loggedInUserID int) []map[string]interface{} {
	var messageList []map[string]interface{}
	for _, message := range messages {
		messageData := map[string]interface{}{
			"id":       message.ID,
			"content":  message.Content,
			"sentAt":   message.SentAt.Format("3:04 PM, Jan 2"),
			"userID":   message.UserID,
			"isMine":   message.UserID == loggedInUserID,
		}
		messageList = append(messageList, messageData)
	}
	return messageList
}
