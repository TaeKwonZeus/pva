package handlers

import (
	"encoding/json"
	"errors"
	"github.com/TaeKwonZeus/pva/data"
	"github.com/TaeKwonZeus/pva/encryption"
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
		log.Println("could not get user")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	if !data.CheckPermission(user.Role, data.PermissionManagePasswords) {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	key, err := encryption.NewAesKey()
	if err != nil {
		serverError(w, err)
		return
	}

	vaultKeyEncrypted, err := encryption.RsaEncrypt(key, user.PublicKey, []byte(body.Name))
	if err != nil {
		serverError(w, err)
		return
	}

	err = e.DB.AddVault(&data.Vault{Name: body.Name, OwnerId: user.Id}, vaultKeyEncrypted)
	if errors.Is(err, data.ErrorConflict) {
		log.Println("Unexpected conflict:", err)
		http.Error(w, "conflict while creating vault", http.StatusConflict)
		return
	}
	if err != nil {
		serverError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
