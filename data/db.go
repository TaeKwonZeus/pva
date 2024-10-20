package data

import (
	"database/sql"
	_ "embed"
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/mattn/go-sqlite3"
)

//go:embed startup.sql
var startupQuery string

type db struct {
	pool *sqlx.DB
}

func IsErrConflict(err error) bool {
	var sqlite3Err sqlite3.Error
	return errors.As(err, &sqlite3Err) && (errors.Is(sqlite3Err.ExtendedCode, sqlite3.ErrConstraintUnique) ||
		errors.Is(sqlite3Err.ExtendedCode, sqlite3.ErrConstraintPrimaryKey))
}

func IsErrNotFound(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}

func (d *db) getIndex(id int) (index Index, err error) {
	vaults, err := d.getVaults(id)
	if err != nil {
		return
	}

	index.Vaults = vaults

	return
}

func (d *db) getUserCount() (n int, err error) {
	row := d.pool.QueryRow("SELECT COUNT(*) FROM users")
	err = row.Scan(&n)
	return
}

func (d *db) createUser(user User) (id int, err error) {
	res, err := d.pool.NamedExec(
		`INSERT INTO users (username, role, salt, public_key, private_key_encrypted)
		VALUES (:username, :role, :salt, :public_key, :private_key_encrypted)`, user)
	if err != nil {
		return 0, err
	}
	i, err := res.LastInsertId()
	return int(i), err
}

func (d *db) getUser(id int) (user User, err error) {
	user.ID = id
	err = d.pool.Get(&user, "SELECT * FROM users WHERE id=?", id)
	return
}

func (d *db) getUserByUsername(username string) (user User, err error) {
	user.Username = username
	err = d.pool.Get(&user, `SELECT * FROM users WHERE username=?`, username)
	return
}

func (d *db) getAdmins() (users []User, err error) {
	users = []User{}
	err = d.pool.Select(&users, "SELECT * FROM users WHERE role=?", RoleAdmin)
	return
}

func (d *db) createVault(vault Vault) (id int, err error) {
	res, err := d.pool.NamedExec("INSERT INTO vaults (name) VALUES (:name)", vault)
	if err != nil {
		return 0, err
	}
	i, err := res.LastInsertId()
	return int(i), err
}

type vaultKey struct {
	UserId       int    `db:"user_id"`
	VaultId      int    `db:"vault_id"`
	KeyEncrypted []byte `db:"key_encrypted"`
}

func (d *db) createVaultKeys(keys ...vaultKey) error {
	tx, err := d.pool.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.NamedExec(`INSERT INTO vault_keys (user_id, vault_id, key_encrypted)
		VALUES (:user_id, :vault_id, :key_encrypted) ON CONFLICT DO NOTHING`, keys)

	return tx.Commit()
}

func (d *db) getVaultKey(id, userId int) (key []byte, err error) {
	err = d.pool.Get(&key, "SELECT key_encrypted FROM vault_keys WHERE user_id=? AND vault_id=?", userId, id)
	return
}

func (d *db) getVault(id, userId int) (vault Vault, err error) {
	err = d.pool.Get(&vault, `SELECT v.*, vk.key_encrypted FROM vaults v
        INNER JOIN vault_keys vk ON v.id = vk.vault_id WHERE user_id=? AND vault_id=?`, userId, id)
	if err != nil {
		return
	}
	vault.Passwords, err = d.getPasswords(id)
	return
}

// GetVaults retrieves all pairs of vaults and data keys the user with userId has access to.
func (d *db) getVaults(userId int) (vaults []Vault, err error) {
	vaults = []Vault{}
	err = d.pool.Select(&vaults, `SELECT v.*, vk.key_encrypted FROM vaults v
        INNER JOIN vault_keys vk ON v.id = vk.vault_id WHERE user_id=?`, userId)
	if err != nil {
		return
	}
	for i := range vaults {
		vaults[i].Passwords, err = d.getPasswords(vaults[i].ID)
		if err != nil {
			return
		}
	}
	return
}

