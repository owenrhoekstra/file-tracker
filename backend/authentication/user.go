package authentication

import (
	"github.com/go-webauthn/webauthn/webauthn"
)

type User struct {
	ID    []byte
	Email string
}

func (u User) WebAuthnID() []byte {
	return u.ID
}

func (u User) WebAuthnName() string {
	return u.Email
}

func (u User) WebAuthnDisplayName() string {
	return u.Email
}

func (u User) WebAuthnCredentials() []webauthn.Credential {
	creds, err := getCredentialsByUserID(u.ID)
	if err != nil {
		return []webauthn.Credential{}
	}
	return creds
}

func (u User) WebAuthnIcon() string {
	return ""
}
