package middleware

import (
	"bytes"
	"io"
	"log"
	"net/http"
)

func RequestPayloadLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// 1. Read the raw bytes from the stream
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("[ERROR] Failed to read request body: %v", err)
			// http.Error(w, "Failed to read request", http.StatusBadRequest)
			// return
		}

		// 2. Log the payload (avoiding empty logs for GET requests)
		if len(bodyBytes) > 0 {
			log.Printf("[INFO] %s %s | Payload: %s", r.Method, r.URL.Path, string(bodyBytes))
		}

		// 3. CRITICAL: Restore the body so the next handler can read it
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		next.ServeHTTP(w, r)
	})
}
