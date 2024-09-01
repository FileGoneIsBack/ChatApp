package database

import (
    "database/sql"
	"log"
)

// GetUsers retrieves all users from the database
func (conn *Instance) GetUsers() ([]User, error) {
    rows, err := Container.conn.Query("SELECT id, username, email FROM users")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var users []User
    for rows.Next() {
        var user User
        err := rows.Scan(&user.ID, &user.Username, &user.Email)
        if err != nil {
            return nil, err
        }
        users = append(users, user)
    }

    return users, nil
}

// GetUserByUsername retrieves a user by their username from the database
func (conn *Instance) GetUsersByUsername(username string) (*User, error) {
    var user User
    err := Container.conn.QueryRow(`
        SELECT id, username, email, lastName, firstName, DOB, website, status, bio, tmz
        FROM users 
        WHERE username = ?`, 
        username,
    ).Scan(
        &user.ID, 
        &user.Username, 
        &user.Email, 
        &user.LastName, 
        &user.FirstName, 
        &user.DOB, 
        &user.Website, 
        &user.Status,
        &user.Bio,
        &user.TMZ,
    )
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, ErrUserNotFound
        }
        return nil, err
    }

    return &user, nil
}
// GetUserByUsername retrieves a user by their username from the database
func (conn *Instance) GetUsernameByID(ID int) (string, error) {
    var username string
    err := Container.conn.QueryRow("SELECT username FROM users WHERE id = ?", ID).Scan(&username)
    if err != nil {
        if err == sql.ErrNoRows {
            return "", ErrUserNotFound
        }
        return "", err
    }

    return username, nil
}
// GetUsersByUsernamePattern retrieves users whose usernames match the given pattern
func (conn *Instance) GetUsersByUsernamePattern(usernamePattern string) ([]User, error) {
    query := "SELECT id, username, email, avatar FROM users WHERE username LIKE ?"
    rows, err := Container.conn.Query(query, usernamePattern+"%")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var users []User
    for rows.Next() {
        var user User
        err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Avatar)
        if err != nil {
            return nil, err
        }
        users = append(users, user)
    }

    return users, nil
}


// GetUserByID retrieves a user by their ID from the database
func (conn *Instance) GetUserByID(userID int) (*User, error) {
    var user User
    err := Container.conn.QueryRow(`
        SELECT id, username, email, lastName, firstName, DOB, website, status, bio, tmz, avatar
        FROM users 
        WHERE id = ?`, 
        userID,
    ).Scan(
        &user.ID, 
        &user.Username, 
        &user.Email, 
        &user.LastName, 
        &user.FirstName, 
        &user.DOB, 
        &user.Website, 
        &user.Status,
        &user.Bio,
        &user.TMZ,
        &user.Avatar,
    )
    if err != nil {
        if err == sql.ErrNoRows {
            log.Printf("User not found for ID: %d", userID)
            return nil, ErrUserNotFound
        }
        log.Printf("Error fetching user for ID %d: %v", userID, err)
        return nil, err
    }
    return &user, nil
}
