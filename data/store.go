package data

import (
	"database/sql"
	"time"
)

// Store abstracts away cryptographic operations on data from db.
type Store struct {
	db          *db
	passwordKey []byte
}

func NewStore(path string, passwordKey []byte) (*Store, error) {
	pool, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	_, err = pool.Exec(startupQuery)
	if err != nil {
		return nil, err
	}

	return &Store{db: &db{pool}, passwordKey: passwordKey}, nil
}

func (s *Store) Close() error {
	return s.db.pool.Close()
}

func (s *Store) EncryptPassword(password string) ([]byte, error) {
	return aesEncrypt([]byte(password), s.passwordKey, nil)
}

func (s *Store) DecryptPassword(passwordEncrypted []byte) (string, error) {
	password, err := aesDecrypt(passwordEncrypted, s.passwordKey, nil)
	return string(password), err
}

func (s *Store) VerifyPassword(username string, password string) (verified bool, user *User) {
	user, err := s.GetUserByUsername(username)
	if err != nil {
		return false, nil
	}
	if user == nil {
		return false, nil
	}

	key := deriveKey(password, user.salt)
	if _, err := aesDecrypt(user.privateKeyEncrypted, key, nil); err != nil {
		return false, nil
	}

	return true, user
}

func (s *Store) GetUser(id int) (*User, error) {
	return s.db.getUser(id)
}

func (s *Store) GetUserByUsername(username string) (*User, error) {
	return s.db.getUserByUsername(username)
}

func (s *Store) CreateUser(user *User, password string) error {
	privateKey, publicKey, err := newKeypair()
	if err != nil {
		return err
	}

	user.salt, err = generateSalt()
	if err != nil {
		return err
	}

	key := deriveKey(password, user.salt)
	user.privateKeyEncrypted, err = aesEncrypt(privateKey, key, nil)
	user.publicKey = publicKey
	if err != nil {
		return err
	}

	return s.db.createUser(user)
}

func (s *Store) CreateVault(vault *Vault, owner *User) error {
	key, err := newAesKey()
	if err != nil {
		return err
	}

	vaultKeyEncrypted, err := rsaEncrypt(key, owner.publicKey, []byte("vault"))
	if err != nil {
		return err
	}

	vault.OwnerId = owner.Id
	return s.db.createVault(vault, vaultKeyEncrypted)
}

func decryptVault(vnk *vaultAndKey, user *User, userKey []byte) (*Vault, error) {
	privateKey, err := aesDecrypt(user.privateKeyEncrypted, userKey, nil)
	if err != nil {
		return nil, err
	}

	vaultKey, err := rsaDecrypt(vnk.keyEncrypted, privateKey, []byte("vault"))
	if err != nil {
		return nil, err
	}

	for _, p := range vnk.vault.Passwords {
		passwordDecrypted, err := aesDecrypt(p.passwordEncrypted, vaultKey, nil)
		if err != nil {
			return nil, err
		}

		p.Password = string(passwordDecrypted)
	}

	return vnk.vault, nil
}

func (s *Store) GetVault(id int, user *User, userKey []byte) (*Vault, error) {
	vnk, err := s.db.getVault(id, user.Id)
	if err != nil {
		return nil, err
	}
	if vnk == nil {
		return nil, nil
	}

	return decryptVault(vnk, user, userKey)
}

func (s *Store) GetVaults(user *User, userKey []byte) ([]*Vault, error) {
	vnks, err := s.db.getVaults(user.Id)
	if err != nil {
		return nil, err
	}

	out := make([]*Vault, len(vnks))

	for i, vnk := range vnks {
		out[i], err = decryptVault(vnk, user, userKey)
		if err != nil {
			return nil, err
		}
	}

	return out, nil
}

func (s *Store) CreatePassword(password *Password, vaultId int, user *User, userKey []byte) error {
	password.CreatedAt = time.Now()
	password.UpdatedAt = time.Now()

	vnk, err := s.db.getVault(vaultId, user.Id)
	if err != nil {
		return err
	}

	privateKey, err := aesDecrypt(user.privateKeyEncrypted, userKey, nil)
	if err != nil {
		return err
	}
	vaultKey, err := rsaDecrypt(vnk.keyEncrypted, privateKey, []byte("vault"))
	if err != nil {
		return err
	}
	password.passwordEncrypted, err = aesEncrypt([]byte(password.Password), vaultKey, nil)
	if err != nil {
		return err
	}

	return s.db.createPassword(password, vaultId)
}
