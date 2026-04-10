package authentication

import (
	"file-tracker-backend/database"

	"github.com/go-webauthn/webauthn/webauthn"
)

// SAVE credential after registration
func saveCredential(userID []byte, cred *webauthn.Credential) error {
	_, err := database.DB.Exec(`
		INSERT INTO credentials (
			user_id,
			credential_id,
			public_key,
			attestation_type,
			transports,
			sign_count
		) VALUES ($1,$2,$3,$4,$5,$6)
	`,
		userID,
		cred.ID,
		cred.PublicKey,
		cred.AttestationType,
		cred.Transport,
		cred.Authenticator.SignCount,
	)

	return err
}

// GET credentials for login (IMPORTANT for WebAuthn)
func getCredentialsByUserID(userID []byte) ([]webauthn.Credential, error) {
	rows, err := database.DB.Query(`
		SELECT credential_id, public_key, attestation_type, transports, sign_count
		FROM credentials
		WHERE user_id = $1
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var creds []webauthn.Credential

	for rows.Next() {
		var cred webauthn.Credential

		err := rows.Scan(
			&cred.ID,
			&cred.PublicKey,
			&cred.AttestationType,
			&cred.Transport,
			&cred.Authenticator.SignCount,
		)
		if err != nil {
			return nil, err
		}

		creds = append(creds, cred)
	}

	return creds, nil
}
