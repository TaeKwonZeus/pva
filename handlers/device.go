package handlers

import (
	"encoding/json"
	"github.com/TaeKwonZeus/pva/data"
	"github.com/charmbracelet/log"
	"net"
	"net/http"
	"strconv"
)

func (e *Env) NewDeviceHandler(w http.ResponseWriter, r *http.Request) {
	_, ok := authenticateNoKey(w, r, data.PermissionManageDevices)
	if !ok {
		return
	}

	var body data.Device
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if net.ParseIP(body.IP) == nil {
		http.Error(w, "invalid IP address", http.StatusBadRequest)
		return
	}

	err := e.Store.CreateDevice(&body)
	if data.IsErrConflict(err) {
		http.Error(w, "device already exists", http.StatusConflict)
		return
	}
	if err != nil {
		log.Error("error creating device", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (e *Env) GetDevicesHandler(w http.ResponseWriter, r *http.Request) {
	_, ok := authenticateNoKey(w, r, data.PermissionViewDevices)
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

func (e *Env) UpdateDeviceHandler(w http.ResponseWriter, r *http.Request) {
	_, ok := authenticateNoKey(w, r, data.PermissionManageDevices)
	if !ok {
		return
	}

	var body data.Device
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if net.ParseIP(body.IP) == nil {
		http.Error(w, "invalid IP address", http.StatusBadRequest)
		return
	}

	if err := e.Store.UpdateDevice(&body); err != nil {
		log.Error("error updating device", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (e *Env) DeleteDeviceHandler(w http.ResponseWriter, r *http.Request) {
	_, ok := authenticateNoKey(w, r, data.PermissionManageDevices)
	if !ok {
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "invalid device id", http.StatusBadRequest)
		return
	}

	err = e.Store.DeleteDevice(id)
	if data.IsErrNotFound(err) {
		http.Error(w, "device not found", http.StatusNotFound)
		return
	}
	if err != nil {
		log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
