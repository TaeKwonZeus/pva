package handlers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/TaeKwonZeus/pva/crypt"
	"github.com/TaeKwonZeus/pva/data"
	"github.com/charmbracelet/log"
	"net/http"
	"time"
)

type authToken struct {
	Issuer   string `json:"iss"`
	Audience string `json:"aud"`
	IssuedAt int64  `json:"iat"`
	Expires  int64  `json:"exp"`
	UserId   int    `json:"uid"`
	Key      string `json:"key"`
}

func decryptToken(encrypted string, key []byte) (t authToken, err error) {
	bytes, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return
	}
	plaintext, err := crypt.AesDecrypt(bytes, key)
	if err != nil {
		return
	}
	err = json.Unmarshal(plaintext, &t)
	return
}

func (t authToken) encryptedString(key []byte) (string, error) {
	j, err := json.Marshal(t)
	if err != nil {
		return "", err
	}
	e, err := crypt.AesEncrypt(j, key)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(e), nil
}

func (t authToken) valid(r *http.Request) error {
	if t.Issuer != r.Host {
		return errors.New("invalid issuer")
	}
	if t.Audience != r.RemoteAddr {
		return errors.New("invalid audience")
	}
	if t.IssuedAt > time.Now().Unix() || t.Expires < time.Now().Unix() {
		return errors.New("token expired")
	}
	if t.UserId == 0 || t.Key == "" {
		return errors.New("invalid id or key")
	}
	return nil
}

func (e *Env) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenCookie, err := r.Cookie("token")
		if err != nil || tokenCookie.Value == "" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		token := tokenCookie.Value

		t, err := decryptToken(token, e.TokenKey)
		if err != nil {
			log.Warn(err.Error(), "token", token)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if err = t.valid(r); err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		user, err := e.Store.GetUser(t.UserId)
		if data.IsErrNotFound(err) {
			http.Error(w, "user not found", http.StatusUnauthorized)
			return
		}
		if err != nil {
			http.Error(w, "failed to get user", http.StatusUnauthorized)
			return
		}

		keyBytes, err := base64.StdEncoding.DecodeString(t.Key)
		if err != nil {
			http.Error(w, "invalid key", http.StatusUnauthorized)
			return
		}

		_, err = user.DecryptPrivateKey(keyBytes)
		if err != nil {
			http.Error(w, "invalid private key", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user", user)

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

	user, err := e.Store.GetUserByUsername(c.Username)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	userKey := user.DeriveKey(c.Password)
	_, err = user.DecryptPrivateKey(userKey)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	remember := r.URL.Query().Get("remember") == "true"

	var exp time.Time
	if remember {
		exp = time.Now().AddDate(0, 1, 0)
	} else {
		exp = time.Now().Add(time.Hour * 4)
	}

	token, err := authToken{
		Issuer:   r.Host,
		Audience: r.RemoteAddr,
		IssuedAt: time.Now().Unix(),
		Expires:  exp.Unix(),
		UserId:   user.ID,
		Key:      base64.StdEncoding.EncodeToString(userKey),
	}.encryptedString(e.TokenKey)
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

func authenticate(w http.ResponseWriter, r *http.Request, permission data.Permission) (user *data.User, ok bool) {
	user, ok = r.Context().Value("user").(*data.User)
	if !ok {
		log.Error("could not get user from context")
		w.WriteHeader(http.StatusInternalServerError)
		return nil, false
	}

	if permission != data.PermissionNone && !data.CheckPermission(user.Role, permission) {
		http.Error(w, "permission not satisfied: "+string(permission), http.StatusForbidden)
		return nil, false
	}
	return user, true
}
