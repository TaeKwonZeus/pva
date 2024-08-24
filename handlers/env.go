package handlers

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/TaeKwonZeus/pva/data"
	"github.com/TaeKwonZeus/pva/encryption"
	"github.com/golang-jwt/jwt"
	"log"
	"net/http"
	"strconv"
)

type Env struct {
	DB   *data.DB
	Keys *encryption.Keys
}

func serverError(w http.ResponseWriter, err error) {
	log.Println("Server failure:", err)
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

// AuthMiddleware verifies the JWT token, fetches the calling user from e.DB, decrypts password in JWT and stores
// the user and the password in r.Context().
func (e *Env) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenCookie, err := r.Cookie("token")
		if err != nil || tokenCookie.Value == "" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		token := tokenCookie.Value

		t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return e.Keys.SigningKey(), nil
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		claims, ok := t.Claims.(jwt.MapClaims)
		sub, subOk := claims["sub"].(string)
		passwd, passwdOk := claims["passwd"].(string)
		if !ok || !t.Valid || !subOk || !passwdOk {
			log.Println(ok, t.Valid, subOk, passwdOk)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		id, err := strconv.Atoi(sub)
		if err != nil {
			http.Error(w, "Failed to parse sub", http.StatusUnauthorized)
			return
		}

		passwdBytes, err := base64.StdEncoding.DecodeString(passwd)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		password, err := encryption.AesDecrypt(passwdBytes, e.Keys.PasswordKey(), nil)
		if err != nil {
			http.Error(w, "Failed to decrypt password", http.StatusUnauthorized)
			return
		}
		user, err := e.DB.GetUser(id)
		if err != nil {
			http.Error(w, "Failed to fetch user", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user", user)
		ctx = context.WithValue(ctx, "password", string(password))
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
