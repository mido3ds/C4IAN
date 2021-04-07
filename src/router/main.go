package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/AkihiroSuda/go-netfilter-queue"
	"github.com/akamensky/argparse"
	"github.com/jsimonetti/rtnetlink/rtnl"
	"github.com/mdlayher/ethernet"
	"github.com/mdlayher/raw"
)

func main() {
	var ctx Context
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)

	err := parseArgs(&ctx)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	ctx.msec = NewMSecLayer(ctx.passphrase)

	addIptablesRule(ctx.queueNumber)

	registerInterruptHandler()

	// start modules
	go startForwarder(&ctx)
	startController()
}

type Context struct {
	ifaceName   string
	passphrase  string
	queueNumber int

	msec *MSecLayer
}

func parseArgs(ctx *Context) error {
	parser := argparse.NewParser("router", "Sets forwarding table in linux to route packets in adhoc-network")
	ifaceName := parser.String("i", "iface", &argparse.Options{Required: true, Help: "Interface name"})
	passphrase := parser.String("p", "pass", &argparse.Options{Required: true, Help: "Passphrase for MSec (en/de)cryption"})
	queueNumber := parser.Int("q", "queue-num", &argparse.Options{Required: false, Help: "Packets queue number", Default: 0})

	err := parser.Parse(os.Args)
	if err != nil {
		return errors.New(parser.Usage(err))
	}

	ctx.ifaceName = *ifaceName
	ctx.passphrase = *passphrase
	ctx.queueNumber = *queueNumber

	return nil
}

func registerInterruptHandler() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go interruptSignalHandler(c)
}

func interruptSignalHandler(c chan os.Signal) {
	<-c

	removeIptablesRule()

	log.Println("closing gracefully")
	os.Exit(0)
}

func addIptablesRule(queueNumber int) {
	exec.Command("iptables", "-t", "filter", "-D", "OUTPUT", "-j", "NFQUEUE").Run()
	cmd := exec.Command("iptables", "-t", "filter", "-A", "OUTPUT", "-j", "NFQUEUE", "--queue-num", strconv.Itoa(queueNumber))
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal("couldn't add iptables rule, err: ", err, ",stderr: ", string(stdoutStderr))
	}
	log.Println("added NFQUEUE rule to OUTPUT chain in iptables")
}

func removeIptablesRule() {
	cmd := exec.Command("iptables", "-t", "filter", "-D", "OUTPUT", "-j", "NFQUEUE")
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal("couldn't remove iptables rule, err: ", err, ",stderr: ", string(stdoutStderr))
	}
	log.Println("remove NFQUEUE rule to OUTPUT chain in iptables")
}

func startController() {
	// TODO
	select {}
}

func startForwarder(ctx *Context) {
	// get interfaces
	iface, err := net.InterfaceByName(ctx.ifaceName)
	if err != nil {
		log.Fatal("couldn't get interface ", ctx.ifaceName, " error: ", err)
	}

	// connect to rtnl
	rtnlConn, err := rtnl.Dial(nil)
	if err != nil {
		log.Fatal("can't establish netlink connection: ", err)
	}
	defer rtnlConn.Close()

	// open interface
	err = rtnlConn.LinkUp(iface)
	if err != nil {
		log.Fatal("can't link-up the interface", err)
	}
	log.Print(ctx.ifaceName, " is up")

	// connect raw to iface
	ifaceRawConn, err := raw.ListenPacket(iface, mSecEtherType, nil)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer ifaceRawConn.Close()

	log.Println("initalized forwarder")

	go receiveFromMACLayer(ifaceRawConn, iface.MTU, ctx.msec)
	receiveFromIPLayer(rtnlConn, iface, ifaceRawConn, ctx.msec)
}

// Make use of an unassigned EtherType to differentiate between MSec traffic and other traffic
// https://www.iana.org/assignments/ieee-802-numbers/ieee-802-numbers.xhtml
const mSecEtherType = 0x7031

// receiveFromMACLayer continuously receives messages from the interface,
// then either repeats it over loopback (if this is destination), or forwards it for another node.
// The messages may be up to the interface's MTU in size.
func receiveFromMACLayer(ifaceRawConn net.PacketConn, mtu int, msec *MSecLayer) {
	var f ethernet.Frame
	b := make([]byte, mtu)

	log.Println("started receiving from MAC layer")
	for {
		err := readRaw(ifaceRawConn, &f, b)
		if err != nil {
			log.Fatal("Couldn't read from interface, err: ", err)
		}

		// TODO: determine if to receive it or forward it
		packet, err := msec.DecryptAll(f.Payload)
		if err != nil {
			log.Fatal("Couldn't decrypt received message, err: ", err)
		}

		// i'm destination, receive message by injecting it in loopback
		err = loopbackRaw(packet)
		if err != nil {
			log.Fatal("Couldn't write to lo interface: ", err)
		}
	}
}

