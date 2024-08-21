package handlers

import (
	"database/sql"
	"github.com/TaeKwonZeus/pva/encryption"
)

type Env struct {
	Pool *sql.DB
	Keys *encryption.Keys
}
