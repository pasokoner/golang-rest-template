package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	mux.NotFound(app.notFound)
	mux.MethodNotAllowed(app.methodNotAllowed)

	mux.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	mux.Use(app.logAccess)
	mux.Use(app.recoverPanic)
	mux.Use(app.authenticate)

	mux.Get("/status", app.status)
	mux.Post("/signup", app.handlerCreateUser)
	mux.Post("/login", app.createAuthenticationToken)

	mux.Group(func(mux chi.Router) {
		mux.Use(app.requireAuthenticatedUser)

		mux.Get("/users/me", app.handlerGetUser)

		mux.Put("/security/email", app.handlerUpdateUserEmail)
		mux.Put("/security/password", app.handlerUpdateUserPassword)
	})

	mux.Group(func(mux chi.Router) {
		mux.Use(app.requireBasicAuthentication)

		mux.Get("/basic-auth-protected", app.protected)
	})

	return mux
}
