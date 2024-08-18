package handlers

import "database/sql"

type Env struct {
	Pool       *sql.DB
	SigningKey []byte
}
