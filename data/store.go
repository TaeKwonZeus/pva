package data

import (
	"database/sql"
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

func (s *Store) CreateUser(user *User, password string) (id int, err error) {
	privateKey, publicKey, err := newKeypair()
	if err != nil {
		return 0, err
	}

	user.salt, err = generateSalt()
	if err != nil {
		return 0, err
	}

	key := deriveKey(password, user.salt)
	user.privateKeyEncrypted, err = aesEncrypt(privateKey, key, nil)
	user.publicKey = publicKey
	if err != nil {
		return 0, err
	}

	err = s.db.createUser(user)
	if err != nil {
		return 0, err
	}

	return user.Id, nil
}

func (s *Store) CreateVault(vault *Vault, owner *User) (id int, err error) {
	key, err := newAesKey()
	if err != nil {
		return 0, err
	}

	vaultKeyEncrypted, err := rsaEncrypt(key, owner.publicKey, []byte("vault"))
	if err != nil {
		return 0, err
	}

	vault.OwnerId = owner.Id
	err = s.db.createVault(vault, vaultKeyEncrypted)
	if err != nil {
		return 0, err
	}

	return vault.Id, nil
}

func decryptVault(vnk *vaultAndKey, user *User, passwordKey []byte) (*Vault, error) {
	privateKey, err := aesDecrypt(user.privateKeyEncrypted, passwordKey, nil)
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

func (s *Store) GetVault(id int, user *User, password string) (*Vault, error) {
	vnk, err := s.db.getVault(id, user.Id)
	if err != nil {
		return nil, err
	}

	key := deriveKey(password, user.salt)
	return decryptVault(vnk, user, key)
}

func (s *Store) GetVaults(user *User, password string) ([]*Vault, error) {
	vnks, err := s.db.getVaults(user.Id)
	if err != nil {
		return nil, err
	}

	out := make([]*Vault, len(vnks))

	key := deriveKey(password, user.salt)
	for i, vnk := range vnks {
		out[i], err = decryptVault(vnk, user, key)
		if err != nil {
			return nil, err
		}
	}

	return out, nil
}
