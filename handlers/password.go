package handlers

import (
	"encoding/json"
	"github.com/TaeKwonZeus/pva/data"
	"github.com/charmbracelet/log"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func (e *Env) NewVaultHandler(w http.ResponseWriter, r *http.Request) {
	var body data.Vault
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, ok := authenticate(w, r, data.PermissionManagePasswords)
	if !ok {
		return
	}

	err := e.Store.CreateVault(body, user)
	if data.IsErrConflict(err) {
		http.Error(w, "vault already exists", http.StatusConflict)
		return
	}
	if err != nil {
		log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (e *Env) GetVaultsHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := authenticate(w, r, data.PermissionViewPasswords)
	if !ok {
		return
	}

	vaults, err := e.Store.GetVaults(user)
	if err != nil {
		log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(vaults)
	if err != nil {
		log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (e *Env) UpdateVaultHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid vault id", http.StatusBadRequest)
		return
	}

	var body data.Vault
	if err = json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	body.ID = id

	user, ok := authenticate(w, r, data.PermissionManagePasswords)
	if !ok {
		return
	}

	if !e.Store.CheckVaultOwnership(id, user) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	err = e.Store.UpdateVault(body)
	if data.IsErrNotFound(err) {
		http.Error(w, "vault not found", http.StatusNotFound)
		return
	}
	if data.IsErrConflict(err) {
		http.Error(w, "vault already exists", http.StatusConflict)
		return
	}
	if err != nil {
		log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (e *Env) DeleteVaultHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid vault id", http.StatusBadRequest)
		return
	}

	user, ok := authenticate(w, r, data.PermissionManagePasswords)
	if !ok {
		return
	}

	if !e.Store.CheckVaultOwnership(id, user) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	err = e.Store.DeleteVault(id)
	if err != nil {
		log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (e *Env) ShareVaultHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid vault id", http.StatusBadRequest)
		return
	}

	targetUsername := r.URL.Query().Get("target")

	user, ok := authenticate(w, r, data.PermissionManagePasswords)
	if !ok {
		return
	}

	if !e.Store.CheckVaultOwnership(id, user) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	target, err := e.Store.GetUserByUsername(targetUsername)
	if data.IsErrNotFound(err) {
		http.Error(w, "target not found", http.StatusNotFound)
		return
	}
	if err != nil {
		log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = e.Store.ShareVault(id, target, user)
	if err != nil {
		log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (e *Env) NewPasswordHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid vault id", http.StatusBadRequest)
		return
	}

	var body data.Password
	if err = json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, ok := authenticate(w, r, data.PermissionManagePasswords)
	if !ok {
		return
	}

	if !e.Store.CheckVaultOwnership(id, user) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	err = e.Store.CreatePassword(body, id, user)
	if data.IsErrConflict(err) {
		http.Error(w, "password already exists in the same vault", http.StatusConflict)
		return
	}
	if err != nil {
		log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (e *Env) UpdatePasswordHandler(w http.ResponseWriter, r *http.Request) {
	vaultId, err := strconv.Atoi(chi.URLParam(r, "vaultId"))
	if err != nil {
		http.Error(w, "invalid vault id", http.StatusBadRequest)
		return
	}
	passwordId, err := strconv.Atoi(chi.URLParam(r, "passwordId"))
	if err != nil {
		http.Error(w, "invalid password id", http.StatusBadRequest)
		return
	}

	var body data.Password
	if err = json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	body.ID = passwordId

	user, ok := authenticate(w, r, data.PermissionManagePasswords)
	if !ok {
		return
	}

	if !e.Store.CheckVaultOwnership(vaultId, user) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	err = e.Store.UpdatePassword(body, vaultId, user)
	if data.IsErrConflict(err) {
		http.Error(w, "password already exists in the same vault", http.StatusConflict)
		return
	}
	if err != nil {
		log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (e *Env) DeletePasswordHandler(w http.ResponseWriter, r *http.Request) {
	vaultId, err := strconv.Atoi(chi.URLParam(r, "vaultId"))
	if err != nil {
		http.Error(w, "invalid vault id", http.StatusBadRequest)
		return
	}
	passwordId, err := strconv.Atoi(chi.URLParam(r, "passwordId"))
	if err != nil {
		http.Error(w, "invalid password id", http.StatusBadRequest)
		return
	}

	user, ok := authenticate(w, r, data.PermissionManagePasswords)
	if !ok {
		return
	}

	if !e.Store.CheckVaultOwnership(vaultId, user) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	err = e.Store.DeletePassword(passwordId)
	if err != nil {
		log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
