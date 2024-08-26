package handlers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/TaeKwonZeus/pva/data"
	"github.com/golang-jwt/jwt"
	"log"
	"net/http"
	"strconv"
	"time"
)

type jwtClaims struct {
	jwt.StandardClaims

	Passwd string `json:"passwd,omitempty"`
}

func (j *jwtClaims) Valid() error {
	if j.Passwd == "" {
		return errors.New("jwt claims missing passwd")
	}
	return j.StandardClaims.Valid()
}

// AuthMiddleware verifies the JWT token, fetches the calling user from e.db, decrypts password in JWT and stores
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

		claims, ok := t.Claims.(*jwtClaims)
		if !t.Valid || !ok {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		id, err := strconv.Atoi(claims.Subject)
		if err != nil {
			http.Error(w, "Failed to parse sub", http.StatusUnauthorized)
			return
		}

		user, err := e.Store.GetUser(id)
		if err != nil {
			http.Error(w, "Failed to get user", http.StatusUnauthorized)
			return
		}
		if user == nil {
			http.Error(w, "User not found", http.StatusUnauthorized)
			return
		}

		passwdBytes, err := base64.StdEncoding.DecodeString(claims.Passwd)
		if err != nil {
			log.Println("Failed to decode passwd")
			http.Error(w, "Failed to decode passwd", http.StatusUnauthorized)
			return
		}

		password, err := e.Store.DecryptPassword(passwdBytes)
		if err != nil {
			http.Error(w, "Failed to decrypt password", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user", user)
		ctx = context.WithValue(ctx, "password", password)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

type credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (e *Env) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var c credentials
	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	verified, user := e.Store.VerifyPassword(c.Username, c.Password)
	if !verified {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	passwd, err := e.Store.EncryptPassword(c.Password)
	if err != nil {
		serverError(w, err)
		return
	}

	remember := r.URL.Query().Get("remember") == "true"

	var exp time.Time
	if remember {
		exp = time.Now().AddDate(0, 1, 0)
	} else {
		exp = time.Now().Add(time.Hour * 4)
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":    strconv.Itoa(user.Id),
		"iss":    r.Host,
		"aud":    r.RemoteAddr,
		"iat":    time.Now().Unix(),
		"exp":    exp.Unix(),
		"passwd": base64.StdEncoding.EncodeToString(passwd),
	}).SignedString(e.Keys.SigningKey())
	if err != nil {
		serverError(w, err)
		return
	}

	cookie := &http.Cookie{
		Name:     "token",
		Value:    token,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	}
	if remember {
		cookie.Expires = exp
	}
	http.SetCookie(w, cookie)
}

func (e *Env) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var c credentials
	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = e.Store.CreateUser(&data.User{
		Username: c.Username,
		Role:     data.RoleAdmin,
	}, c.Password)
	if data.IsErrConflict(err) {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}
	if err != nil {
		serverError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (_ *Env) Revoke(w http.ResponseWriter, _ *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "token",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
}
