package database

import (
    "time"
)

type Message struct {
    ID        int
    ChatID    int
    UserID  int
    Content   string
    SentAt    time.Time
}

func (conn *Instance) NewMessage(message *Message) error {
    _, err := Container.conn.Exec("INSERT INTO messages (chat_id, user_id, content) VALUES (?, ?, ?)",
        message.ChatID, message.UserID, message.Content)
    if err != nil {
        return err
    }
    return nil
}

func (conn *Instance) GetMessagesForChat(chatID int) ([]Message, error) {
    rows, err := Container.conn.Query("SELECT id, chat_id, user_id, content, sent_at FROM messages WHERE chat_id = ?", chatID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var messages []Message
    for rows.Next() {
        var msg Message
        err := rows.Scan(&msg.ID, &msg.ChatID, &msg.UserID, &msg.Content, &msg.SentAt)
        if err != nil {
            return nil, err
        }
        messages = append(messages, msg)
    }

    return messages, nil
}

func (conn *Instance) GetNewMessagesForChat(chatID int, lastMessageID int) ([]Message, error) {
    rows, err := Container.conn.Query(`
        SELECT id, chat_id, user_id, content, sent_at 
        FROM messages 
        WHERE chat_id = ? AND id > ? 
        ORDER BY id ASC`, chatID, lastMessageID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var messages []Message
    for rows.Next() {
        var msg Message
        err := rows.Scan(&msg.ID, &msg.ChatID, &msg.UserID, &msg.Content, &msg.SentAt)
        if err != nil {
            return nil, err
        }
        messages = append(messages, msg)
    }

    return messages, nil
}

