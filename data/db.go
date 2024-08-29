package data

import (
	"database/sql"
	_ "embed"
	"encoding/base64"
	"errors"
	"github.com/mattn/go-sqlite3"
	"log"
	"time"
)

//go:embed startup.sql
var startupQuery string

type db struct {
	pool *sql.DB
}

var conflictErrors = []sqlite3.ErrNoExtended{sqlite3.ErrConstraintUnique, sqlite3.ErrConstraintPrimaryKey}

func IsErrConflict(err error) bool {
	var sqlite3Err sqlite3.Error
	if !errors.As(err, &sqlite3Err) {
		return false
	}
	for _, e := range conflictErrors {
		if errors.Is(e, sqlite3Err.ExtendedCode) {
			return true
		}
	}
	return false
}

func (d *db) createUser(user *User) error {
	_, err := d.pool.Exec(
		`INSERT INTO users (username, role, salt, public_key, private_key_encrypted)
		VALUES (?, ?, ?, ?, ?)`,
		user.Username,
		user.Role,
		base64.StdEncoding.EncodeToString(user.salt),
		base64.StdEncoding.EncodeToString(user.publicKey),
		base64.StdEncoding.EncodeToString(user.privateKeyEncrypted),
	)
	return err
}

func (d *db) getUser(id int) (user *User, err error) {
	user = &User{Id: id}

	var salt string
	var publicKey string
	var privateKeyEncrypted string

	row := d.pool.QueryRow(`SELECT username, role, salt, public_key, private_key_encrypted
		FROM users WHERE id=?`, id)
	err = row.Scan(&user.Username, &user.Role, &salt, &publicKey, &privateKeyEncrypted)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	user.salt, err = base64.StdEncoding.DecodeString(salt)
	if err != nil {
		log.Println("failed to decode salt")
		return nil, err
	}
	user.publicKey, err = base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		log.Println("failed to decode public key")
		return nil, err
	}
	user.privateKeyEncrypted, err = base64.StdEncoding.DecodeString(privateKeyEncrypted)
	if err != nil {
		log.Println("failed to decode private key")
		return nil, err
	}

	return
}

func (d *db) getUserByUsername(username string) (user *User, err error) {
	user = &User{Username: username}

	var salt string
	var publicKey string
	var privateKeyEncrypted string

	row := d.pool.QueryRow(`SELECT id, role, salt, public_key, private_key_encrypted
		FROM users WHERE username=?`, username)
	err = row.Scan(&user.Id, &user.Role, &salt, &publicKey, &privateKeyEncrypted)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	user.salt, err = base64.StdEncoding.DecodeString(salt)
	if err != nil {
		log.Println("failed to decode salt")
		return nil, err
	}
	user.publicKey, err = base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		log.Println("failed to decode public key")
		return nil, err
	}
	user.privateKeyEncrypted, err = base64.StdEncoding.DecodeString(privateKeyEncrypted)
	if err != nil {
		log.Println("failed to decode private key")
		return nil, err
	}

	return
}

func (d *db) createVault(vault *Vault, vaultKeyEncrypted []byte) error {
	tx, err := d.pool.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	res, err := tx.Exec("INSERT INTO vaults (name, owner_id) VALUES (?, ?)",
		vault.Name, vault.OwnerId)
	if err != nil {
		return err
	}

	id, _ := res.LastInsertId()

	_, err = tx.Exec("INSERT INTO vault_keys (user_id, vault_id, vault_key_encrypted) VALUES (?, ?, ?)",
		vault.OwnerId, id, base64.StdEncoding.EncodeToString(vaultKeyEncrypted))
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

type vaultAndKey struct {
	vault        *Vault
	keyEncrypted []byte
}

func (d *db) getVault(id int, userId int) (vnk *vaultAndKey, err error) {
	// Get vault by id
	vault := &Vault{Id: id}
	row := d.pool.QueryRow("SELECT name, owner_id FROM vaults where id=?", id)
	err = row.Scan(&vault.Name, &vault.OwnerId)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Get vault key
	var keyString string
	row = d.pool.QueryRow("SELECT vault_key_encrypted FROM vault_keys WHERE user_id=? AND vault_id=?", userId, id)
	err = row.Scan(&keyString)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	key, err := base64.StdEncoding.DecodeString(keyString)
	if err != nil {
		log.Println("failed to decode vault key")
		return nil, err
	}

	// Get all passwords in vault
	rows, err := d.pool.Query(`SELECT id, name, description, password_encrypted, created_at, updated_at
		FROM passwords WHERE vault_id=?`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		password := new(Password)

		var createdAtTimestamp int64
		var updatedAtTimestamp int64

		if err = rows.Scan(&password.Id, &password.Name, &password.Description, &password.passwordEncrypted,
			&createdAtTimestamp, &updatedAtTimestamp); err != nil {
			return nil, err
		}

		password.CreatedAt = time.Unix(createdAtTimestamp, 0)
		password.UpdatedAt = time.Unix(updatedAtTimestamp, 0)

		vault.Passwords = append(vault.Passwords, password)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &vaultAndKey{vault: vault, keyEncrypted: key}, nil
}

// GetVaults retrieves all pairs of vaults and data keys the user with userId has access to.
func (d *db) getVaults(userId int) (vnks []*vaultAndKey, err error) {
	rows, err := d.pool.Query("SELECT vault_id FROM vault_keys WHERE user_id=?", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		if err = rows.Scan(&id); err != nil {
			return nil, err
		}

		vault, err := d.getVault(id, userId)
		if err != nil {
			return nil, err
		}
		if vault != nil {
			vnks = append(vnks, vault)
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return
}

func (d *db) deleteVault(id int) error {
	_, err := d.pool.Exec("DELETE FROM vaults WHERE id=?", id)
	return err
}

func (d *db) createPassword(password *Password, vaultId int) error {
	_, err := d.pool.Exec(
		`INSERT INTO passwords (name, description, password_encrypted, created_at, updated_at, vault_id)
		VALUES (?, ?, ?, ?, ?, ?)`,
		password.Name,
		password.Description,
		password.passwordEncrypted,
		password.CreatedAt.Unix(),
		password.UpdatedAt.Unix(),
		vaultId,
	)
	return err
}

func (d *db) deletePassword(id int) error {
	_, err := d.pool.Exec("DELETE FROM passwords WHERE id=?", id)
	return err
}
