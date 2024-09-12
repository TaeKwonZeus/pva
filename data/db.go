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

func IsErrConflict(err error) bool {
	var sqlite3Err sqlite3.Error
	return errors.As(err, &sqlite3Err) && (errors.Is(sqlite3Err.ExtendedCode, sqlite3.ErrConstraintUnique) ||
		errors.Is(sqlite3Err.ExtendedCode, sqlite3.ErrConstraintPrimaryKey))
}

func IsErrNotFound(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}

func (d *db) getIndex(id int) (index *Index, err error) {
	index = new(Index)

	vnks, err := d.getVaults(id)
	if err != nil {
		return nil, err
	}

	index.Vaults = make([]*Vault, len(vnks))
	for i, vnk := range vnks {
		index.Vaults[i] = vnk.vault
	}

	return
}

func (d *db) getUserCount() (n int, err error) {
	row := d.pool.QueryRow("SELECT COUNT(*) FROM users")
	err = row.Scan(&n)
	return
}

func (d *db) createUser(user *User) error {
	res, err := d.pool.Exec(
		`INSERT INTO users (username, role, salt, public_key, private_key_encrypted)
		VALUES (?, ?, ?, ?, ?)`,
		user.Username,
		user.Role,
		base64.StdEncoding.EncodeToString(user.salt),
		base64.StdEncoding.EncodeToString(user.publicKey),
		base64.StdEncoding.EncodeToString(user.privateKeyEncrypted),
	)
	if err != nil {
		return err
	}
	id, _ := res.LastInsertId()
	user.ID = int(id)
	return nil
}

func (d *db) getUser(id int) (user *User, err error) {
	user = &User{ID: id}

	var salt string
	var publicKey string
	var privateKeyEncrypted string

	row := d.pool.QueryRow(`SELECT username, role, salt, public_key, private_key_encrypted
		FROM users WHERE id=?`, id)
	err = row.Scan(&user.Username, &user.Role, &salt, &publicKey, &privateKeyEncrypted)
	if err != nil {
		return nil, err
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
	err = row.Scan(&user.ID, &user.Role, &salt, &publicKey, &privateKeyEncrypted)
	if err != nil {
		return nil, err
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

func (d *db) getAdmins() (users []*User, err error) {
	rows, err := d.pool.Query(`SELECT id, username, salt, public_key, private_key_encrypted FROM users
        WHERE role=?`, RoleAdmin)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		user := User{Role: RoleAdmin}

		var salt string
		var publicKey string
		var privateKeyEncrypted string

		if err = rows.Scan(&user.ID, &user.Username, &salt, &publicKey, &privateKeyEncrypted); err != nil {
			return nil, err
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

		users = append(users, &user)
	}
	if err = rows.Err(); err != nil {
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
	vault.ID = int(id)

	_, err = tx.Exec("INSERT INTO vault_keys (user_id, vault_id, vault_key_encrypted) VALUES (?, ?, ?)",
		vault.OwnerId, id, base64.StdEncoding.EncodeToString(vaultKeyEncrypted))
	if err != nil {
		return err
	}

	return tx.Commit()
}

type vaultKey struct {
	userId       int
	vaultId      int
	keyEncrypted []byte
}

func (d *db) createVaultKeys(keys ...*vaultKey) error {
	tx, err := d.pool.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`INSERT INTO vault_keys (user_id, vault_id, vault_key_encrypted) VALUES (?, ?, ?) ON CONFLICT DO NOTHING`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, key := range keys {
		_, err = stmt.Exec(key.userId, key.vaultId, base64.StdEncoding.EncodeToString(key.keyEncrypted))
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

type vaultAndKey struct {
	vault        *Vault
	keyEncrypted []byte
}

func (d *db) getVaultKey(id, userId int) (key []byte, err error) {
	var keyString string
	row := d.pool.QueryRow("SELECT vault_key_encrypted FROM vault_keys WHERE user_id=? AND vault_id=?", userId, id)
	err = row.Scan(&keyString)
	if err != nil {
		return nil, err
	}

	key, err = base64.StdEncoding.DecodeString(keyString)
	if err != nil {
		log.Println("failed to decode vault key")
		return nil, err
	}

	return
}

func (d *db) getVault(id, userId int) (vnk *vaultAndKey, err error) {
	// Get vault by id
	vault := &Vault{ID: id}
	row := d.pool.QueryRow("SELECT name, owner_id FROM vaults where id=?", id)
	err = row.Scan(&vault.Name, &vault.OwnerId)
	if err != nil {
		return nil, err
	}

	// Get vault key
	key, err := d.getVaultKey(vault.ID, userId)
	if err != nil {
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

		if err = rows.Scan(&password.ID, &password.Name, &password.Description, &password.passwordEncrypted,
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

func (d *db) updateVault(vault *Vault) error {
	if vault.Name != "" {
		_, err := d.pool.Exec("UPDATE vaults SET name=? WHERE id=?", vault.Name, vault.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *db) deleteVault(id int) error {
	// All passwords get cascade deleted by sqlite
	_, err := d.pool.Exec("DELETE FROM vaults WHERE id=?", id)
	return err
}

func (d *db) createPassword(password *Password, vaultId int) error {
	res, err := d.pool.Exec(
		`INSERT INTO passwords (name, description, password_encrypted, created_at, updated_at, vault_id)
		VALUES (?, ?, ?, ?, ?, ?)`,
		password.Name,
		password.Description,
		password.passwordEncrypted,
		password.CreatedAt.Unix(),
		password.UpdatedAt.Unix(),
		vaultId,
	)
	if err != nil {
		return err
	}
	id, _ := res.LastInsertId()
	password.ID = int(id)
	return nil
}

func (d *db) updatePassword(password *Password) error {
	tx, err := d.pool.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var updated bool

	if password.Name != "" {
		updated = true
		_, err = tx.Exec("UPDATE passwords SET name=? WHERE id=?", password.Name, password.ID)
		if err != nil {
			return err
		}
	}
	if password.Description != "" {
		updated = true
		_, err = tx.Exec("UPDATE passwords SET description=? WHERE id=?", password.Description, password.ID)
		if err != nil {
			return err
		}
	}
	if password.passwordEncrypted != nil {
		updated = true
		_, err = tx.Exec("UPDATE passwords SET password_encrypted=? WHERE id=?", password.passwordEncrypted, password.ID)
		if err != nil {
			return err
		}
	}
	if updated {
		_, err = tx.Exec("UPDATE passwords SET updated_at=? WHERE id=?", password.UpdatedAt.Unix(), password.ID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (d *db) deletePassword(id int) error {
	_, err := d.pool.Exec("DELETE FROM passwords WHERE id=?", id)
	return err
}

func (d *db) createDevice(device *Device) error {
	_, err := d.pool.Exec("INSERT INTO devices (ip, name, description) VALUES (?, ?, ?)",
		device.IP, device.Name, device.Description)
	return err
}

func (d *db) getDevices() (devices []*Device, err error) {
	rows, err := d.pool.Query("SELECT id, ip, name, description FROM devices")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var device Device
		if err = rows.Scan(&device.ID, &device.IP, &device.Name, &device.Description); err != nil {
			return nil, err
		}
		devices = append(devices, &device)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return
}

func (d *db) updateDevice(device *Device) error {
	_, err := d.pool.Exec("UPDATE devices SET ip=?, name=?, description=? WHERE id=?", device.IP, device.Name, device.Description, device.ID)
	return err
}

func (d *db) deleteDevice(id int) error {
	_, err := d.pool.Exec("DELETE FROM devices WHERE id=?", id)
	return err
}
