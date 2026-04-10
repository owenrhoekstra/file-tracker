package authentication

import (
	"crypto/rand"
	"database/sql"
	"file-tracker-backend/database"
)

var allowedEmails = map[string]bool{
	"owenhoekstra@icloud.com": true,
}

// simple guard
func isAllowed(email string) bool {
	return allowedEmails[email]
}

// Generate a proper user ID (UUID)
func generateUserID() []byte {
	id := make([]byte, 16)
	_, _ = rand.Read(id)
	return id
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

	// not found → create with generated ID
	u.ID = generateUserID()

	_, err = database.DB.Exec(`
		INSERT INTO users (id, email)
		VALUES ($1, $2)
		ON CONFLICT (email) DO NOTHING
	`, u.ID, u.Email)

	if err != nil {
		return nil, err
	}

	// After insert (or conflict), fetch the actual user to ensure consistency
	err = database.DB.QueryRow(`
		SELECT id, email
		FROM users
		WHERE email = $1
	`, email).Scan(&u.ID, &u.Email)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	return u, nil
}
