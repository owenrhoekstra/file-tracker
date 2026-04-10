package authentication

import (
	"encoding/json"
	"net/http"
)

type EmailRequest struct {
	Email string `json:"email"`
}

func decodeEmail(r *http.Request) (string, error) {
	var req EmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return "", err
	}
	return req.Email, nil
}
