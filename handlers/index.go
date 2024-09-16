package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/TaeKwonZeus/pva/data"
	"github.com/charmbracelet/log"
	"net/http"
)

type entryType string

const (
	entryVault    entryType = "vault"
	entryPassword entryType = "password"
)

type entry struct {
	Title string    `json:"title"`
	Url   string    `json:"url"`
	Type  entryType `json:"type"`
}

func (e *Env) GetIndexHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := authenticate(w, r, data.PermissionNone)
	if !ok {
		return
	}

	index, err := e.Store.GetIndex(user.ID)
	if err != nil {
		log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res := make([]entry, 0)

	for _, vault := range index.Vaults {
		res = append(res, entry{
			Title: vault.Name,
			Url:   fmt.Sprintf("/passwords?vault=%d", vault.ID),
			Type:  entryVault,
		})
		for _, password := range vault.Passwords {
			res = append(res, entry{
				Title: password.Name,
				Url:   fmt.Sprintf("/passwords?vault=%d?password=%d", vault.ID, password.ID),
				Type:  entryPassword,
			})
		}
	}

	if err = json.NewEncoder(w).Encode(res); err != nil {
		log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
