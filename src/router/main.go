package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/AkihiroSuda/go-netfilter-queue"
	"github.com/akamensky/argparse"
	"github.com/jsimonetti/rtnetlink/rtnl"

	"github.com/mdlayher/ethernet"
	"github.com/mdlayher/raw"
)

func main() {
	parser := argparse.NewParser("router", "Sets forwarding table in linux to route packets in adhoc-network")
	ifaceName := parser.String("i", "iface", &argparse.Options{Required: true, Help: "Interface name"})
	loIfaceName := parser.String("", "lo", &argparse.Options{Required: false, Help: "Loopback interface name", Default: "lo"})
	queueNumber := parser.Int("q", "queue-num", &argparse.Options{Required: false, Help: "Packets queue number", Default: 0})

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)

	iface, err := net.InterfaceByName(*ifaceName)
	if err != nil {
		log.Fatal("couldn't get interface ", *ifaceName, " error:", err)
	}

	loIface, err := net.InterfaceByName(*loIfaceName)
	if err != nil {
		log.Fatal("couldn't get interface ", *loIfaceName, " error:", err)
	}

	log.Println("interface:", *ifaceName)
	log.Println("queue number:", *queueNumber)

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
	log.Print(*ifaceName, " is up")

	// open lo interface
	err = rtnlConn.LinkUp(loIface)
	if err != nil {
		log.Fatal("can't link-up lo interface", err)
	}
	log.Print(*loIfaceName, " is up")

	// connect raw to iface
	ifaceRawConn, err := raw.ListenPacket(iface, msecEtherType, nil)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer ifaceRawConn.Close()

	go readAllRaw(ifaceRawConn, iface.MTU, loIface)
	writeAllRaw(rtnlConn, iface, ifaceRawConn)
}

// Make use of an unassigned EtherType to differentiate between MSec traffic and other traffic
// https://www.iana.org/assignments/ieee-802-numbers/ieee-802-numbers.xhtml
const msecEtherType = 0x7031

// readAllRaw continuously receives messages over a connection, then repeats it over loopback. The messages
// may be up to the interface's MTU in size.
func readAllRaw(ifaceRawConn net.PacketConn, mtu int, loIface *net.Interface) {
	var f ethernet.Frame
	b := make([]byte, mtu)

	loConn, err := raw.ListenPacket(loIface, msecEtherType, nil)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer loConn.Close()

	for {
		n, addr, err := ifaceRawConn.ReadFrom(b)
		if err != nil {
			log.Fatalf("failed to receive message: %v", err)
		}

		// Unpack Ethernet II frame into Go representation.
		if err := (&f).UnmarshalBinary(b[:n]); err != nil {
			log.Fatalf("failed to unmarshal ethernet frame: %v", err)
		}

		log.Printf("[%s] %s", addr.String(), string(f.Payload))

		writeRaw(loConn, loIface.HardwareAddr, ethernet.Broadcast, decrypt(f.Payload))
	}
}

func writeAllRaw(rtnlConn *rtnl.Conn, iface *net.Interface, ifaceRawConn *raw.Conn) {
	// get packets from kernel
	nfq, err := netfilter.NewNFQueue(0, 200, netfilter.NF_DEFAULT_PACKET_SIZE)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer nfq.Close()
	packets := nfq.GetPackets()

	for {
		select {
		case p := <-packets:
			fmt.Println(p.Packet)
			p.SetVerdict(netfilter.NF_ACCEPT)

			writeRaw(ifaceRawConn, iface.HardwareAddr, ethernet.Broadcast, encrypt(p.Packet.Data()))
		}
	}
}

// writeRaw continuously sends a message over a connection at regular intervals,
// sourced from specified hardware address.
func writeRaw(ifaceRawConn net.PacketConn, source net.HardwareAddr, dest net.HardwareAddr, msg []byte) {
	// Message is broadcast to all machines in same network segment.
	f := &ethernet.Frame{
		Destination: dest,
		Source:      source,
		EtherType:   msecEtherType,
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

	if _, err := ifaceRawConn.WriteTo(b, addr); err != nil {
		log.Fatalf("failed to send message: %v", err)
	}
}

func encrypt(msg []byte) []byte {
	msg[0]++
	return msg
}

func decrypt(msg []byte) []byte {
	msg[0]--
	return msg
}
