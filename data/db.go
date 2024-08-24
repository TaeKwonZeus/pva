package data

import (
	"database/sql"
	_ "embed"
	"encoding/base64"
	"errors"
	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed startup.sql
var startupQuery string

type DB struct {
	db *sql.DB
}

// ErrorConflict is returned on primary key or unique violations.
var ErrorConflict = errors.Join(sqlite3.ErrConstraintPrimaryKey, sqlite3.ErrConstraintUnique)

func NewDB(path string) (*DB, error) {
	pool, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	_, err = pool.Exec(startupQuery)
	if err != nil {
		return nil, err
	}

	return &DB{db: pool}, nil
}

func (d *DB) Close() error {
	return d.db.Close()
}

func (d *DB) AddUser(user *User) error {
	res, err := d.db.Exec(
		`INSERT INTO users (username, salt, public_key, private_key_encrypted, role)
		VALUES (?, ?, ?, ?, ?)`,
		user.Username,
		base64.StdEncoding.EncodeToString(user.Salt),
		base64.StdEncoding.EncodeToString(user.PublicKey),
		base64.StdEncoding.EncodeToString(user.PrivateKeyEncrypted),
		user.Role,
	)
	if err != nil {
		return err
	}

	id, _ := res.LastInsertId()
	user.Id = int(id)
	return nil
}

func (d *DB) GetUser(id int) (user *User, err error) {
	user = &User{Id: id}

	var salt string
	var publicKey string
	var privateKeyEncrypted string

	row := d.db.QueryRow(`SELECT username, salt, public_key, private_key_encrypted, role
		FROM users WHERE id=?`, id)
	err = row.Scan(&user.Username, &salt, &publicKey, &privateKeyEncrypted, &user.Role)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	user.Salt, err = base64.StdEncoding.DecodeString(salt)
	if err != nil {
		return nil, err
	}
	user.PublicKey, err = base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return nil, err
	}
	user.PrivateKeyEncrypted, err = base64.StdEncoding.DecodeString(privateKeyEncrypted)
	if err != nil {
		return nil, err
	}

	return
}

func (d *DB) GetUserByUsername(username string) (user *User, err error) {
	user = &User{Username: username}

	var salt string
	var publicKey string
	var privateKeyEncrypted string

	row := d.db.QueryRow(`SELECT id, salt, public_key, private_key_encrypted, role
		FROM users WHERE username=?`, username)
	err = row.Scan(&user.Id, &salt, &publicKey, &privateKeyEncrypted, &user.Role)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	user.Salt, err = base64.StdEncoding.DecodeString(salt)
	if err != nil {
		return nil, err
	}
	user.PublicKey, err = base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return nil, err
	}
	user.PrivateKeyEncrypted, err = base64.StdEncoding.DecodeString(privateKeyEncrypted)
	if err != nil {
		return nil, err
	}

	return
}

// AddVault adds a new vault and a new record in vault_keys
func (d *DB) AddVault(vault *Vault, vaultKeyEncrypted []byte) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	res, err := tx.Exec(`INSERT INTO vaults (name, owner_id) VALUES (?, ?)`,
		vault.Name, vault.OwnerId)
	if err != nil {
		return err
	}
	id, _ := res.LastInsertId()
	vault.Id = int(id)

	_, err = tx.Exec(`INSERT INTO vault_keys (user_id, vault_id, vault_key_encrypted) VALUES (?, ?, ?)`,
		vault.OwnerId, vault.Id, base64.StdEncoding.EncodeToString(vaultKeyEncrypted))
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
