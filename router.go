package main

import (
	"github.com/TaeKwonZeus/pva/frontend"
	"github.com/TaeKwonZeus/pva/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func newRouter(env *handlers.Env) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Route("/api/auth", func(r chi.Router) {
		r.Use(middleware.RequestID)

		r.Post("/login", env.LoginHandler)
		r.Post("/register", env.RegisterHandler)
		r.Post("/revoke", env.Revoke)
	})

	r.Route("/api", func(r chi.Router) {
		r.Use(middleware.RequestID)
		r.Use(env.AuthMiddleware)

		r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("pong"))
		})

		r.Get("/vaults", env.GetVaultsHandler)
		r.Post("/vaults/new", env.NewVaultHandler)
		r.Post("/vaults/{id}/new", env.NewPasswordHandler)
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
