package network

import (
	"github.com/charmbracelet/log"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"net"
	"os"
	"sync"
	"time"
)

var (
	Timeout time.Duration
)

func GetOutboundIP() (ip net.IP, err error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil, err
	}
	ip = conn.LocalAddr().(*net.UDPAddr).IP
	return
}

type Device struct {
	IP   net.IP
	Name string
	MAC  net.HardwareAddr
}

func Scan() (devices []Device, err error) {
	conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	for _, ip := range localIPs() {
		go sendICMP(conn, ip)
	}

	done := make(chan error)
	res := newResults()

	go func() {
		for {
			go recvICMP(conn, done, res.c)
		}
	}()

	select {
	case err = <-done:
		return nil, err
	case <-time.After(Timeout):
		return res.get(), nil
	}
}

func localIPs() []net.IP {
	var ips []net.IP
	for i := byte(1); i < 254; i++ {
		ips = append(ips, net.IPv4(192, 168, 0, i))
	}
	return ips
}

func sendICMP(conn *icmp.PacketConn, ip net.IP) {
	msg, _ := (&icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   os.Getpid() & 0xffff,
			Seq:  1,
			Data: []byte("ping"),
		},
	}).Marshal(nil)

	if _, err := conn.WriteTo(msg, &net.IPAddr{IP: ip}); err != nil {
		log.Error("send error", "ip", ip.String(), "err", err)
	}
}

func recvICMP(conn *icmp.PacketConn, done chan<- error, res chan<- Device) {
	rb := make([]byte, 1024)
	n, peer, err := conn.ReadFrom(rb)
	if err != nil {
		done <- err
		return
	}
	peerIP := net.ParseIP(peer.String())

}

type results struct {
	devices map[string]Device
	mu      sync.RWMutex
	c       chan<- Device
}

func (r *results) get() []Device {
	r.mu.RLock()
	defer r.mu.RUnlock()

	devices := make([]Device, 0, len(r.devices))
	for _, d := range r.devices {
		devices = append(devices, d)
	}
	return devices
}

func newResults() *results {
	c := make(chan Device)
	r := &results{
		devices: make(map[string]Device),
		c:       c,
	}
	go func() {
		for device := range c {
			ip := device.IP.String()

			r.mu.RLock()
			if _, ok := r.devices[ip]; ok {
				log.Info("ICMP double receive", "ip", ip)
				r.mu.RUnlock()
				continue
			}
			r.mu.RUnlock()

			r.mu.Lock()
			r.devices[ip] = device
			r.mu.Unlock()
		}
	}()
	return r
}
