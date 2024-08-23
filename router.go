package main

import (
	"github.com/TaeKwonZeus/pva/encryption"
	"github.com/TaeKwonZeus/pva/frontend"
	"github.com/TaeKwonZeus/pva/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func newRouter(e *handlers.Env) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Route("/api/auth", func(r chi.Router) {
		r.Use(middleware.RequestID)

		r.Post("/login", e.LoginHandler)
		r.Post("/register", e.RegisterHandler)
		r.Post("/revoke", e.Revoke)
	})

	r.Route("/api", func(r chi.Router) {
		r.Use(middleware.RequestID)
		r.Use(encryption.AuthMiddleware(e.Keys.SigningKey()))

		r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("pong"))
		})
	})

	r.Route("/", func(r chi.Router) {
		r.Get("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFileFS(w, r, frontend.Embed(), "favicon.ico")
		})
		r.Mount("/assets/", http.FileServerFS(frontend.Embed()))
		r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFileFS(w, r, frontend.Embed(), "index.html")
		})
	})

	return r
}
