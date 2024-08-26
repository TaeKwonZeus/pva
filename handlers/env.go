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

func serverError(w http.ResponseWriter, err error) {
	log.Println("Server failure:", err)
	http.Error(w, err.Error(), http.StatusInternalServerError)
}
