package network

import (
	"encoding/binary"
	"errors"
	"github.com/charmbracelet/log"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/sync/errgroup"
	"net"
	"net/netip"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var outboundIP netip.Addr

func OutboundIP() (netip.Addr, error) {
	err := sync.OnceValue(func() error {
		conn, err := net.Dial("udp", "8.8.8.8:80")
		if err != nil {
			return err
		}
		outboundIP, _ = netip.AddrFromSlice(conn.LocalAddr().(*net.UDPAddr).IP)
		return nil
	})()
	return outboundIP, err
}

type Device struct {
	IP   netip.Addr
	Name string
	MAC  net.HardwareAddr
}

type options struct {
	timeout time.Duration
	mask    net.IPMask
}

var defaultOptions = options{
	timeout: time.Second * 1,
	mask:    net.IPv4Mask(255, 255, 255, 0),
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

func WithMask(mask net.IPMask) Option {
	return func(o *options) error {
		o.mask = mask
		return nil
	}
}

var cachedResult []Device

func StartAutoDiscovery(interval time.Duration, opts ...Option) {
	go func() {
		timer := time.NewTicker(interval)
		defer timer.Stop()

		var err error
		for range timer.C {
			cachedResult = nil

			cachedResult, err = Scan(opts...)
			if err != nil {
				log.Error("scanning error", "err", err.Error())
			} else {
				log.Debug("auto device scan complete")
			}
		}
	}()
	log.Info("auto device discovery started")
}

func Devices() ([]Device, error) {
	if cachedResult != nil {
		return cachedResult, nil
	}
	return Scan()
}

func Scan(opts ...Option) (devices []Device, err error) {
	opt := defaultOptions
	for _, o := range opts {
		if err = o(&opt); err != nil {
			return nil, err
		}
	}

	hostIP, err := OutboundIP()
	if err != nil {
		return nil, err
	}

	log.Debug("starting scan")

	conn, err := icmp.ListenPacket("udp4", "0.0.0.0")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	var wg sync.WaitGroup
	var sent atomic.Int32
	for _, ip := range localIPs(hostIP, opt.mask) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := sendICMP(conn, ip); err != nil {
				log.Warn("ICMP send error", "ip", ip, "err", err)
				return
			}
			log.Debug("sent packet", "ip", ip)
			sent.Add(1)
		}()
	}
	wg.Wait()

	n := sent.Load()
	log.Debugf("sent %d echo packets", n)

	done := make(chan error)
	res := newResults()
	var eg errgroup.Group

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

func localIPs(hostIP netip.Addr, mask net.IPMask) []netip.Addr {
	host := binary.BigEndian.Uint32(hostIP.AsSlice())
	netmask := binary.BigEndian.Uint32(mask)
	networkAddr := host & netmask
	broadcastAddr := networkAddr | ^netmask

	var ips []netip.Addr
	for i := networkAddr + 1; i <= broadcastAddr; i++ {
		var ip [4]byte
		binary.BigEndian.PutUint32(ip[:], i)
		ips = append(ips, netip.AddrFrom4(ip))
	}
	return ips
}

func sendICMP(conn *icmp.PacketConn, ip netip.Addr) error {
	msg, _ := (&icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   os.Getpid() & 0xffff,
			Seq:  1,
			Data: nil,
		},
	}).Marshal(nil)

	_, err := conn.WriteTo(msg, &net.UDPAddr{IP: ip.AsSlice()})
	return err
}

func recvICMP(conn *icmp.PacketConn, res chan<- Device) error {
	rb := make([]byte, 1500)
	n, peer, err := conn.ReadFrom(rb)
	if err != nil {
		return err
	}
	peerIP, err := netip.ParseAddr(strings.TrimSuffix(peer.String(), ":0"))
	if err != nil {
		return err
	}

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
	devices map[netip.Addr]Device
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
		devices: make(map[netip.Addr]Device),
		c:       c,
	}
	go func() {
		for device := range c {
			r.mu.Lock()
			if _, ok := r.devices[device.IP]; ok {
				log.Debug("ICMP double receive", "ip", device.IP)
			} else {
				r.devices[device.IP] = device
			}
			r.mu.Unlock()
		}
	}()
	return r
}
