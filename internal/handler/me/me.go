package me

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	dm "nyeh-back/internal/domain/me"
	"nyeh-back/internal/middleware"
)

type MeHandler struct{}

func NewMeHandler() *MeHandler {
	return &MeHandler{}
}

// AdminHandler godoc
//
//	@Summary	Admin Check
//	@Accept		json
//	@Produce	json
//	@Router		/me [get]
func (h *MeHandler) GetMeHandler(w http.ResponseWriter, r *http.Request) {

	email, ok := r.Context().Value(middleware.UserEmailKey).(string)

	if !ok {
		http.Error(w, "Unauthorized: Could not find user email", http.StatusUnauthorized)
		return
	}

	response := dm.NyehResponse{
		Status:  "ok",
		Message: fmt.Sprintf("YOu %v got the fLag", email),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// PostMeHandler godoc
//
//	@Summary	Post to Me Check
//	@Accept		json
//	@Produce	json
//	@Router		/me [post]
func (h *MeHandler) PostMeHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Read the stream (This works because your middleware beautifully restored it!)
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}

	// 2. Convert bytes to string so it prints human-readable JSON
	payloadString := string(bodyBytes)
	if payloadString == "" {
		payloadString = "No payload provided"
	}

	response := dm.NyehResponse{
		Status:  "ok",
		Message: fmt.Sprintf("Your request payload is: %s", payloadString),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
