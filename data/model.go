package data

import (
	"github.com/TaeKwonZeus/pva/crypt"
	"log"
	"slices"
)

type Role string

const (
	// RoleAdmin Can do anything with any resource
	RoleAdmin Role = "admin"

	// RoleManager Can view and edit resources with access
	RoleManager Role = "manager"

	// RoleViewer Can view resources with access
	RoleViewer Role = "viewer"
)

type Permission string

const (
	PermissionNone            Permission = "none"
	PermissionViewPasswords              = "passwords.view"
	PermissionManagePasswords            = "passwords.manage"
	PermissionViewDevices                = "devices.view"
	PermissionManageDevices              = "devices.manage"
	PermissionViewDocuments              = "documents.view"
	PermissionManageDocuments            = "documents.manage"
)

var permissions = map[Role][]Permission{
	RoleManager: {
		PermissionViewPasswords,
		PermissionManagePasswords,
		PermissionViewDevices,
		PermissionManageDevices,
		PermissionViewDocuments,
		PermissionManageDocuments,
	},
	RoleViewer: {
		PermissionViewPasswords,
		PermissionViewDevices,
		PermissionViewDocuments,
	},
}

func CheckPermission(role Role, permission Permission) bool {
	if role == RoleAdmin {
		return true
	}
	perms, ok := permissions[role]
	if !ok {
		log.Println("Invalid role:", role)
		return false
	}
	return slices.Contains(perms, permission)
}

type User struct {
	ID       int    `json:"id,omitempty" db:"id"`
	Username string `json:"username" db:"username"`
	Role     Role   `json:"role" db:"role"`

	Salt                []byte `json:"-" db:"salt"`
	PublicKey           []byte `json:"-" db:"public_key"`
	PrivateKey          []byte `json:"-"`
	PrivateKeyEncrypted []byte `json:"-" db:"private_key_encrypted"`
}

func (u *User) DecryptPrivateKey(key []byte) (privateKey []byte, err error) {
	u.PrivateKey, err = crypt.AesDecrypt(u.PrivateKeyEncrypted, key)
	if err != nil {
		return
	}
	return u.PrivateKey, nil
}

func (u *User) SetPrivateKey(privateKey []byte) {
	u.PrivateKey = privateKey
}

func (u *User) DeriveKey(password string) []byte {
	return crypt.DeriveKey(password, u.Salt)
}

type Vault struct {
	ID        int        `json:"id,omitempty" db:"id"`
	Name      string     `json:"name" db:"name"`
	Passwords []Password `json:"passwords"`

	KeyEncrypted []byte `json:"-" db:"key_encrypted"`
}

type Password struct {
	ID          int    `json:"id,omitempty" db:"id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	Password    string `json:"password,omitempty"`

	PasswordEncrypted []byte `json:"-" db:"password_encrypted"`
}

type Device struct {
	ID          int    `json:"id" db:"id"`
	IP          string `json:"ip" db:"ip"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	//NetworkName string `json:"networkName"`
	//MAC       string `json:"mac"`
	Connected bool `json:"connected"`
}

type Document struct {
	ID          int          `json:"id" db:"id"`
	Name        string       `json:"name" db:"name"`
	Payload     string       `json:"payload"`
	FileName    string       `json:"-" db:"file_name"`
	Attachments []Attachment `json:"attachments"`

	PayloadEncrypted []byte `json:"-" db:"payload_encrypted"`
	KeyEncrypted     []byte `json:"-" db:"key_encrypted"`
}

type Attachment struct {
	ID       int    `json:"id" db:"id"`
	Name     string `json:"name" db:"name"`
	FileName string `json:"-" db:"file_name"`
}

type Index struct {
	Vaults []Vault `json:"vaults"`
}
