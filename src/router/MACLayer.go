package main

import (
	"log"
	"net"

	"github.com/mdlayher/ethernet"
	"github.com/mdlayher/raw"
)

type MACLayerConn struct {
	packetConn net.PacketConn
	source     net.HardwareAddr
	etherType  ethernet.EtherType

	// dirty optimization
	// DON'T use one Conn for multiple readers!
	f *ethernet.Frame
	b []byte
}

func NewMACLayerConn(iface *net.Interface, etherType ethernet.EtherType) (*MACLayerConn, error) {
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
	f := &ethernet.Frame{
		Destination: dest,
		Source:      c.source,
		EtherType:   c.etherType,
		Payload:     packet,
	}

	b, err := f.MarshalBinary()
	if err != nil {
		log.Panic("failed to write to device driver, err: ", err)
	}

	_, err = c.packetConn.WriteTo(b, &raw.Addr{HardwareAddr: dest})
	if err != nil {
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
		log.Panic("failed to read from device driver, err: ", err)
	}

	return c.f.Payload
}

func (c *MACLayerConn) Close() {
	c.packetConn.Close()
}
