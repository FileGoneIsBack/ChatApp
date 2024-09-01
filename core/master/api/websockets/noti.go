package websockets

import (
    "log"
    "login/core/database"
    "github.com/gorilla/websocket"
    "time"
    "sync"
)

type NotificationManager struct {
    conn            *websocket.Conn
    chatID          int
    lastMessageID   int
    loggedInUserID  int
    done            chan struct{}
    mu              sync.Mutex
}

func NewNotificationManager(conn *websocket.Conn, chatID int, loggedInUserID int) *NotificationManager {
    return &NotificationManager{
        conn:           conn,
        chatID:         chatID,
        loggedInUserID: loggedInUserID,
        done:           make(chan struct{}),
    }
}

func (nm *NotificationManager) StartPolling() {
    go func() {
        for {
            select {
            case <-nm.done:
                return
            default:
                // Fetch new messages since the last message ID
                messages, err := database.Container.GetNewMessagesForChat(nm.chatID, nm.lastMessageID)
                if err != nil {
                    log.Printf("Error fetching new messages: %v", err)
                    nm.handleConnectionError("Error fetching new messages")
                    return
                }

                // If there are new messages, send them
                if len(messages) > 0 {
                    // Update the lastMessageID to the latest one fetched
                    nm.lastMessageID = messages[len(messages)-1].ID

                    // Create JSON response
                    response := generateMessageJSON(messages, nm.loggedInUserID)

                    // Send JSON response over WebSocket
                    if err := nm.conn.WriteJSON(response); err != nil {
                        log.Printf("Error sending new messages: %v", err)
                        nm.handleConnectionError("Error sending new messages")
                        return
                    }
                }

                // Sleep or wait before checking for new messages again
                time.Sleep(10 * time.Second) // Reduced interval
            }
        }
    }()
}

func (nm *NotificationManager) StopPolling() {
    nm.mu.Lock()
    defer nm.mu.Unlock()
    select {
    case <-nm.done:
        // Channel already closed
        return
    default:
        close(nm.done)
    }
}

func (nm *NotificationManager) cleanup() {
    nm.StopPolling()
    nm.conn.Close()
}

func (nm *NotificationManager) handleConnectionError(message string) {
    if nm.conn != nil {
        nm.conn.WriteMessage(websocket.TextMessage, []byte(message))
    }
}

func generateMessageJSON(messages []database.Message, loggedInUserID int) []map[string]interface{} {
    var messageList []map[string]interface{}
    for _, message := range messages {
        messageData := map[string]interface{}{
            "id":      message.ID,
            "content": message.Content,
            "sentAt":  message.SentAt.Format("3:04 PM, Jan 2"),
            "userID":  message.UserID,
            "isMine":  message.UserID == loggedInUserID,
        }
        messageList = append(messageList, messageData)
    }
    return messageList
}
