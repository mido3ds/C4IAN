package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/AkihiroSuda/go-netfilter-queue"
	"github.com/akamensky/argparse"
	"github.com/jsimonetti/rtnetlink/rtnl"
)

func main() {
	parser := argparse.NewParser("router", "Sets forwarding table in linux to route packets in adhoc-network")
	ifaceName := parser.String("i", "iface", &argparse.Options{Required: false, Help: "Interface name", Default: "wlan0"})
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

	log.Println("interface:", *ifaceName)
	log.Println("queue number:", *queueNumber)

	// connect to rtnl
	conn, err := rtnl.Dial(nil)
	if err != nil {
		log.Fatal("can't establish netlink connection: ", err)
	}
	defer conn.Close()

	// open interface
	err = conn.LinkUp(iface)
	if err != nil {
		log.Fatal("can't link-up the interface", err)
	}
	log.Print(*ifaceName, " is up")

	// get packets from kernel
	nfq, err := netfilter.NewNFQueue(0, 200, netfilter.NF_DEFAULT_PACKET_SIZE)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer nfq.Close()
	packets := nfq.GetPackets()

	for true {
		select {
		case p := <-packets:
			handlePacket(p, conn, iface)
		}
	}
}

func handlePacket(p netfilter.NFPacket, conn *rtnl.Conn, iface *net.Interface) {
	// TODO
	// get the destination from the packet,
	// hold it pack untill you get its routing (next hops),
	// set the next hop in the linux forwarding table see (https://pkg.go.dev/github.com/jsimonetti/rtnetlink@v0.0.0-20210122163228-8d122574c736/rtnl#Conn.RouteReplace),
	// then accept the packet

	fmt.Println(p.Packet)
	p.SetVerdict(netfilter.NF_ACCEPT) // accept all for now
}
