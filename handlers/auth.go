package handlers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/TaeKwonZeus/pva/data"
	"github.com/charmbracelet/log"
	"github.com/golang-jwt/jwt"
	"net/http"
	"strconv"
	"time"
)

type jwtClaims struct {
	jwt.StandardClaims

	Passwd string `json:"passwd,omitempty"`
}

func (j jwtClaims) Valid() error {
	if j.Passwd == "" {
		return errors.New("jwt claims missing passwd")
	}
	return j.StandardClaims.Valid()
}

// AuthMiddleware verifies the JWT token, fetches the calling user from e.db, decrypts password in JWT and stores
// the user and the password in r.Context().
func (e *Env) AuthMiddleware(decrypt bool) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenCookie, err := r.Cookie("token")
			if err != nil || tokenCookie.Value == "" {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			token := tokenCookie.Value

			t, err := jwt.ParseWithClaims(token, new(jwtClaims), func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}

				return e.Keys.SigningKey(), nil
			})
			if err != nil {
				log.Warn(err.Error(), "token", token)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			claims, ok := t.Claims.(*jwtClaims)
			if !t.Valid || !ok {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}
			id, err := strconv.Atoi(claims.Subject)
			if err != nil {
				http.Error(w, "failed to parse sub", http.StatusUnauthorized)
				return
			}

			user, err := e.Store.GetUser(id)
			if err != nil {
				http.Error(w, "failed to get user", http.StatusUnauthorized)
				return
			}
			if user == nil {
				http.Error(w, "user not found", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "user", user)

			if decrypt {
				passwdBytes, err := base64.StdEncoding.DecodeString(claims.Passwd)
				if err != nil {
					log.Error("failed to decode", "passwd", claims.Passwd)
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				password, err := e.Store.DecryptPassword(passwdBytes)
				if err != nil {
					log.Warn("failed to decrypt", "passwd", claims.Passwd)
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				userKey := user.DeriveKey(password)
				ctx = context.WithValue(ctx, "userKey", userKey)
			}

			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
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
		log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	remember := r.URL.Query().Get("remember") == "true"

	var exp time.Time
	if remember {
		exp = time.Now().AddDate(0, 1, 0)
	} else {
		exp = time.Now().Add(time.Hour * 4)
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwtClaims{
		StandardClaims: jwt.StandardClaims{
			Subject:   strconv.Itoa(user.ID),
			Issuer:    r.Host,
			Audience:  r.RemoteAddr,
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: exp.Unix(),
		},
		Passwd: base64.StdEncoding.EncodeToString(passwd),
	}).SignedString(e.Keys.SigningKey())
	if err != nil {
		log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
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

	user := data.User{
		Username: c.Username,
	}

	n, err := e.Store.GetUserCount()
	if err != nil {
		log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if n == 0 {
		user.Role = data.RoleAdmin
	} else {
		user.Role = data.RoleViewer
	}

	err = e.Store.CreateUser(&user, c.Password)
	if data.IsErrConflict(err) {
		http.Error(w, "user already exists", http.StatusConflict)
		return
	}
	if err != nil {
		log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
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

func authenticate(w http.ResponseWriter, r *http.Request, permission data.Permission) (user *data.User, userKey []byte, ok bool) {
	user, ok = r.Context().Value("user").(*data.User)
	if !ok {
		log.Error("could not get user from context")
		w.WriteHeader(http.StatusInternalServerError)
		return nil, nil, false
	}

	userKey, ok = r.Context().Value("userKey").([]byte)
	if !ok {
		log.Warn("could not get user key from context")
		w.WriteHeader(http.StatusInternalServerError)
		return nil, nil, false
	}

	if permission >= 0 && !data.CheckPermission(user.Role, permission) {
		w.WriteHeader(http.StatusForbidden)
		return nil, nil, false
	}
	return user, userKey, true
}

func authenticateNoKey(w http.ResponseWriter, r *http.Request, permission data.Permission) (user *data.User, ok bool) {
	user, ok = r.Context().Value("user").(*data.User)
	if !ok {
		log.Error("could not get user from context")
		w.WriteHeader(http.StatusInternalServerError)
		return nil, false
	}

	if permission >= 0 && !data.CheckPermission(user.Role, permission) {
		w.WriteHeader(http.StatusForbidden)
		return nil, false
	}

	return user, true
}
