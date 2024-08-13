package main

import (
	"fmt"
	"log"
	"net/http"
	"path"

	"github.com/TaeKwonZeus/pva/config"
)

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

	cfg, err := config.NewConfig(path.Join(fileDirectory, configFilename))
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Listening at https://localhost:%d", cfg.Port)
	err = http.ListenAndServeTLS(
		fmt.Sprintf(":%d", cfg.Port),
		path.Join(fileDirectory, certFilename),
		path.Join(fileDirectory, keyFilename),
		newRouter(),
	)
	if err != nil {
		log.Fatal(err)
	}
}
