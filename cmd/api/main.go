package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"nyeh-back/internal/api"
	v1 "nyeh-back/internal/api/v1"
	"nyeh-back/internal/core"

	"nyeh-back/internal/infra"

	handler "nyeh-back/internal/handler/http"
	handler_auth "nyeh-back/internal/handler/http/auth"

	firestore "nyeh-back/internal/repository/firestore"
	"nyeh-back/internal/repository/redis"
	usecase "nyeh-back/internal/usecase"

	"time"
)

func main() {
	core.LoadEnv()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	addr := fmt.Sprintf(":%v", core.Settings.PORT)

	infra.InitRedis(
		core.Settings.REDIS_ADDR,
		core.Settings.REDIS_USER,
		core.Settings.REDIS_PASS,
	)
	defer infra.RedisClient.Close()

	firestoreClient := infra.InitFirestore(ctx)
	defer firestoreClient.Close()

	// API V1 DI
	v1Deps := v1.V1RouterDependencies{
		AuthHandler: handler_auth.NewAuthHandler(usecase.NewAuthUsecase(redis.NewRedisSessionRepo(infra.RedisClient))),
		BioHandler:  handler.NewBioHandler(usecase.NewBioUsecase(firestore.NewFirestoreRepo(firestoreClient))),
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
