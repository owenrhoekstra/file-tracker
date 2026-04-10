package authentication

import (
	"encoding/json"
	"log"
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
		log.Println("CheckEmailHandler: getUser error for", req.Email, ":", err)
		http.Error(w, "user lookup failed", http.StatusInternalServerError)
		return
	}

	hasPasskey := false

	if u != nil {
		creds, err := getCredentialsByUserID(u.ID)
		if err == nil && len(creds) > 0 {
			hasPasskey = true
			log.Println("CheckEmailHandler:", req.Email, "has", len(creds), "credential(s)")
		} else if err != nil {
			log.Println("CheckEmailHandler: getCredentials error for user ID", string(u.ID), ":", err)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"allowed":    true,
		"hasPasskey": hasPasskey,
	})
}
