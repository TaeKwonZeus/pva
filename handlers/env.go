package handlers

import (
	"github.com/TaeKwonZeus/pva/data"
)

type Env struct {
	Store    *data.Store
	TokenKey []byte
}
