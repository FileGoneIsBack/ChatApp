package database

import (
    "database/sql"
    "errors"
    "time"
)

type Chat struct {
    ID        int
    Name      string
    UserID    int
    PeopleID  int
    CreatedAt time.Time
}

var (
    ErrChatNotFound    = errors.New("chat not found in the database")
    ErrMessageNotFound = errors.New("message not found in the database")
)

// CreateChat inserts a new chat into the database
func (conn *Instance) CreateChat(chat *Chat) error {
    query := `INSERT INTO chats (name, user_id, people_id, created_at) VALUES (?, ?, ?, ?)`
    
    result, err := Container.conn.Exec(query, chat.Name, chat.UserID, chat.PeopleID, chat.CreatedAt)
    if err != nil {
        return err
    }

    // Retrieve the last inserted ID
    chatID, err := result.LastInsertId()
    if err != nil {
        return err
    }

    chat.ID = int(chatID)
    return nil
}


// GetChatByID retrieves a chat by its ID from the database
func (conn *Instance) GetChatByID(chatID int) (*Chat, []Message, error) {
    // Retrieve chat details
    var chat Chat
    err := Container.conn.QueryRow("SELECT id, name, user_id, people_id, created_at FROM chats WHERE id = ?", chatID).Scan(&chat.ID, &chat.Name, &chat.UserID, &chat.PeopleID, &chat.CreatedAt)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil, ErrChatNotFound
        }
        return nil, nil, err
    }

    // Retrieve messages for the chat
    messages, err := conn.GetMessagesForChat(chat.ID)
    if err != nil {
        return nil, nil, err
    }

    return &chat, messages, nil
}


// GetChatAndMessages retrieves chat and messages between two users
func (conn *Instance) GetChatAndMessages(userID, peopleID int) (*Chat, []Message, error) {
    username, err := conn.GetUsernameByID(peopleID)
    if err != nil {
        return nil, nil, err
    }

    chat, err := conn.GetChatByUserIDs(userID, peopleID)
    if err != nil {
        if errors.Is(err, ErrChatNotFound) {
            newChat := Chat{
                Name:      username,
                UserID:    userID,
                PeopleID:  peopleID,
                CreatedAt: time.Now(),
            }
            err := conn.CreateChat(&newChat)
            if err != nil {
                return nil, nil, err
            }
            return &newChat, nil, nil
        }
        return nil, nil, err
    }


    messages, err := conn.GetMessagesForChat(chat.ID)
    if err != nil {
        return nil, nil, err
    }

    return chat, messages, nil
}


func (conn *Instance) GetChatByUserIDs(userID, peopleID int) (*Chat, error) {
    // Query for the chat where either (userID, peopleID) or (peopleID, userID) match
    query := `
        SELECT * FROM chats
        WHERE (user_id = ? AND people_id = ?)
        OR (user_id = ? AND people_id = ?)
        LIMIT 1
    `
    var chat Chat
    err := conn.conn.QueryRow(query, userID, peopleID, peopleID, userID).Scan(&chat.ID, &chat.Name, &chat.UserID, &chat.PeopleID, &chat.CreatedAt)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, ErrChatNotFound
        }
        return nil, err
    }
    return &chat, nil
}

// GetRecentChats retrieves the last 5 chats involving the specified user, handling both sender and receiver cases.
func (conn *Instance) GetRecentChats(userID int) ([]Chat, error) {
    query := `
        SELECT id, name, 
               CASE WHEN user_id = ? THEN people_id ELSE user_id END AS chatPartnerID,
               created_at
        FROM chats
        WHERE (user_id = ? OR people_id = ?)
        AND user_id != people_id
        ORDER BY created_at DESC
        LIMIT 5
    `
    
    rows, err := conn.conn.Query(query, userID, userID, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var chats []Chat
    for rows.Next() {
        var chat Chat
        if err := rows.Scan(&chat.ID, &chat.Name, &chat.PeopleID, &chat.CreatedAt); err != nil {
            return nil, err
        }
        chats = append(chats, chat)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return chats, nil
}
