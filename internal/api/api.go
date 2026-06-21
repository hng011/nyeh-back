package api

import (
	"net/http"
	v1 "nyeh-back/internal/api/v1"

	"nyeh-back/internal/middleware"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

const (
	BASE_PATH_V1 string = "/api/v1"
)

func Setup(v1Deps v1.V1RouterDependencies) http.Handler {
	r := chi.NewRouter()

	// TODO: Parameterize this
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"*"},
		AllowCredentials: false,
	}))

	// Global Middleware
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(middleware.RequestPayloadLogger)

	r.Route(BASE_PATH_V1, func(r chi.Router) { v1.SetupApiV1(r, BASE_PATH_V1, v1Deps) })

	return r
}
