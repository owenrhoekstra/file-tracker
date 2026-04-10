package authentication

import (
	"github.com/go-webauthn/webauthn/webauthn"
)

var webAuthn *webauthn.WebAuthn

func InitWebAuthn() {
	var err error

	webAuthn, err = webauthn.New(&webauthn.Config{
		RPDisplayName: "FileLogix",
		RPID:          "orh-home-server.tailac3f56.ts.net",
		RPOrigins:     []string{"https://orh-home-server.tailac3f56.ts.net"},
	})
	if err != nil {
		panic(err)
	}
}
