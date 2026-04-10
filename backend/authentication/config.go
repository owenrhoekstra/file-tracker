package authentication

import (
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/joho/godotenv"
	"log"
	"os"
)

var webAuthn *webauthn.WebAuthn

func InitWebAuthn() {
	// Load .env file if it exists (ignore error if file doesn't exist)
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Get RPID from environment variable, fallback to default
	rpid := os.Getenv("WEBAUTHN_RPID")
	if rpid == "" {
		rpid = "orh-home-server.tailac3f56.ts.net" // default for dev
	}

	// Get RP Origins from environment variable, fallback to default
	rpOrigin := os.Getenv("WEBAUTHN_RP_ORIGIN")
	if rpOrigin == "" {
		rpOrigin = "https://orh-home-server.tailac3f56.ts.net" // default for dev
	}

	log.Printf("WebAuthn Config - RPID: %s, RPOrigin: %s", rpid, rpOrigin)

	webAuthn, err = webauthn.New(&webauthn.Config{
		RPDisplayName: "FileLogix",
		RPID:          rpid,
		RPOrigins:     []string{rpOrigin},
	})
	if err != nil {
		panic(err)
	}
}
