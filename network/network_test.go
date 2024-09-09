package network

import (
	"github.com/charmbracelet/log"
	"testing"
	"time"
)

func TestScan(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	devices, err := Scan(WithTimeout(time.Second))
	if err != nil {
		t.Fatal(err)
	}
	t.Log("devices:", devices)
}
