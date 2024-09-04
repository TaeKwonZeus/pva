package handlers

import (
	"encoding/json"
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

	if err = json.NewEncoder(w).Encode(index); err != nil {
		log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func authenticate(w http.ResponseWriter, r *http.Request, permission data.Permission) (user *data.User, userKey []byte, ok bool) {
	user, ok = r.Context().Value("user").(*data.User)
	if !ok {
		log.Warn("could not get user")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return nil, nil, false
	}

	userKey, ok = r.Context().Value("userKey").([]byte)
	if !ok {
		log.Warn("could not get user key")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return nil, nil, false
	}

	if permission >= 0 && !data.CheckPermission(user.Role, permission) {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return nil, nil, false
	}
	return user, userKey, true
}
