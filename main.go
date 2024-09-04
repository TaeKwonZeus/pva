package main

import (
	"errors"
	"fmt"
	"github.com/TaeKwonZeus/pva/data"
	"github.com/TaeKwonZeus/pva/handlers"
	"github.com/charmbracelet/log"
	"io/fs"
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
	if err := setupDirectory(); err != nil {
		log.Fatal("failed to set up working directory; please run as root", "dir", directory)
	}

	cfg, err := config.NewConfig(path.Join(directory, configFilename))
	if err != nil {
		log.Fatal(err)
	}

	keys, err := data.NewKeys()
	if err != nil {
		log.Fatal(err)
	}
	defer keys.Erase()

	store, err := data.NewStore(path.Join(directory, dbFilename), keys.PasswordKey())
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()

	env := &handlers.Env{Store: store, Keys: keys}

	log.Info("starting server", "port", cfg.Port)
	err = http.ListenAndServeTLS(
		fmt.Sprintf(":%d", cfg.Port),
		path.Join(directory, certFilename),
		path.Join(directory, keyFilename),
		newRouter(env),
	)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			log.Fatal("missing cert or key",
				"cert", certFilename, "key", keyFilename, "dir", directory)
		}
	}
}
