package functions

import (
	"encoding/json"
	"log"
	"login/core/database"
	"login/core/models/server"
	"login/core/models/validation"
	"login/core/sessions"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func init() {
	Route.NewSub(server.NewRoute("/chat/{user_id}", func(w http.ResponseWriter, r *http.Request) {
		ok, user := sessions.IsLoggedIn(w, r)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Extract the user ID from the URL path
		vars := mux.Vars(r)
		peopleID := vars["user_id"]
		userID := user.ID
		peopleIDInt, err := strconv.Atoi(peopleID)
		if err != nil {
			log.Printf("Invalid user_id in URL: %v", err)
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		if r.Method == http.MethodPost {
			// Example: Fetch chat history from database based on userID
			chats, messages, err := database.Container.GetChatAndMessages(userID, peopleIDInt)
			if err != nil {
				Sanitize.LogMessage("success", "Updated User: "+user.Username, err)
				http.Error(w, "Error fetching chats", http.StatusInternalServerError)
				return
			}

			response := map[string]interface{}{
				"chats":    chats,
				"messages": messages,
			}

			// Set the response header and encode the response as JSON
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			if err := json.NewEncoder(w).Encode(response); err != nil {
				http.Error(w, "Error encoding response to JSON", http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))
}
