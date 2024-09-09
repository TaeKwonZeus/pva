package data

import (
	"log"
	"slices"
	"time"
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

type Permission int

const (
	PermissionViewPasswords Permission = iota
	PermissionManagePasswords
	PermissionViewDevices
	PermissionManageDevices
)

var permissions = map[Role][]Permission{
	RoleManager: {PermissionViewPasswords, PermissionManagePasswords, PermissionViewDevices, PermissionManageDevices},
	RoleViewer:  {PermissionViewPasswords, PermissionViewDevices},
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
	ID       int    `json:"id,omitempty"`
	Username string `json:"username"`
	Role     Role   `json:"role"`

	salt                []byte
	publicKey           []byte
	privateKeyEncrypted []byte
}

func (u *User) DeriveKey(password string) []byte {
	return deriveKey(password, u.salt)
}

type Vault struct {
	ID        int         `json:"id,omitempty"`
	Name      string      `json:"name"`
	OwnerId   int         `json:"ownerId,omitempty"`
	Passwords []*Password `json:"passwords"`
}

type Password struct {
	ID          int       `json:"id,omitempty"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Password    string    `json:"password,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`

	passwordEncrypted []byte
}

type Device struct {
	ID          int    `json:"id"`
	IP          string `json:"ip"`
	Name        string `json:"name"`
	Description string `json:"description"`
	NetworkName string `json:"networkName"`
	MAC         string `json:"mac"`
}

type Index struct {
	Vaults []*Vault `json:"vaults"`
}
