package v1

import (
	docs "nyeh-back/docs"
	handler "nyeh-back/internal/handler"
	auth "nyeh-back/internal/handler/auth"
	"nyeh-back/internal/handler/me"
	mid "nyeh-back/internal/middleware"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

func SetupApiV1(r chi.Router, base_path string) {

	docs.SwaggerInfo.Title = "NYEH-BACK API DOCS"
	docs.SwaggerInfo.BasePath = base_path
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	r.Get("/swagger/*", httpSwagger.Handler())
	r.Get("/healthCheck", handler.HealthCheckHandler)
	authRouter(r, "/auth")

	r.Group(func(r chi.Router) {
		r.Use(mid.RequireAuth)
		r.Get("/me", me.MeHandler)
		r.Post("/me", me.PostMeHandler)
	})
}

func authRouter(r chi.Router, baseUrl string) {
	r.Route(baseUrl, func(r chi.Router) {

		r.Route("/google", func(r chi.Router) {
			r.Get("/login", auth.GoogleLoginHandler)
			r.Get("/callback", auth.GoogleCallbackHandler)
		})

	})
}
