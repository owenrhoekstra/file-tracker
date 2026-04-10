package main

import (
	"file-tracker-backend/authentication"
	"file-tracker-backend/database"
	"net/http"
)

func main() {
	authentication.InitWebAuthn()
	database.Init()

	http.HandleFunc("/api/auth/check-email", authentication.CheckEmailHandler)
	http.HandleFunc("/api/auth/passkey/register-challenge", authentication.RegisterChallengeHandler)
	http.HandleFunc("/api/auth/passkey/register-verify", authentication.RegisterVerifyHandler)
	http.HandleFunc("/api/auth/passkey/login-challenge", authentication.LoginChallengeHandler)
	http.HandleFunc("/api/auth/passkey/login-verify", authentication.LoginVerifyHandler)
	http.ListenAndServe(":8000", nil)
}
