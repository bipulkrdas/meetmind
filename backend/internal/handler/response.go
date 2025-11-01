package handler

import (
    "encoding/json"
    "net/http"
	"errors"
    "github.com/google/uuid"
)

func respondWithError(w http.ResponseWriter, code int, message string) {
    respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
    response, _ := json.Marshal(payload)

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    w.Write(response)
}

func getUserIDFromContext(r *http.Request) (uuid.UUID, error) {
    userIDStr, ok := r.Context().Value("userID").(string)
    if !ok {
        return uuid.Nil, errors.New("userID not found in context")
    }
    return uuid.Parse(userIDStr)
}
