package network

import (
	"encoding/base64"
	"github.com/charmbracelet/log"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/sync/errgroup"
	"net"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

var (
	Timeout time.Duration = time.Second * 5
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

	var sent atomic.Int32
	for _, ip := range localIPs() {
		go func() {
			err := sendICMP(conn, ip)
			if err != nil {
				log.Error("ICMP send error", "ip", ip, "err", err)
				return
			}
			sent.Add(1)
		}()
	}

	done := make(chan error)
	res := newResults()

	var eg errgroup.Group

	n := sent.Load()
	log.Debugf("sent %d echo packets", n)
	for range n {
		eg.Go(func() error {
			return recvICMP(conn, res.c)
		})
	}
	go func() {
		done <- eg.Wait()
	}()

	select {
	case err = <-done:
		return res.get(), err
	case <-time.After(Timeout):
		log.Warn("ICMP timeout")
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

func sendICMP(conn *icmp.PacketConn, ip net.IP) error {
	msg, _ := (&icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   os.Getpid() & 0xffff,
			Seq:  1,
			Data: []byte("ping"),
		},
	}).Marshal(nil)

	log.Debug("sending ICMP Echo", "ip", ip.String(), "msg", base64.StdEncoding.EncodeToString(msg))
	_, err := conn.WriteTo(msg, &net.IPAddr{IP: ip})
	return err
}

func recvICMP(conn *icmp.PacketConn, res chan<- Device) error {
	rb := make([]byte, 1500)
	n, peer, err := conn.ReadFrom(rb)
	if err != nil {
		return err
	}
	peerIP := net.ParseIP(peer.String())

	msg, err := icmp.ParseMessage(1, rb[:n])
	if err != nil {
		return err
	}
	_, ok := msg.Body.(*icmp.Echo)
	if !ok {
		switch b := msg.Body.(type) {
		case *icmp.DstUnreach:
			var dest string
			switch msg.Code {
			case 0:
				dest = "network"
			case 1:
				dest = "host"
			case 2:
				dest = "protocol"
			case 3:
				dest = "port"
			case 4:
				dest = "must-fragment"
			default:
				dest = "dest"
			}
			log.Warn("ICMP unreachable", "dest", dest)
		case *icmp.PacketTooBig:
			log.Warn("ICMP packet too big", "mtu", b.MTU)
		default:
			log.Warn("ICMP non-echo response", "response", b)
		}
		return nil
	}

	res <- Device{
		IP: peerIP,
		// TODO find MAC and Name
	}

	return nil
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
