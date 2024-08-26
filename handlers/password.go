package handlers

import (
	"encoding/json"
	"github.com/TaeKwonZeus/pva/data"
	"log"
	"net/http"
)

func (e *Env) NewVaultHandler(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, ok := r.Context().Value("user").(*data.User)
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	if !data.CheckPermission(user.Role, data.PermissionManagePasswords) {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	_, err := e.Store.CreateVault(&data.Vault{
		Name: body.Name,
	}, user)
	if data.IsErrConflict(err) {
		http.Error(w, "Vault already exists", http.StatusConflict)
		return
	}
	if err != nil {
		serverError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (e *Env) GetVaultsHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(*data.User)
	if !ok {
		log.Println("could not get user")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	password, ok := r.Context().Value("password").(string)
	if !ok {
		log.Println("could not get password")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	if !data.CheckPermission(user.Role, data.PermissionViewPasswords) {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	vaults, err := e.Store.GetVaults(user, password)
	if err != nil {
		serverError(w, err)
		return
	}

	err = json.NewEncoder(w).Encode(vaults)
	if err != nil {
		serverError(w, err)
		return
	}
}
