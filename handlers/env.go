package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/TaeKwonZeus/pva/data"
	"github.com/charmbracelet/log"
	"net/http"
)

type Env struct {
	Store *data.Store
	Keys  *data.Keys
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

	type entry struct {
		Title string `json:"title"`
		Url   string `json:"url"`
	}
	res := make([]entry, 0)

	for _, vault := range index.Vaults {
		res = append(res, entry{Title: vault.Name, Url: fmt.Sprintf("/vaults/%d", vault.Id)})
		for _, password := range vault.Passwords {
			res = append(res, entry{Title: password.Name, Url: fmt.Sprintf("/vaults/%d/%d", vault.Id, password.Id)})
		}
	}

	if err = json.NewEncoder(w).Encode(res); err != nil {
		log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func authenticate(w http.ResponseWriter, r *http.Request, permission data.Permission) (user *data.User, userKey []byte, ok bool) {
	user, ok = r.Context().Value("user").(*data.User)
	if !ok {
		log.Error("could not get user from context")
		w.WriteHeader(http.StatusInternalServerError)
		return nil, nil, false
	}

	userKey, ok = r.Context().Value("userKey").([]byte)
	if !ok {
		log.Warn("could not get user key from context")
		w.WriteHeader(http.StatusInternalServerError)
		return nil, nil, false
	}

	if permission >= 0 && !data.CheckPermission(user.Role, permission) {
		w.WriteHeader(http.StatusForbidden)
		return nil, nil, false
	}
	return user, userKey, true
}