func (d *db) updateVault(vault Vault) error {
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

func (d *db) createPassword(password Password, vaultId int) (id int, err error) {
	res, err := d.pool.Exec(
		`INSERT INTO passwords (name, description, password_encrypted, vault_id)
		VALUES (?, ?, ?, ?)`,
		password.Name,
		password.Description,
		password.PasswordEncrypted,
		vaultId,
	)
	if err != nil {
		return 0, err
	}
	i, err := res.LastInsertId()
	return int(i), err
}

func (d *db) getPasswords(vaultId int) (passwords []Password, err error) {
	passwords = []Password{}
	err = d.pool.Select(&passwords, "SELECT id, name, description, password_encrypted FROM passwords WHERE vault_id=?", vaultId)
	return
}

func (d *db) updatePassword(password Password) error {
	tx, err := d.pool.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if password.Name != "" {
		_, err = tx.Exec("UPDATE passwords SET name=? WHERE id=?", password.Name, password.ID)
		if err != nil {
			return err
		}
	}
	if password.Description != "" {
		_, err = tx.Exec("UPDATE passwords SET description=? WHERE id=?", password.Description, password.ID)
		if err != nil {
			return err
		}
	}
	if password.PasswordEncrypted != nil {
		_, err = tx.Exec("UPDATE passwords SET password_encrypted=? WHERE id=?", password.PasswordEncrypted, password.ID)
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

func (d *db) createDevice(device Device) (id int, err error) {
	res, err := d.pool.NamedExec("INSERT INTO devices (ip, name, description) VALUES (:ip, :name, :description)",
		device)
	if err != nil {
		return 0, err
	}
	i, err := res.LastInsertId()
	return int(i), err
}

func (d *db) getDevices() (devices []Device, err error) {
	devices = []Device{}
	err = d.pool.Select(&devices, "SELECT id, ip, name, description FROM devices")
	return
}

func (d *db) updateDevice(device Device) error {
	_, err := d.pool.Exec("UPDATE devices SET ip=?, name=?, description=? WHERE id=?", device.IP, device.Name, device.Description, device.ID)
	return err
}

func (d *db) deleteDevice(id int) error {
	_, err := d.pool.Exec("DELETE FROM devices WHERE id=?", id)
	return err
}

func (d *db) createDocument(document Document) (id int, err error) {
	res, err := d.pool.Exec("INSERT INTO documents (name, payload_encrypted) VALUES (?, ?)",
		document.Name, document.PayloadEncrypted)
	if err != nil {
		return 0, err
	}
	i, err := res.LastInsertId()
	return int(i), err
}

type documentKey struct {
	UserId       int    `db:"user_id"`
	DocumentId   int    `db:"vault_id"`
	KeyEncrypted []byte `db:"key_encrypted"`
}

func (d *db) createDocumentKeys(keys ...documentKey) error {
	tx, err := d.pool.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.NamedExec(`INSERT INTO document_keys (user_id, document_id, key_encrypted)
		VALUES (:user_id, :document_id, :key_encrypted)`, keys)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (d *db) getDocuments(userId int) (docs []Document, err error) {
	docs = []Document{}
	err = d.pool.Select(&docs, `SELECT d.*, dk.key_encrypted FROM documents d
        INNER JOIN document_keys dk on d.id = dk.document_id WHERE dk.user_id=?`, userId)

	// TODO add attachments
	return
}

func (d *db) getDocument(id, userId int) (doc Document, err error) {
	err = d.pool.Get(&doc, `SELECT d.*, dk.key_encrypted FROM documents d
        INNER JOIN document_keys dk on d.id = dk.document_id WHERE user_id=? AND document_id=?`, userId, id)

	// TODO add attachments
	return doc, err
}

func (d *db) updateDocument(id, userId int) error {
	// TODO implement
	return nil
}
