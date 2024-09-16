package main

import (
	"errors"
	"fmt"
	"github.com/TaeKwonZeus/pva/crypt"
	"github.com/TaeKwonZeus/pva/data"
	"github.com/TaeKwonZeus/pva/handlers"
	"github.com/TaeKwonZeus/pva/network"
	"github.com/charmbracelet/log"
	"io/fs"
	stdlog "log"
	"net/http"
	"path"
	"time"

	"github.com/TaeKwonZeus/pva/config"
)

const (
	dbFilename     = "db.sqlite"
	configFilename = "config.json"
	certFilename   = "cert.pem"
	keyFilename    = "key.pem"
)

type lw struct{}

func (l lw) Write(p []byte) (n int, err error) {
	log.Info(string(p))
	return
}

func main() {
	// Route logs from other packages to charm logger
	stdlog.SetFlags(0)
	stdlog.SetOutput(lw{})

	if err := setupDirectory(); err != nil {
		log.Fatal("failed to set up working directory; please run as root", "dir", directory)
	}

	cfg, err := config.NewConfig(path.Join(directory, configFilename))
	if err != nil {
		log.Fatal(err)
	}

	store, err := data.NewStore(path.Join(directory, dbFilename))
	if err != nil {
		log.Fatal("error setting up store", "err", err)
	}
	defer store.Close()

	tokenKey, err := crypt.NewAesKey()
	if err != nil {
		log.Fatal("error setting up AES key", "err", err)
	}
	env := &handlers.Env{Store: store, TokenKey: tokenKey}

	ip, err := network.OutboundIP()
	if err != nil {
		log.Fatal(err)
	}

	network.StartAutoDiscovery(time.Minute * 2)

	log.Infof("starting server on https://%s:%d", ip, cfg.Port)
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
		log.Fatal(err)
	}
}
