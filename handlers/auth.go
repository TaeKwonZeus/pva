package handlers

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/TaeKwonZeus/pva/data"
	"github.com/TaeKwonZeus/pva/encryption"
	"github.com/golang-jwt/jwt"
	"net/http"
	"strconv"
	"time"
)

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

	remember := r.URL.Query().Get("remember") == "true"

	user, err := e.DB.GetUserByUsername(c.Username)
	if err != nil {
		serverError(w, err)
		return
	}

	key := encryption.DeriveKey(c.Password, user.Salt)
	_, err = encryption.AesDecrypt(user.PrivateKeyEncrypted, key, nil)
	if err != nil {
		http.Error(w, "Failed to verify identity", http.StatusUnauthorized)
		return
	}

	// User password encrypted with e.Keys.PasswordKey()
	passwd, err := encryption.AesEncrypt([]byte(c.Password), e.Keys.PasswordKey(), nil)
	if err != nil {
		serverError(w, err)
		return
	}

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

	privateKey, publicKey, err := encryption.NewKeypair()
	if err != nil {
		serverError(w, err)
		return
	}
	salt, err := encryption.GenerateSalt()
	if err != nil {
		serverError(w, err)
		return
	}
	key := encryption.DeriveKey(c.Password, salt)
	privateKeyEncrypted, err := encryption.AesEncrypt(privateKey, key, nil)
	if err != nil {
		serverError(w, err)
		return
	}

	err = e.DB.AddUser(&data.User{
		Username:            c.Username,
		Salt:                salt,
		PublicKey:           publicKey,
		PrivateKeyEncrypted: privateKeyEncrypted,
		Role:                data.RoleAdmin,
	})
	if errors.Is(err, data.ErrorConflict) {
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
