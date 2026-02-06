package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Cinnamoon-dev/blue-gopher/internal/customerrors"
)

func RespondJSON(w http.ResponseWriter, status int, payload any) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")

	if payload != nil {
		json.NewEncoder(w).Encode(payload)
	}
}

func RespondError(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	var httpError *customerrors.HTTPError
	if errors.As(err, &httpError) {
		RespondJSON(w, httpError.Status, map[string]string{"error": httpError.Message})
		return
	}

	RespondJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
}
