package handlers

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/TaeKwonZeus/pva/encryption"
	"github.com/golang-jwt/jwt"
	"net/http"
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

	var salt string
	var privateKeyEncrypted string
	err = e.Pool.QueryRow("SELECT salt, private_key_encrypted FROM users WHERE username = ?", c.Username).
		Scan(&salt, &privateKeyEncrypted)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	saltBytes, err := base64.StdEncoding.DecodeString(salt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	privateKeyEncryptedBytes, err := base64.StdEncoding.DecodeString(privateKeyEncrypted)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	key := encryption.DeriveKey(c.Password, saltBytes)
	_, err = encryption.AesDecrypt(privateKeyEncryptedBytes, key, nil)
	if err != nil {
		http.Error(w, "Failed to verify identity", http.StatusUnauthorized)
		return
	}

	passwd, err := encryption.AesEncrypt([]byte(c.Password), e.Keys.PasswordKey(), nil)
	if err != nil {
		http.Error(w, "Failed to encrypt password", http.StatusInternalServerError)
		return
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":    c.Username,
		"iss":    r.Host,
		"aud":    r.RemoteAddr,
		"iat":    time.Now().Unix(),
		"exp":    time.Now().Add(time.Hour * 2).Unix(),
		"passwd": base64.StdEncoding.EncodeToString(passwd),
	}).SignedString(e.Keys.SigningKey())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})
	w.WriteHeader(http.StatusOK)
}

func (e *Env) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var c credentials
	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !errors.Is(e.Pool.QueryRow("SELECT username FROM users WHERE username = ?", c.Username).Scan(), sql.ErrNoRows) {
		http.Error(w, "User already exists", http.StatusBadRequest)
		return
	}

	privateKey, publicKey, err := encryption.CreateKeypair()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	salt, err := encryption.GenerateSalt()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	key := encryption.DeriveKey(c.Password, salt)
	privateKeyEncrypted, err := encryption.AesEncrypt(privateKey, key, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = e.Pool.Exec(
		"INSERT INTO users (username, salt, public_key, private_key_encrypted, role) VALUES (?, ?, ?, ?, ?)",
		c.Username,
		base64.StdEncoding.EncodeToString(salt),
		base64.StdEncoding.EncodeToString(publicKey),
		base64.StdEncoding.EncodeToString(privateKeyEncrypted),
		"None",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (_ *Env) Revoke(w http.ResponseWriter, _ *http.Request) {
	// TODO fix
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
	})
	w.WriteHeader(http.StatusOK)
}
