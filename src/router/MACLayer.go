package main

import (
	"net"

	"github.com/mdlayher/ethernet"
	"github.com/mdlayher/raw"
)

type MACLayerConn struct {
	packetConn net.PacketConn
	source     net.HardwareAddr

	// dirty optimization
	// DON'T use one Conn for multiple readers!
	f *ethernet.Frame
	b []byte
}

func NewMACLayerConn(iface *net.Interface, etherType uint16) (*MACLayerConn, error) {
	packetConn, err := raw.ListenPacket(iface, etherType, nil)
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
	}, nil
}

func (c *MACLayerConn) Write(packet []byte, dest net.HardwareAddr) error {
	f := &ethernet.Frame{
		Destination: dest,
		Source:      c.source,
		EtherType:   ZIDEtherType,
		Payload:     packet,
	}

	b, err := f.MarshalBinary()
	if err != nil {
		return err
	}

	_, err = c.packetConn.WriteTo(b, &raw.Addr{HardwareAddr: dest})
	return err
}

func (c *MACLayerConn) Read() ([]byte, error) {
	n, _, err := c.packetConn.ReadFrom(c.b)
	if err != nil {
		return nil, err
	}

	if err = c.f.UnmarshalBinary(c.b[:n]); err != nil {
		return nil, err
	}

	return c.f.Payload, nil
}

func (c *MACLayerConn) Close() {
	c.packetConn.Close()
}