// receiveFromIPLayer periodically forwards packets from IP to MAC
// after encrypting them and determining their destination
func receiveFromIPLayer(rtnlConn *rtnl.Conn, iface *net.Interface, ifaceRawConn *raw.Conn, msec *MSecLayer) {
	// get packets from kernel
	nfq, err := netfilter.NewNFQueue(0, 200, netfilter.NF_DEFAULT_PACKET_SIZE)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer nfq.Close()

	packets := nfq.GetPackets()

	log.Println("started receiving from IP layer")
	for {
		select {
		case p := <-packets:
			packet := p.Packet.Data()

			// if looped back or local, accept it
			if isPacketLoopedBack(packet) || isPacketLocal(packet) {
				p.SetVerdict(netfilter.NF_ACCEPT)
				continue
			}

			// steal packet
			p.SetVerdict(netfilter.NF_DROP)

			// determine forwarding dest
			destIP, err := extractIP(packet)
			if err != nil {
				log.Println("Couldn't extract ip from packet, ignore it")
				continue
			}
			destHWAddr, err := getForwardDest(&destIP)
			if err != nil {
				log.Fatal("Couldn't determine packets destination: ", err)
			}

			// encrypt
			encryptedPacket, err := msec.Encrypt(packet)
			if err != nil {
				log.Fatal("Couldn't encrypt packet, err: ", err)
			}

			// hand it directly to the interface
			err = writeRaw(ifaceRawConn, iface.HardwareAddr, destHWAddr, encryptedPacket)
			if err != nil {
				log.Fatal("Couldn't write to the interface: ", err)
			}
		}
	}
}

func extractIP(packet []byte) (net.IP, error) {
	// TODO: check of version and errors
	// TODO: support ipv6
	return net.IPv4(packet[16], packet[17], packet[18], packet[19]), nil
}

func getForwardDest(destIP *net.IP) (net.HardwareAddr, error) {
	// TODO: lookup forwarding table for given ipaddr
	return ethernet.Broadcast, nil
}

// writeRaw sends a message over an interface, sourced from specified hardware address.
func writeRaw(ifaceRawConn net.PacketConn, source net.HardwareAddr, dest net.HardwareAddr, msg []byte) error {
	// Message is broadcast to all machines in same network segment.
	f := &ethernet.Frame{
		Destination: dest,
		Source:      source,
		EtherType:   mSecEtherType,
		Payload:     msg,
	}

	b, err := f.MarshalBinary()
	if err != nil {
		log.Fatalf("failed to marshal ethernet frame: %v", err)
	}

	// Required by Linux, even though the Ethernet frame has a destination.
	// Unused by BSD.
	addr := &raw.Addr{
		HardwareAddr: dest,
	}

	_, err = ifaceRawConn.WriteTo(b, addr)
	return err
}

// readRaw reads a message from an interface, returns the sender HWAddr
func readRaw(ifaceRawConn net.PacketConn, f *ethernet.Frame, b []byte) error {
	// read from interface
	n, _, err := ifaceRawConn.ReadFrom(b)
	if err != nil {
		return err
	}

	// unpack Ethernet II frame
	if err = f.UnmarshalBinary(b[:n]); err != nil {
		return err
	}

	return nil
}

var localip = [4]byte{127, 0, 0, 1}

var loopbackRawAddr = syscall.SockaddrInet4{
	Port: 0,
	Addr: localip,
}

func loopbackRaw(packet []byte) error {
	markAslLoopedBack(packet)
	fd, _ := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_RAW)
	return syscall.Sendto(fd, packet, 0, &loopbackRawAddr)
}

func markAslLoopedBack(packet []byte) {
	// TODO: support ipv6
	packet[1] |= byte(1)
}

func isPacketLoopedBack(packet []byte) bool {
	// TODO: support ipv6
	return (packet[1] & byte(1)) == 1
}

func isPacketLocal(packet []byte) bool {
	// TODO: support ipv6
	return packet[16] == localip[0]
}
