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
)

var permissions = map[Role][]Permission{
	RoleManager: {PermissionViewPasswords, PermissionManagePasswords},
	RoleViewer:  {PermissionViewPasswords},
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
	Id       int    `json:"id"`
	Username string `json:"username"`
	Role     Role   `json:"role"`

	salt                []byte
	publicKey           []byte
	privateKeyEncrypted []byte
}

type Vault struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	OwnerId   int    `json:"ownerId"`
	Passwords []*Password
}

type Password struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Password    string    `json:"password"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`

	passwordEncrypted []byte
}
