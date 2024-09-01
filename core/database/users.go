package database

import (
    "database/sql"
    "errors"
    "golang.org/x/crypto/bcrypt"
)

type User struct {
    ID              int
    Username, Password,
    Email, LastName, 
    FirstName, DOB, Bio,
    Website, TMZ,
    Avatar          string
    Status          bool
}

var (
    ErrDuplicateUser   = errors.New("duplicate user")
    ErrUserNotFound    = errors.New("user couldn't be found in the database")
    ErrInvalidPassword = errors.New("invalid password")
)

// NewUser inserts a new user into the database
func (conn *Instance) NewUser(user *User) error {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }

    _, err = Container.conn.Exec("INSERT INTO users (username, password_hash, email) VALUES (?, ?, ?)", user.Username, hashedPassword, user.Email)
    if err != nil {
        return err
    }

    return nil
}

// UpdateUserStatus updates the status of a user in the database
func (conn *Instance) UpdateUserStatus(userID int, status bool) error {
    _, err := Container.conn.Exec("UPDATE users SET status = ? WHERE id = ?", status, userID)
    if err != nil {
        return err
    }
    return nil
}

// AuthenticateUser checks if the username and password match a user in the database
func (conn *Instance) AuthenticateUser(username, password string) (*User, error) {
    var hashedPassword string
    var user User

    err := Container.conn.QueryRow(`
        SELECT id, username, password_hash, email, lastName, firstName, DOB, website, status 
        FROM users 
        WHERE username = ?`, 
        username,
    ).Scan(
        &user.ID, 
        &user.Username, 
        &hashedPassword, 
        &user.Email, 
        &user.LastName, 
        &user.FirstName, 
        &user.DOB, 
        &user.Website, 
        &user.Status,
    )
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, ErrUserNotFound
        }
        return nil, err
    }

    // Verify the password hash
    err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
    if err != nil {
        return nil, ErrInvalidPassword
    }

    // Clear the password field before returning
    user.Password = ""
    return &user, nil
}



// UpdateUser updates a user's information in the database.
func (conn *Instance) UpdateUser(user *User) error {
    _, err := Container.conn.Exec(`
        UPDATE users
        SET firstName = ?, lastName = ?, email = ?, DOB = ?, website = ?
        WHERE id = ?
    `, user.FirstName, user.LastName, user.Email, user.DOB, user.Website, user.ID)

    if err != nil {
        return err
    }

    return nil
}

// UpdateUser updates a user's information in the database.
func (conn *Instance) UpdateAvatar(path string, id int) error {
    _, err := Container.conn.Exec(`
        UPDATE users
        SET avatar = ?
        WHERE id = ?
    `, path, id)

    if err != nil {
        return err
    }

    return nil
}

// UpdateUser updates a user's information in the database.
func (conn *Instance) GetAvatar(id int) (string, error) {
    var avatar string
    err := conn.conn.QueryRow(`
        SELECT avatar 
        FROM users
        WHERE id = ?
    `, id).Scan(&avatar)

    if err != nil {
        return "", err
    }

    return avatar, nil
}

