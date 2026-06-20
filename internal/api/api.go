package api

import (
	"net/http"
	v1 "nyeh-back/internal/api/v1"
	"nyeh-back/internal/middleware"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

const (
	BASE_PATH_V1 string = "/api/v1"
)

func Setup() http.Handler {
	r := chi.NewRouter()

	// Global Middleware
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(middleware.RequestPayloadLogger)

	r.Route(BASE_PATH_V1, func(r chi.Router) { v1.SetupApiV1(r, BASE_PATH_V1) })

	return r
}
