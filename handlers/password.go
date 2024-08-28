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

	err := e.Store.CreateVault(&data.Vault{
		Name: body.Name,
	}, user)
	if data.IsErrConflict(err) {
		http.Error(w, "Vault already exists", http.StatusConflict)
		return
	}
	if err != nil {
		log.Println("Server failure:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

	userKey, ok := r.Context().Value("userKey").([]byte)
	if !ok {
		log.Println("could not get user key")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	if !data.CheckPermission(user.Role, data.PermissionViewPasswords) {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	vaults, err := e.Store.GetVaults(user, userKey)
	if err != nil {
		log.Println("Server failure:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(vaults)
	if err != nil {
		log.Println("Server failure:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (e *Env) NewPasswordHandler(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Password    string `json:"password"`
		VaultId     int    `json:"vaultId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, ok := r.Context().Value("user").(*data.User)
	if !ok {
		log.Println("could not get user")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	userKey, ok := r.Context().Value("userKey").([]byte)
	if !ok {
		log.Println("could not get user key")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	if !data.CheckPermission(user.Role, data.PermissionManagePasswords) {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	err := e.Store.CreatePassword(&data.Password{
		Name:        body.Name,
		Description: body.Description,
		Password:    body.Password,
	}, body.VaultId, user, userKey)
	if data.IsErrConflict(err) {
		http.Error(w, "Password already exists in the same vault", http.StatusConflict)
		return
	}
	if err != nil {
		log.Println("Server failure:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
