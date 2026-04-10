package authentication

import (
	"encoding/json"
	"file-tracker-backend/database"
	"net/http"
)

type LoginVerifyRequest struct {
	Email     string `json:"email"`
	SessionID string `json:"sessionId"`
}

func LoginChallengeHandler(w http.ResponseWriter, r *http.Request) {
	var req EmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if req.Email == "" {
		http.Error(w, "missing email", http.StatusBadRequest)
		return
	}

	if !isAllowed(req.Email) {
		http.Error(w, "email not allowed", http.StatusForbidden)
		return
	}

	u, err := getUser(req.Email)
	if err != nil {
		http.Error(w, "user lookup failed", http.StatusInternalServerError)
		return
	}

	options, sessionData, err := webAuthn.BeginLogin(u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// store session in DB-backed store (returns sessionId internally)
	sessionID := loginSessions.set(req.Email, sessionData)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"options":   options,
		"sessionId": sessionID,
	})
}

func LoginVerifyHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginVerifyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.SessionID == "" {
		http.Error(w, "missing email or sessionId", http.StatusBadRequest)
		return
	}

	sessionData, ok := loginSessions.get(req.SessionID)
	if !ok || sessionData == nil {
		http.Error(w, "missing session", http.StatusBadRequest)
		return
	}

	u, err := getUser(req.Email)
	if err != nil {
		http.Error(w, "user lookup failed", http.StatusInternalServerError)
		return
	}

	_, err = webAuthn.FinishLogin(u, *sessionData, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// cleanup session (DB-backed)
	_, _ = database.DB.Exec(`
		DELETE FROM webauthn_sessions
		WHERE id = $1
	`, req.SessionID)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}
