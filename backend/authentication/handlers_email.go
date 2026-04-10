package authentication

import (
	"encoding/json"
	"net/http"
)

func CheckEmailHandler(w http.ResponseWriter, r *http.Request) {
	var req EmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if !isAllowed(req.Email) {
		http.Error(w, "email not allowed", http.StatusForbidden)
		return
	}

	u, err := getUser(req.Email)
	if err != nil {
		u = nil
	}

	hasPasskey := false

	if u != nil {
		creds, err := getCredentialsByUserID(u.ID)
		if err == nil && len(creds) > 0 {
			hasPasskey = true
		}
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"allowed":    true,
		"hasPasskey": hasPasskey,
	})
}
