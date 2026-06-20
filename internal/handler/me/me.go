package me

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	d "nyeh-back/internal/domain/me"
)

// AdminHandler godoc
//
//	@Summary	Admin Check
//	@Accept		json
//	@Produce	json
//	@Router		/me [get]
func MeHandler(w http.ResponseWriter, r *http.Request) {
	response := d.NyehResponse{
		Status:  "ok",
		Message: "YOu got the fLag",
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
func PostMeHandler(w http.ResponseWriter, r *http.Request) {
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

	response := d.NyehResponse{
		Status:  "ok",
		Message: fmt.Sprintf("Your request payload is: %s", payloadString),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
