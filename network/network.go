package network

import (
	"encoding/binary"
	"github.com/charmbracelet/log"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/sync/errgroup"
	"maps"
	"net"
	"net/netip"
	"os"
	"slices"
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
	IP    netip.Addr
	Name  string
	MAC   net.HardwareAddr
	IsTCP bool
}

var cachedResult []Device

func instantTick(interval time.Duration) <-chan time.Time {
	c := make(chan time.Time)
	go func() {
		c <- time.Now()
		for t := range time.Tick(interval) {
			c <- t
		}
		close(c)
	}()
	return c
}

func StartAutoDiscovery(mask net.IPMask, timeout time.Duration, interval time.Duration) {
	go func() {
		var err error
		for range instantTick(interval) {
			cachedResult, err = Scan(mask, timeout)
			if err != nil {
				log.Error("scanning error", "err", err.Error())
			} else {
				log.Debug("auto device scan complete")
			}
		}
	}()
	log.Info("auto device discovery started")
}

func Devices() []Device {
	for cachedResult == nil {
	}
	return slices.Clone(cachedResult)
}

func Scan(mask net.IPMask, timeout time.Duration) (devices []Device, err error) {
	hostIP, err := OutboundIP()
	if err != nil {
		return nil, err
	}

	log.Debug("starting scan", "host", hostIP.String())

	ips := localIPs(hostIP, mask)
	res := make(map[netip.Addr]Device)

	icmpScanner, err := scanICMP(ips, timeout)
	if err != nil {
		return nil, err
	}
	tcpScanner := scanTCP(ips, timeout)
	t := time.After(timeout)

	for {
		select {
		case dev, k := <-icmpScanner:
			if !k {
				continue
			}
			if _, ok := res[dev.IP]; !ok {
				res[dev.IP] = dev
			}
		case dev, k := <-tcpScanner:
			if !k {
				continue
			}
			res[dev.IP] = dev
		case <-t:
			devices = slices.Collect(maps.Values(res))
			log.Info("scan complete", "devices", devices)
			return devices, nil
		}
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

func scanICMP(ips []netip.Addr, timeout time.Duration) (<-chan Device, error) {
	c := make(chan Device)

	conn, err := icmp.ListenPacket("udp4", "0.0.0.0")
	if err != nil {
		return nil, err
	}

	if err = conn.SetDeadline(time.Now().Add(timeout)); err != nil {
		return nil, err
	}

	go func() {
		defer conn.Close()

		var wg sync.WaitGroup
		var sent atomic.Int32
		for _, ip := range ips {
			wg.Add(1)
			go func() {
				defer wg.Done()
				if err := sendICMP(conn, ip); err != nil {
					log.Debug("ICMP send error", "ip", ip, "err", err)
					return
				}
				log.Debug("sent packet", "ip", ip)
				sent.Add(1)
			}()
		}
		wg.Wait()

		n := sent.Load()
		log.Debugf("sent %d echo packets", n)

		var eg errgroup.Group
		for range n {
			eg.Go(func() error {
				return recvICMP(conn, c)
			})
		}
		if err := eg.Wait(); err != nil {
			log.Warn("ICMP recv error", "err", err)
		}
		close(c)
	}()

	return c, nil
}

func sendICMP(conn *icmp.PacketConn, ip netip.Addr) error {
	msg, _ := (&icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:  os.Getpid() & 0xffff,
			Seq: 1,
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

func scanTCP(ips []netip.Addr, timeout time.Duration) <-chan Device {
	c := make(chan Device)

	go func() {
		var eg errgroup.Group
		for _, ip := range ips {
			eg.Go(func() error {
				err := dialTCP(ip, timeout, c)
				if err != nil {
					log.Debug("TCP dial error", "ip", ip, "err", err)
					return err
				}
				log.Debug("TCP dial complete", "ip", ip)
				return nil
			})
		}
		if err := eg.Wait(); err != nil {
			log.Warn("TCP dial error", "err", err)
		}
		close(c)
	}()

	return c
}

func dialTCP(ip netip.Addr, timeout time.Duration, res chan<- Device) error {
	conn, err := net.Dial("tcp", ip.String()+":http")
	if err != nil {
		return err
	}

	if err = conn.SetDeadline(time.Now().Add(timeout)); err != nil {
		return err
	}

	host, _, err := net.SplitHostPort(conn.RemoteAddr().String())
	if err != nil {
		return err
	}

	addr, err := netip.ParseAddr(host)
	if err != nil {
		return err
	}

	res <- Device{
		IP:    addr,
		IsTCP: true,
	}

	return conn.Close()
}
