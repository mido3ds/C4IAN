package main

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

type odmrp_join_query struct {
	valid_     int
	seqno_     int
	hop_count_ int
	piggyback_ int
}

type odmrp_join_reply struct {
	valid_ int
	seqno_ int
	count_ int
	// struct prev_hop_pairs pairs_[OD_MAX_NUM_PREV_HOP_PAIRS];
}

// string to odmrpaddr_t
type hdr_odmrp struct {
	valid_ack_       int
	mcastgroup_addr_ string
	prev_hop_addr_   string
	join_query_      odmrp_join_query
	join_reply_      odmrp_join_reply
}

func new_packet() Packet {
	var p Packet
	p.time_to_live = DEFAULT_TIME_TO_LIVE
	p.hops_traveled = 0
	return p
}

func (ip IP) get_address_ip_type() Cast {
	is_matched, err := regexp.MatchString(MULTICAST_PATTERN, string(ip))
	if err != nil {
		if is_matched {
			return MULTICAST
		}
	}
	is_matched, err = regexp.MatchString(MULTICAST_PATTERN, string(ip))
	if err != nil {
		if is_matched {
			return BROADCAST
		}
	}
	return UNICAST
}

func (p *Packet) to_string() string {
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
	return ("Packet Casting Mode: " + cast_mode + "\n Source Address: " + string(p.source_addr) +
		", Destination Address: " + string(p.dest_addr) + "\nTime To Live: " + strconv.Itoa(p.time_to_live) +
		", Hops Traveled: " + strconv.Itoa(p.hops_traveled))
}
