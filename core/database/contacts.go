package database

import (

)

type Contacts struct {
    ID        int
    UserID    string
    ContactID int
    Status    string
}

func (conn *Instance) SendReq(UserID, ContactID string) error {
    query := `INSERT INTO contacts (user_id, contact_id, status) VALUES (?, ?, 'pending')`
    
    _, err := Container.conn.Exec(query, UserID, ContactID)
    if err != nil {
        return err
    }

    return nil
}

// GetPending retrieves contacts with a status of "pending" for a given ContactID
func (conn *Instance) GetPending(ContactID int) (*Contacts, error) {
    var pendingContacts *Contacts
    query := `SELECT id, user_id, contact_id, status FROM contacts WHERE contact_id = ? AND status = 'pending'`
    
    _, err := Container.conn.Exec(query, ContactID)
    if err != nil {
        return nil, err
    }

    return pendingContacts, nil
}
