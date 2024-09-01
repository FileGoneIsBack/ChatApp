package websockets

import (
    "log"
    "login/core/database"
    "login/core/models/server"
    "login/core/sessions"
    "net/http"
    "strconv"
    "github.com/gorilla/websocket"
)

func init() {
    Route.NewSub(server.NewRoute("/messages", func(w http.ResponseWriter, r *http.Request) {
        // Check if user is logged in
        ok, user := sessions.IsLoggedIn(w, r)
        if !ok {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        // Upgrade the HTTP connection to a WebSocket connection
        conn, err := upgrader.Upgrade(w, r, nil)
        if err != nil {
            log.Printf("Failed to upgrade connection: %v", err)
            return
        }
        defer conn.Close()

        // Read chatID from the cookie
        chatIDCookie, err := r.Cookie("chatID")
        if err != nil {
            if err == http.ErrNoCookie {
                conn.WriteMessage(websocket.TextMessage, []byte("chatID cookie not found"))
                return
            }
            conn.WriteMessage(websocket.TextMessage, []byte("Invalid chat ID"))
            return
        }

        // Convert chatID from string to int
        chatID, err := strconv.Atoi(chatIDCookie.Value)
        if err != nil {
            conn.WriteMessage(websocket.TextMessage, []byte("Invalid chat ID"))
            return
        }

        // Load existing messages for the chat
        messages, err := database.Container.GetMessagesForChat(chatID)
        if err != nil {
            log.Printf("Error fetching existing messages: %v", err)
            conn.WriteMessage(websocket.TextMessage, []byte("Error fetching existing messages"))
            return
        }

        // Send existing messages to the client
        response := generateMessageJSON(messages, user.ID)
        if err := conn.WriteJSON(response); err != nil {
            log.Printf("Error sending existing messages: %v", err)
            return
        }

        // Initialize and start NotificationManager to listen for new messages
        nm := NewNotificationManager(conn, chatID, user.ID)
        nm.StartPolling()

        // Ensure you stop polling when the WebSocket connection is closed
        defer nm.StopPolling()
    }))
}

