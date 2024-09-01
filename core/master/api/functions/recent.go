package functions

import (
    "encoding/json"
    "log"
    "net/http"
    "login/core/database"
    "login/core/models/server"
    "login/core/sessions"
)

func init() {
    Route.NewSub(server.NewRoute("/recent", func(w http.ResponseWriter, r *http.Request) {
        // Check if user is logged in
        ok, user := sessions.IsLoggedIn(w, r)
        if !ok {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
    
        // Fetch recent chats for the logged-in user
        chats, err := database.Container.GetRecentChats(user.ID)
        if err != nil {
            log.Printf("Error fetching recent chats for user %d: %v", user.ID, err)
            http.Error(w, "Error fetching recent chats", http.StatusInternalServerError)
            return
        }
    
        if len(chats) == 0 {
            // Return an empty JSON array if no chats are found
            w.Header().Set("Content-Type", "application/json; charset=utf-8")
            w.Write([]byte("[]"))
            return
        }
    
        // Create JSON response
        response := generateChatsJSON(chats)
    
        // Return the JSON response
        w.Header().Set("Content-Type", "application/json; charset=utf-8")
        if err := json.NewEncoder(w).Encode(response); err != nil {
            log.Printf("Error encoding response to JSON: %v", err)
            http.Error(w, "Error encoding response to JSON", http.StatusInternalServerError)
            return
        }
    }))
    
}

func generateChatsJSON(chats []database.Chat) []map[string]interface{} {
    var chatList []map[string]interface{}
    for _, chat := range chats {
        chatData := map[string]interface{}{
            "chatID":    chat.ID,
            "name":      chat.Name,
            "peopleID":  chat.PeopleID,
            "createdAt": chat.CreatedAt.Format("3:04 PM, Jan 2"),
        }
        chatList = append(chatList, chatData)
    }
    return chatList
}
