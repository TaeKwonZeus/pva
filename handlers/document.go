package handlers

import (
	"encoding/json"
	"github.com/TaeKwonZeus/pva/data"
	"net/http"
)

func (e *Env) GetDocumentHandler(w http.ResponseWriter, r *http.Request) {
	_, ok := authenticate(w, r, data.PermissionViewDocuments)
	if !ok {
		return
	}
}

func (e *Env) NewDocumentHandler(w http.ResponseWriter, r *http.Request) {
	_, ok := authenticate(w, r, data.PermissionManageDocuments)
	if !ok {
		return
	}

	var body struct {
		Name string `json:"name"`
		Body string `json:"body"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

}

func (e *Env) UpdateDocumentHandler(w http.ResponseWriter, r *http.Request) {
	_, ok := authenticate(w, r, data.PermissionManageDocuments)
	if !ok {
		return
	}
}

func (e *Env) DeleteDocumentHandler(w http.ResponseWriter, r *http.Request) {
	_, ok := authenticate(w, r, data.PermissionManageDocuments)
	if !ok {
		return
	}
}

func (e *Env) ShareDocumentHandler(w http.ResponseWriter, r *http.Request) {
	_, ok := authenticate(w, r, data.PermissionManageDocuments)
	if !ok {
		return
	}
}
