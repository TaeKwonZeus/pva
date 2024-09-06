package handlers

import (
	"github.com/TaeKwonZeus/pva/data"
	"github.com/charmbracelet/log"
	"net/http"
)

type Env struct {
	Store *data.Store
	Keys  *data.Keys
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
