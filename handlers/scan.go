package handlers

import (
	"encoding/json"
	"github.com/TaeKwonZeus/pva/data"
	"github.com/charmbracelet/log"
	"net/http"
)

func (e *Env) GetDevicesHandler(w http.ResponseWriter, r *http.Request) {
	_, _, ok := authenticate(w, r, data.PermissionViewDevices)
	if !ok {
		return
	}

	devices, err := e.Store.GetDevices()
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = json.NewEncoder(w).Encode(devices); err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
