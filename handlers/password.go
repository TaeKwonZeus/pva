package handlers

import (
	"encoding/json"
	"github.com/TaeKwonZeus/pva/data"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"strconv"
)

func (e *Env) NewVaultHandler(w http.ResponseWriter, r *http.Request) {
	var body data.Vault
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, _, ok := authenticate(w, r, data.PermissionManagePasswords)
	if !ok {
		return
	}

	err := e.Store.CreateVault(&body, user)
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
	user, userKey, ok := authenticate(w, r, data.PermissionViewPasswords)
	if !ok {
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

func (e *Env) UpdateVaultHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid vault id", http.StatusBadRequest)
		return
	}

	var body data.Vault
	if err = json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	body.Id = id

	user, userKey, ok := authenticate(w, r, data.PermissionManagePasswords)
	if !ok {
		return
	}

	if !e.Store.CheckVaultOwnership(id, user, userKey) {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	err = e.Store.UpdateVault(&body)
	if data.IsErrNotFound(err) {
		http.Error(w, "Vault not found", http.StatusNotFound)
		return
	}
	if data.IsErrConflict(err) {
		http.Error(w, "Vault already exists", http.StatusConflict)
		return
	}
	if err != nil {
		log.Println("Server failure:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (e *Env) DeleteVaultHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid vault id", http.StatusBadRequest)
		return
	}

	user, userKey, ok := authenticate(w, r, data.PermissionManagePasswords)
	if !ok {
		return
	}

	if !e.Store.CheckVaultOwnership(id, user, userKey) {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	err = e.Store.DeleteVault(id)
	if err != nil {
		log.Println("Server failure:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (e *Env) ShareVaultHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid vault id", http.StatusBadRequest)
		return
	}

	targetUsername := r.URL.Query().Get("target")

	user, userKey, ok := authenticate(w, r, data.PermissionManagePasswords)
	if !ok {
		return
	}

	if !e.Store.CheckVaultOwnership(id, user, userKey) {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	target, err := e.Store.GetUserByUsername(targetUsername)
	if data.IsErrNotFound(err) {
		http.Error(w, "Target not found", http.StatusNotFound)
		return
	}
	if err != nil {
		log.Println("Server failure:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = e.Store.ShareVault(id, target, user, userKey)
	if err != nil {
		log.Println("Server failure:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (e *Env) NewPasswordHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid vault id", http.StatusBadRequest)
		return
	}

	var body data.Password
	if err = json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, userKey, ok := authenticate(w, r, data.PermissionManagePasswords)
	if !ok {
		return
	}

	if !e.Store.CheckVaultOwnership(id, user, userKey) {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	err = e.Store.CreatePassword(&body, id, user, userKey)
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

func (e *Env) UpdatePasswordHandler(w http.ResponseWriter, r *http.Request) {
	vaultId, err := strconv.Atoi(chi.URLParam(r, "vaultId"))
	if err != nil {
		http.Error(w, "Invalid vault id", http.StatusBadRequest)
		return
	}
	passwordId, err := strconv.Atoi(chi.URLParam(r, "passwordId"))
	if err != nil {
		http.Error(w, "Invalid password id", http.StatusBadRequest)
		return
	}

	var body data.Password
	if err = json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	body.Id = passwordId

	user, userKey, ok := authenticate(w, r, data.PermissionManagePasswords)
	if !ok {
		return
	}

	if !e.Store.CheckVaultOwnership(vaultId, user, userKey) {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	err = e.Store.UpdatePassword(&body, vaultId, user, userKey)
	if data.IsErrConflict(err) {
		http.Error(w, "Password already exists in the same vault", http.StatusConflict)
		return
	}
	if err != nil {
		log.Println("Server failure:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (e *Env) DeletePasswordHandler(w http.ResponseWriter, r *http.Request) {
	vaultId, err := strconv.Atoi(chi.URLParam(r, "vaultId"))
	if err != nil {
		http.Error(w, "Invalid vault id", http.StatusBadRequest)
		return
	}
	passwordId, err := strconv.Atoi(chi.URLParam(r, "passwordId"))
	if err != nil {
		http.Error(w, "Invalid password id", http.StatusBadRequest)
		return
	}

	user, userKey, ok := authenticate(w, r, data.PermissionManagePasswords)
	if !ok {
		return
	}

	if !e.Store.CheckVaultOwnership(vaultId, user, userKey) {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	err = e.Store.DeletePassword(passwordId)
	if err != nil {
		log.Println("Server failure:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
