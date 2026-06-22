package v1

import (
	docs "nyeh-back/docs"
	handler "nyeh-back/internal/handler"
	auth "nyeh-back/internal/handler/auth"
	me "nyeh-back/internal/handler/me"
	mid "nyeh-back/internal/middleware"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

// RouterDependencies holds all the fully-built handlers
type V1RouterDependencies struct {
	AuthHandler *auth.AuthHandler
	MeHandler   *me.MeHandler
}

func SetupApiV1(r chi.Router, base_path string, v1Deps V1RouterDependencies) {

	docs.SwaggerInfo.Title = "NYEH-BACK API DOCS"
	docs.SwaggerInfo.BasePath = base_path
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	r.Get("/swagger/*", httpSwagger.Handler())
	r.Get("/healthCheck", handler.HealthCheckHandler)
	authRouter(r, "/auth", v1Deps.AuthHandler)

	r.Group(func(r chi.Router) {
		r.Use(mid.RequireAuth)
		r.Get("/me", v1Deps.MeHandler.GetMeHandler)
		r.Post("/me", v1Deps.MeHandler.PostMeHandler)
	})
}

func authRouter(r chi.Router, baseUrl string, authDeps *auth.AuthHandler) {
	r.Route(baseUrl, func(r chi.Router) {

		r.Route("/google", func(r chi.Router) {
			r.Get("/login", authDeps.GoogleLoginHandler)
			r.Get("/callback", authDeps.GoogleCallbackHandler)
		})
		r.Post("/refresh", authDeps.RefreshHandler)
		r.Post("/logout", authDeps.LogoutHandler)
	})
}
