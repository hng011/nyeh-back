package main

import (
	"fmt"
	"log"
	"net/http"

	"nyeh-back/internal/api"
	v1 "nyeh-back/internal/api/v1"
	"nyeh-back/internal/core"

	"nyeh-back/internal/infra"

	// Import Handlers
	"nyeh-back/internal/cache/redis"
	auth "nyeh-back/internal/handler/auth"
	me "nyeh-back/internal/handler/me"

	"time"
)

func main() {
	core.LoadEnv()

	addr := fmt.Sprintf(":%v", core.Settings.PORT)

	infra.InitRedis(
		core.Settings.REDIS_ADDR,
		core.Settings.REDIS_PASSWORD,
	)
	defer infra.RedisClient.Close()

	// API V1 DI
	v1Deps := v1.V1RouterDependencies{
		AuthHandler: auth.NewAuthHandler(redis.NewUserSessionCache(infra.RedisClient)),
		MeHandler:   me.NewMeHandler(),
	}

	server := &http.Server{
		Addr:         addr,
		Handler:      api.Setup(v1Deps),
		ReadTimeout:  10 * time.Second,  // Drop slow clients
		WriteTimeout: 10 * time.Second,  // Drop hung responses
		IdleTimeout:  120 * time.Second, // Clean up dead connections
	}

	messages := []string{
		"Setup Complete",
		fmt.Sprintf("Mode: %v", core.Settings.ENV),
		fmt.Sprintf("Server running on port %v", server.Addr),
	}

	for _, x := range messages {
		log.Printf("[INFO] %v\n", x)
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start:\n%v", err)
	}
}
