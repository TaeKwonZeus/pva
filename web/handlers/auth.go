package handlers

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/argon2"
	"net/http"
	"time"
)

type credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

const (
	nonceSize  = 12
	saltSize   = 32
	keySize    = 32
	rsaKeySize = 4096
)

func (e *Env) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var c credentials
	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var encryptedKey string
	err = e.Pool.QueryRow("SELECT encrypted_private_key FROM users WHERE username = ?", c.Username).Scan(&encryptedKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = aesDecrypt(encryptedKey, c.Password, nil)
	if err != nil {
		http.Error(w, "Failed to verify identity", http.StatusUnauthorized)
		return
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": c.Username,
		"iss": r.Host,
		"aud": r.RemoteAddr,
		"iat": time.Now().Unix(),
	}).SignedString(e.SigningKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"token": "` + token + `"}`))
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

	privateKey, publicKey, err := generateKeypair()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	encryptedPrivateKey, err := aesEncrypt(privateKey, c.Password, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = e.Pool.Exec("INSERT INTO users (username, public_key, encrypted_private_key) VALUES (?, ?, ?)",
		c.Username, publicKey, encryptedPrivateKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func generateKeypair() (private string, public string, err error) {
	prKey, err := rsa.GenerateKey(rand.Reader, rsaKeySize)
	if err != nil {
		return "", "", err
	}
	pubKey := prKey.PublicKey

	private = base64.StdEncoding.EncodeToString(x509.MarshalPKCS1PrivateKey(prKey))
	public = base64.StdEncoding.EncodeToString(x509.MarshalPKCS1PublicKey(&pubKey))
	return
}

func aesEncrypt(plaintext, password string, aad []byte) (string, error) {
	// Derive key
	salt := make([]byte, saltSize)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	key := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, keySize)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, nonceSize)
	if _, err = rand.Read(nonce); err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nil, nonce, []byte(plaintext), aad)
	out := bytes.Join([][]byte{salt, nonce, ciphertext}, nil)

	return base64.StdEncoding.EncodeToString(out), nil
}

func aesDecrypt(ciphertext, password string, aad []byte) (string, error) {
	ciphertextBytes, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}
	salt := ciphertextBytes[:saltSize]
	nonce := ciphertextBytes[saltSize : saltSize+nonceSize]
	ciphertextBytes = ciphertextBytes[saltSize+nonceSize:]

	key := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, keySize)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	plaintext, err := gcm.Open(nil, nonce, ciphertextBytes, aad)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
