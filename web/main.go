package main

import (
	"errors"
	"fmt"
	"github.com/TaeKwonZeus/pva/db"
	"github.com/TaeKwonZeus/pva/handlers"
	"io/fs"
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
		log.Fatal("Failed to set up /var/lib/pva-server; please run as root")
	}

	cfg, err := config.NewConfig(path.Join(directory, configFilename))
	if err != nil {
		log.Fatal(err)
	}

	pool, err := db.NewPool(path.Join(directory, dbFilename))
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	env := &handlers.Env{Pool: pool}

	log.Printf("Listening at https://localhost:%d", cfg.Port)
	err = http.ListenAndServeTLS(
		fmt.Sprintf(":%d", cfg.Port),
		path.Join(directory, certFilename),
		path.Join(directory, keyFilename),
		newRouter(env),
	)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			log.Fatalf("Missing %s or %s. Please provide a valid TLS certificate and private key in %s.",
				certFilename, keyFilename, directory)
		}
		log.Fatal(err)
	}
}
