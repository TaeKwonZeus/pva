package main

import (
	"github.com/TaeKwonZeus/pva/frontend"
	"github.com/TaeKwonZeus/pva/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func newRouter(_ *handlers.Env) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Mount("/api", apiRouter())
	r.Mount("/", frontendRouter())

	return r
}

func apiRouter() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)

	return r
}

func frontendRouter() http.Handler {
	r := chi.NewRouter()

	r.Get("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFileFS(w, r, frontend.Embed, "favicon.ico")
	})

	r.Mount("/assets/", http.FileServerFS(frontend.Embed))

	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFileFS(w, r, frontend.Embed, "index.html")
	})

	return r
}
