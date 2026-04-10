package authentication

import (
	"file-tracker-backend/database"
)

var allowedEmails = map[string]bool{
	"owenhoekstra@icloud.com": true,
}

// simple guard
func isAllowed(email string) bool {
	return allowedEmails[email]
}

// USER LOADING (DB-backed)
func getUser(email string) (*User, error) {
	u := &User{Email: email}

	err := database.DB.QueryRow(`
		SELECT id, email
		FROM users
		WHERE email = $1
	`, email).Scan(&u.ID, &u.Email)

	if err == nil {
		return u, nil
	}

	// not found → create
	u.ID = []byte(email)

	_, err = database.DB.Exec(`
		INSERT INTO users (id, email)
		VALUES ($1, $2)
	`, u.ID, u.Email)

	if err != nil {
		return nil, err
	}

	return u, nil
}
