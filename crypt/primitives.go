package crypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"golang.org/x/crypto/argon2"
)

const (
	nonceSize  = 12
	saltSize   = 32
	aesKeySize = 32
	rsaKeySize = 4096
)

//
//type Keys struct {
//	signingKey []byte
//	// Used to encrypt passwords and send back to clients so they don't have to supply
//	// unencrypted passwords with each request.
//	passwordKey []byte
//}
//
//func NewKeys() (*Keys, error) {
//	keys := &Keys{
//		signingKey:  make([]byte, hmacKeySize),
//		passwordKey: make([]byte, aesKeySize),
//	}
//
//	// FIXME replace later
//	//if _, err := rand.Read(keys.signingKey); err != nil {
//	//	return nil, err
//	//}
//	//if _, err := rand.Read(keys.passwordKey); err != nil {
//	//	return nil, err
//	//}
//
//	return keys, nil
//}
//
//func (k *Keys) SigningKey() []byte {
//	return k.signingKey
//}
//
//func (k *Keys) PasswordKey() []byte {
//	return k.passwordKey
//}
//
//func (k *Keys) Erase() {
//	for i := range k.signingKey {
//		k.signingKey[i] = 0
//	}
//	for i := range k.passwordKey {
//		k.passwordKey[i] = 0
//	}
//}

func NewKeypair() (private []byte, public []byte, err error) {
	prKey, err := rsa.GenerateKey(rand.Reader, rsaKeySize)
	if err != nil {
		return
	}
	pubKey := prKey.PublicKey

	return x509.MarshalPKCS1PrivateKey(prKey), x509.MarshalPKCS1PublicKey(&pubKey), nil
}

func NewAesKey() ([]byte, error) {
	key := make([]byte, aesKeySize)
	if _, err := rand.Read(key); err != nil {
		return nil, err
	}
	return key, nil
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

func AesEncrypt(plaintext, key []byte) ([]byte, error) {
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

	ciphertext := gcm.Seal(nil, nonce, plaintext, nil)
	return bytes.Join([][]byte{nonce, ciphertext}, nil), nil
}

func AesDecrypt(ciphertext, key []byte) ([]byte, error) {
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

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func RsaEncrypt(plaintext, key []byte) ([]byte, error) {
	pk, err := x509.ParsePKCS1PublicKey(key)
	if err != nil {
		return nil, err
	}
	return rsa.EncryptOAEP(sha256.New(), rand.Reader, pk, plaintext, nil)
}

func RsaDecrypt(ciphertext, key []byte) ([]byte, error) {
	pk, err := x509.ParsePKCS1PrivateKey(key)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptOAEP(sha256.New(), rand.Reader, pk, ciphertext, nil)
}
