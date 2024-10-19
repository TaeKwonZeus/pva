package main

import (
	"errors"
	"fmt"
	"github.com/TaeKwonZeus/pva/data"
	"github.com/TaeKwonZeus/pva/handlers"
	"github.com/TaeKwonZeus/pva/network"
	"github.com/charmbracelet/log"
	"io/fs"
	stdlog "log"
	"net"
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
	// log.SetLevel(log.DebugLevel)

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

	//tokenKey, err := crypt.NewAesKey()
	//if err != nil {
	//	log.Fatal("error setting up AES key", "err", err)
	//}
	// FIXME change in prod
	tokenKey := make([]byte, 32)
	env := &handlers.Env{Store: store, TokenKey: tokenKey}

	ip, err := network.OutboundIP()
	if err != nil {
		log.Fatal(err)
	}

	mask := net.IPMask(net.ParseIP(cfg.Scan.Netmask).To4())
	if mask == nil {
		log.Fatal("invalid netmask", "mask", mask)
	}
	timeout := time.Duration(cfg.Scan.Timeout) * time.Second
	if timeout < time.Second {
		log.Fatal("timeout cannot be less than 1 second", "timeout", timeout)
	}
	interval := time.Duration(cfg.Scan.Interval) * time.Second
	if interval < time.Second {
		log.Fatal("interval cannot be less than 1 second", "interval", interval)
	}
	network.StartAutoDiscovery(mask, timeout, interval)

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
