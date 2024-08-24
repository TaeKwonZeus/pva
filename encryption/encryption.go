package encryption

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/argon2"
	"net/http"
)

const (
	nonceSize   = 12
	saltSize    = 32
	aesKeySize  = 32
	hmacKeySize = 32
	rsaKeySize  = 4096
)

type Keys struct {
	signingKey []byte
	// Used to encrypt passwords and send back to clients so they don't have to supply
	// unencrypted passwords with each request.
	passwordKey []byte
}

func NewKeys() (*Keys, error) {
	keys := &Keys{
		signingKey:  make([]byte, hmacKeySize),
		passwordKey: make([]byte, aesKeySize),
	}

	if _, err := rand.Read(keys.signingKey); err != nil {
		return nil, err
	}
	if _, err := rand.Read(keys.passwordKey); err != nil {
		return nil, err
	}

	return keys, nil
}

func (k *Keys) SigningKey() []byte {
	return k.signingKey
}

func (k *Keys) PasswordKey() []byte {
	return k.passwordKey
}

func (k *Keys) Erase() {
	for i := range k.signingKey {
		k.signingKey[i] = 0
	}
	for i := range k.passwordKey {
		k.passwordKey[i] = 0
	}
}

func CreateKeypair() (private []byte, public []byte, err error) {
	prKey, err := rsa.GenerateKey(rand.Reader, rsaKeySize)
	if err != nil {
		return
	}
	pubKey := prKey.PublicKey

	return x509.MarshalPKCS1PrivateKey(prKey), x509.MarshalPKCS1PublicKey(&pubKey), nil
}

func GenerateSalt() ([]byte, error) {
	salt := make([]byte, saltSize)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}
	return salt, nil
}

func DeriveKey(password string, salt []byte) []byte {
	return argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, aesKeySize)
}

func AesEncrypt(plaintext, key, aad []byte) ([]byte, error) {
	nonce := make([]byte, nonceSize)
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nil, nonce, plaintext, aad)
	return bytes.Join([][]byte{nonce, ciphertext}, nil), nil
}

func AesDecrypt(ciphertext, key, aad []byte) ([]byte, error) {
	nonce := ciphertext[:nonceSize]
	ciphertext = ciphertext[nonceSize:]

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	plaintext, err := gcm.Open(nil, nonce, ciphertext, aad)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func AuthMiddleware(signingKey []byte) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
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

				return signingKey, nil
			})
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			claims, ok := t.Claims.(jwt.MapClaims)
			sub, subOk := claims["sub"].(string)
			if !(ok && t.Valid) || !subOk {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			r.Header.Add("username", sub)
			// encrypted password for cryptographic shit
			if passwd, ok := claims["passwd"].(string); ok {
				r.Header.Add("passwd", passwd)
			}

			next.ServeHTTP(w, r)
		})
	}
}
