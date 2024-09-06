package handlers

import (
	"encoding/json"
	"fmt"
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
	user, _, ok := authenticate(w, r, -1)
	if !ok {
		return
	}

	index, err := e.Store.GetIndex(user.Id)
	if err != nil {
		log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res := make([]entry, 0)

	for _, vault := range index.Vaults {
		res = append(res, entry{
			Title: vault.Name,
			Url:   fmt.Sprintf("/passwords?vault=%d", vault.Id),
			Type:  entryVault,
		})
		for _, password := range vault.Passwords {
			res = append(res, entry{
				Title: password.Name,
				Url:   fmt.Sprintf("/passwords?vault=%d?password=%d", vault.Id, password.Id),
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
