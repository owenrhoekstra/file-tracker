package authentication

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"file-tracker-backend/database"
	"log"
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
	log.Println("Generated new UUID for user:", hex.EncodeToString(id))
	return id
}

// USER LOADING (DB-backed)
func getUser(email string) (*User, error) {
	u := &User{Email: email}

	log.Println("Looking up user by email:", email)

	err := database.DB.QueryRow(`
		SELECT id, email
		FROM users
		WHERE email = $1
	`, email).Scan(&u.ID, &u.Email)

	if err == nil {
		log.Println("User found in database, ID:", hex.EncodeToString(u.ID))
		return u, nil
	}

	// not found → create with generated ID
	log.Println("User not found, creating new user")
	u.ID = generateUserID()

	log.Println("Inserting user with ID:", hex.EncodeToString(u.ID), "email:", email)

	_, err = database.DB.Exec(`
		INSERT INTO users (id, email)
		VALUES ($1, $2)
		ON CONFLICT (email) DO NOTHING
	`, u.ID, u.Email)

	if err != nil {
		log.Println("Error inserting user:", err)
		return nil, err
	}

	log.Println("User insert complete, fetching from database")

	// After insert (or conflict), fetch the actual user to ensure consistency
	err = database.DB.QueryRow(`
		SELECT id, email
		FROM users
		WHERE email = $1
	`, email).Scan(&u.ID, &u.Email)

	if err != nil {
		log.Println("Error fetching user after insert:", err)
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	log.Println("User fetched from database, confirmed ID:", hex.EncodeToString(u.ID))

	return u, nil
}
