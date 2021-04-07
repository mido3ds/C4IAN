package main

import (
	"fmt"
	"net"
)

type IPPacket struct {
	version byte
	destIP  net.IP
}

var (
	errTooSmall         = fmt.Errorf("too small IP packet")
	errInvalidIPVersion = fmt.Errorf("invalid ip version")
)

func ParseIPPacket(buffer []byte) (*IPPacket, error) {
	var ip net.IP
	version := byte(buffer[0]) >> 4

	if version == 4 {
		if len(buffer) < 20 {
			return nil, errTooSmall
		}
		ip = net.IPv4(buffer[16], buffer[17], buffer[18], buffer[19])
	} else if version == 6 {
		if len(buffer) < 40 {
			return nil, errTooSmall
		}
		ip = buffer[24:40]
	} else {
		return nil, errInvalidIPVersion
	}

	return &IPPacket{
		version: version,
		destIP:  ip,
	}, nil
}
