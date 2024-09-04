package main

import (
	"github.com/TaeKwonZeus/pva/frontend"
	"github.com/TaeKwonZeus/pva/handlers"
	"github.com/charmbracelet/log"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

type loggingRw struct {
	http.ResponseWriter
	status int
}

func (lrw *loggingRw) WriteHeader(status int) {
	lrw.status = status
	lrw.ResponseWriter.WriteHeader(status)
}

func logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lrw := loggingRw{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(&lrw, r)

		if lrw.status >= 500 {
			log.Info(r.URL.Path, "method", r.Method, "code", lrw.status)
		} else {
			log.Info(r.URL.Path, "method", r.Method, "code", lrw.status)
		}
	})
}

func newRouter(env *handlers.Env) http.Handler {
	r := chi.NewRouter()
	r.Use(logger)

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

		r.Get("/index", env.GetIndexHandler)

		r.Route("/vaults", func(r chi.Router) {
			r.Get("/", env.GetVaultsHandler)
			r.Post("/new", env.NewVaultHandler)
			r.Patch("/{id}", env.UpdateVaultHandler)
			r.Delete("/{id}", env.DeleteVaultHandler)
			r.Patch("/{id}", env.ShareVaultHandler)
			r.Post("/{id}/share", env.ShareVaultHandler)

			r.Post("/{id}/new", env.NewPasswordHandler)
			r.Patch("/{vaultId}/{passwordId}", env.UpdatePasswordHandler)
			r.Delete("/{vaultId}/{passwordId}", env.DeletePasswordHandler)
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
