package db

import (
	"database/sql"
	_ "embed"
	"sync"
)

//go:embed startup.sql
var startupQuery string

func NewPool(path string) (*sql.DB, error) {
	return sync.OnceValues(func() (*sql.DB, error) {
		pool, err := sql.Open("sqlite3", path)
		if err != nil {
			return nil, err
		}

		_, err = pool.Exec(startupQuery)
		if err != nil {
			return nil, err
		}

		return pool, nil
	})()
}
