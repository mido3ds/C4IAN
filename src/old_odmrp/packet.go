package odmrp

import (
	"regexp"
	"strconv"
)

type IP string

// define Cast enum
type Cast int

const (
	UNICAST Cast = iota
	MULTICAST
	BROADCAST
)

type Packet struct {
	cast_mode     Cast
	source_addr   IP
	dest_addr     IP
	hops_traveled int
	time_to_live  int
	payload       string // could be []byte !?
}

func newPacket() Packet {
	var p Packet
	p.time_to_live = DEFAULT_TIME_TO_LIVE
	p.hops_traveled = 0
	return p
}

func (p *Packet) Copy() *Packet {
	new_packet := *p
	return &new_packet
}

func (ip IP) getIpType() Cast {
	is_matched, err := regexp.MatchString(MULTICAST_PATTERN, ip.ToString())
	if err != nil {
		if is_matched {
			return MULTICAST
		}
	}
	is_matched, err = regexp.MatchString(MULTICAST_PATTERN, ip.ToString())
	if err != nil {
		if is_matched {
			return BROADCAST
		}
	}
	return UNICAST
}

func (ip IP) ToString() string {
	return string(ip)
}

func (p *Packet) ToString() string {
	var cast_mode string
	switch mode := p.cast_mode; mode {
	case UNICAST:
		cast_mode = "unicast"
		break
	case MULTICAST:
		cast_mode = "multicast"
		break
	case BROADCAST:
		cast_mode = "broadcast"
		break
	}
	return ("Packet Casting Mode: " + cast_mode + "\n Source Address: " + p.source_addr.ToString() +
		", Destination Address: " + p.dest_addr.ToString() + "\nTime To Live: " + strconv.Itoa(p.time_to_live) +
		", Hops Traveled: " + strconv.Itoa(p.hops_traveled))
}
