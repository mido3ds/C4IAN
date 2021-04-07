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

type JoinQuery struct {
	valid_     bool
	seqno_     int
	hop_count_ int
	piggyback_ bool
}

type JoinReply struct {
	valid_ bool
	seqno_ int
	count_ int
	// struct prev_hop_pairs pairs_[OD_MAX_NUM_PREV_HOP_PAIRS];
}

// string to odmrpaddr_t
type Header struct {
	valid_ack_       int
	mcastgroup_addr_ string
	prev_hop_addr_   string
	join_query_      JoinQuery
	join_reply_      JoinReply
}

func newPacket() Packet {
	var p Packet
	p.time_to_live = DEFAULT_TIME_TO_LIVE
	p.hops_traveled = 0
	return p
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
