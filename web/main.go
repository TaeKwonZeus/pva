package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"path"

	"github.com/TaeKwonZeus/pva/config"
)

//go:embed frontend/dist
var frontendEmbed embed.FS

const (
	dbFilename     = "db.sqlite"
	configFilename = "config.json"
	certFilename   = "cert.pem"
	keyFilename    = "key.pem"
)

func main() {
	log.SetFlags(log.Ldate | log.Lshortfile)

	if err := setupDirectory(); err != nil {
		log.Fatalln(err)
	}

	config, err := config.NewConfig(path.Join(fileDirectory, configFilename))
	if err != nil {
		log.Fatalln(err)
	}

	http.Handle("/api", apiServeMux())
	http.Handle("/", frontendServeMux())

	log.Printf("Listening at https://localhost:%d", config.Port)
	err = http.ListenAndServeTLS(
		fmt.Sprintf(":%d", config.Port),
		path.Join(fileDirectory, certFilename),
		path.Join(fileDirectory, keyFilename),
		http.DefaultServeMux,
	)
	if err != nil {
		log.Fatalln(err)
	}
}

func apiServeMux() *http.ServeMux {
	mux := http.NewServeMux()
	return mux
}

func frontendServeMux() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFileFS(w, r, frontendEmbed, "./favicon.ico")
	})
	mux.HandleFunc("GET /*", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFileFS(w, r, frontendEmbed, "./index.html")
	})
	mux.Handle("GET /assets/*", http.FileServerFS(frontendEmbed))

	return mux
}
