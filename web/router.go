package main

import (
	"fmt"
	"github.com/TaeKwonZeus/pva/frontend"
	"github.com/TaeKwonZeus/pva/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/golang-jwt/jwt"
	"log"
	"net/http"
	"strings"
)

func newRouter(e *handlers.Env) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Mount("/api/auth", authRouter(e))
	r.Mount("/api", apiRouter(e))

	r.Mount("/", frontendRouter())

	return r
}

func authRouter(e *handlers.Env) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)

	r.Post("/login", e.LoginHandler)
	r.Post("/register", e.RegisterHandler)

	return r
}

func authMiddleware(signingKey []byte) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, _ := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
			if token == "" {
				log.Printf("%s: Missing token", r.RemoteAddr)
				http.Error(w, "Missing token", http.StatusUnauthorized)
				return
			}

			t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}

				return signingKey, nil
			})
			if err != nil {
				log.Printf("%s: %v", r.RemoteAddr, err)
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			if claims, ok := t.Claims.(jwt.MapClaims); ok && t.Valid {
				if (claims["iss"] != r.Host) || (claims["aud"] != r.RemoteAddr) {
					http.Error(w, "Host or peer doesn't match", http.StatusUnauthorized)
				}
				r.Header.Add("username", claims["sub"].(string))
			} else {
				log.Printf("%s: Missing \"sub\" claim", r.RemoteAddr)
			}

			next.ServeHTTP(w, r)
		})
	}
}

func apiRouter(e *handlers.Env) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(authMiddleware(e.SigningKey))

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("pong"))
	})

	return r
}

func frontendRouter() http.Handler {
	r := chi.NewRouter()

	r.Get("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFileFS(w, r, frontend.Embed(), "favicon.ico")
	})

	r.Mount("/assets/", http.FileServerFS(frontend.Embed()))

	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFileFS(w, r, frontend.Embed(), "index.html")
	})

	return r
}
