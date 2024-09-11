package network

import (
	"github.com/charmbracelet/log"
	"net"
	"testing"
	"time"
)

func TestScan(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	log.Info("Started")
	devices, err := Scan(WithTimeout(time.Second*3), WithMask(net.IPv4Mask(255, 255, 255, 0)))
	if err != nil {
		t.Fatal(err)
	}
	t.Log("devices:", devices)
}
