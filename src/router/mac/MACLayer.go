package mac

import (
	"log"
	"net"

	"github.com/mdlayher/ethernet"
	"github.com/mdlayher/raw"
)

var BroadcastMACAddr = ethernet.Broadcast

type MACLayerConn struct {
	packetConn net.PacketConn
	source     net.HardwareAddr
	etherType  EtherType

	// dirty optimization
	// DON'T use one Conn for multiple readers!
	f *ethernet.Frame
	b []byte
}

func NewMACLayerConn(iface *net.Interface, etherType EtherType) (*MACLayerConn, error) {
	packetConn, err := raw.ListenPacket(iface, uint16(etherType), nil)
	if err != nil {
		return nil, err
	}

	f := new(ethernet.Frame)
	b := make([]byte, iface.MTU)

	return &MACLayerConn{
		packetConn: packetConn,
		source:     iface.HardwareAddr,
		f:          f,
		b:          b,
		etherType:  etherType,
	}, nil
}

func (c *MACLayerConn) Write(packet []byte, dest net.HardwareAddr) {
	f := ethernet.Frame{
		Destination: dest,
		Source:      c.source,
		EtherType:   ethernet.EtherType(c.etherType),
		Payload:     packet,
	}

	b, err := f.MarshalBinary()
	if err != nil {
		log.Panic("failed to marshal ethernet frame, err: ", err)
	}

	_, err = c.packetConn.WriteTo(b, &raw.Addr{HardwareAddr: dest})
	if err != nil {
		log.Println("packet length: ", len(packet))
		log.Panic("failed to write to device driver, err: ", err)
	}
}

func (c *MACLayerConn) Read() []byte {
	n, _, err := c.packetConn.ReadFrom(c.b)
	if err != nil {
		log.Panic("failed to read from device driver, err: ", err)
	}

	err = c.f.UnmarshalBinary(c.b[:n])
	if err != nil {
		log.Panic("failed to unmarshal ethernet frame, err: ", err)
	}

	return c.f.Payload
}

func (c *MACLayerConn) Close() {
	c.packetConn.Close()
}
