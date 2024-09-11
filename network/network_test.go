package network

import (
	"github.com/charmbracelet/log"
	"testing"
)

func TestScan(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	log.Info("Started")
	devices, err := Scan()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("devices:", devices)
}
