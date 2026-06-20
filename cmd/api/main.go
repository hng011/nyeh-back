package main

import (
	"fmt"
	"log"
	"net/http"
	a "nyeh-back/internal/api"
	c "nyeh-back/internal/core"
	"time"
)

func main() {
	c.LoadEnv()

	addr := fmt.Sprintf(":%v", c.Settings.PORT)

	server := &http.Server{
		Addr:         addr,
		Handler:      a.Setup(),
		ReadTimeout:  10 * time.Second,  // Drop slow clients
		WriteTimeout: 10 * time.Second,  // Drop hung responses
		IdleTimeout:  120 * time.Second, // Clean up dead connections
	}

	messages := []string{
		"Setup Complete",
		fmt.Sprintf("Mode: %v", c.Settings.ENV),
		fmt.Sprintf("Server running on port %v", server.Addr),
	}

	for _, x := range messages {
		log.Printf("[INFO] %v\n", x)
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start:\n%v", err)
	}
}
