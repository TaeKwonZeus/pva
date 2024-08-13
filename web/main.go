package main

import (
	"fmt"
	"github.com/TaeKwonZeus/pva/db"
	"github.com/TaeKwonZeus/pva/handlers"
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

	pool, err := db.NewPool(path.Join(fileDirectory, dbFilename))
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	env := &handlers.Env{Pool: pool}

	log.Printf("Listening at https://localhost:%d", cfg.Port)
	err = http.ListenAndServeTLS(
		fmt.Sprintf(":%d", cfg.Port),
		path.Join(fileDirectory, certFilename),
		path.Join(fileDirectory, keyFilename),
		newRouter(env),
	)
	if err != nil {
		log.Fatal(err)
	}
}
