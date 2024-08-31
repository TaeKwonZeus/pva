package handlers

import (
	"github.com/TaeKwonZeus/pva/data"
	"log"
	"net/http"
)

type Env struct {
	Store *data.Store
	Keys  *data.Keys
}

func authenticate(w http.ResponseWriter, r *http.Request, permission data.Permission) (user *data.User, userKey []byte, ok bool) {
	user, ok = r.Context().Value("user").(*data.User)
	if !ok {
		log.Println("could not get user")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return nil, nil, false
	}

	userKey, ok = r.Context().Value("userKey").([]byte)
	if !ok {
		log.Println("could not get user key")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return nil, nil, false
	}

	if !data.CheckPermission(user.Role, permission) {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return nil, nil, false
	}
	return user, userKey, true
}
