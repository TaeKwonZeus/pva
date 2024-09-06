package main

import (
	"github.com/TaeKwonZeus/pva/frontend"
	"github.com/TaeKwonZeus/pva/handlers"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"strconv"
)

type loggingRw struct {
	http.ResponseWriter
	status int
}

func (lrw *loggingRw) WriteHeader(status int) {
	lrw.status = status
	lrw.ResponseWriter.WriteHeader(status)
}

var (
	infoFmt      = lipgloss.NewStyle().Foreground(lipgloss.Color("27")).Render
	okFmt        = lipgloss.NewStyle().Foreground(lipgloss.Color("120")).Render
	redirectFmt  = lipgloss.NewStyle().Foreground(lipgloss.Color("87")).Render
	clientErrFmt = lipgloss.NewStyle().Foreground(lipgloss.Color("226")).Render
	serverErrFmt = lipgloss.NewStyle().Foreground(lipgloss.Color("204")).Render
)

func logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lrw := loggingRw{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(&lrw, r)

		s := strconv.Itoa(lrw.status)
		log.DefaultStyles()
		switch {
		case lrw.status < 200:
			log.Info(infoFmt(s)+" "+r.Method, "path", r.URL.Path, "peer", r.RemoteAddr)
		case lrw.status < 300:
			log.Info(okFmt(s)+" "+r.Method, "path", r.URL.Path, "peer", r.RemoteAddr)
		case lrw.status < 400:
			log.Info(redirectFmt(s)+" "+r.Method, "path", r.URL.Path, "peer", r.RemoteAddr)
		case lrw.status < 500:
			log.Info(clientErrFmt(s)+" "+r.Method, "path", r.URL.Path, "peer", r.RemoteAddr)
		default:
			log.Error(serverErrFmt(s)+" "+r.Method, "path", r.URL.Path, "peer", r.RemoteAddr)
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
		r.Get("/vite.svg", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFileFS(w, r, frontend.Embed(), "vite.svg")
		})
		r.Mount("/assets/", http.FileServerFS(frontend.Embed()))
		r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFileFS(w, r, frontend.Embed(), "index.html")
		})
	})

	return r
}
