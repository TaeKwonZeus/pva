package network

import (
	"errors"
	"github.com/charmbracelet/log"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
	"golang.org/x/sync/errgroup"
	"net"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"
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

type options struct {
	timeout time.Duration
}

var defaultOptions = options{
	timeout: time.Second * 3,
}

type Option func(*options) error

func WithTimeout(timeout time.Duration) Option {
	return func(o *options) error {
		if timeout.Seconds() < 1 {
			return errors.New("timeout cannot be less than 1 second")
		}
		o.timeout = timeout
		return nil
	}
}

func Scan(opts ...Option) (devices []Device, err error) {
	opt := defaultOptions
	for _, o := range opts {
		if err = o(&opt); err != nil {
			return nil, err
		}
	}

	log.Debug("starting scan")

	conn, err := icmp.ListenPacket("udp4", "0.0.0.0")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	var wg sync.WaitGroup
	var sent atomic.Int32
	for _, ip := range localIPs() {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := sendICMP(conn, ip); err != nil {
				log.Error("ICMP send error", "ip", ip, "err", err)
				return
			}
			sent.Add(1)
		}()
	}
	wg.Wait()

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
		close(res.c)
		return res.get(), err
	case <-time.After(opt.timeout):
		log.Debug("ICMP timeout")
		close(res.c)
		return res.get(), nil
	}
}

func localIPs() []net.IP {
	var ips []net.IP
	for i := byte(1); i < 255; i++ {
		ips = append(ips, net.IPv4(192, 168, 0, i) /* net.IPv4(192, 168, 1, i) */)
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

	_, err := conn.WriteTo(msg, &net.UDPAddr{IP: ip})
	return err
}

func recvICMP(conn *icmp.PacketConn, res chan<- Device) error {
	rb := make([]byte, 1500)
	n, peer, err := conn.ReadFrom(rb)
	if err != nil {
		return err
	}
	peerIP := net.ParseIP(strings.TrimSuffix(peer.String(), ":0"))

	msg, err := icmp.ParseMessage(1, rb[:n])
	if err != nil {
		return err
	}

	if msg.Type != ipv4.ICMPTypeEchoReply && msg.Type != ipv6.ICMPTypeEchoReply {
		log.Debug("non-echo response", "ip", peerIP, "msg", msg)
		return nil
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
			log.Warn("ICMP unreachable", "dest", dest, "ip", peerIP.String())
		case *icmp.PacketTooBig:
			log.Warn("ICMP packet too big", "mtu", b.MTU, "ip", peerIP.String())
		default:
			log.Warn("ICMP non-echo response", "msg", msg, "ip", peerIP.String())
		}
		return nil
	}

	log.Debug("echo response", "ip", peerIP, "msg", *msg)
	res <- Device{
		IP: peerIP,
		// TODO find MAC and Name
	}

	return nil
}

type results struct {
	devices map[string]Device
	mu      sync.Mutex
	c       chan<- Device
}

func (r *results) get() []Device {
	r.mu.Lock()
	defer r.mu.Unlock()

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

			r.mu.Lock()
			if _, ok := r.devices[ip]; ok {
				log.Debug("ICMP double receive", "ip", ip)
			} else {
				r.devices[ip] = device
			}
			r.mu.Unlock()
		}
	}()
	return r
}
