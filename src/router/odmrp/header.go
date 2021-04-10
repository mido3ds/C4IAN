package odmrp

type HopPair struct {
	src_addr_      int
	prev_hop_addr_ int
}

// Packet Types
type JoinQuery struct {
	valid_     bool // Query header is valid ?
	seqno_     int  // unique query identifier
	hop_count_ int
	piggyback_ bool
}

type JoinReply struct {
	valid_ bool // Reply header is valid ?
	seqno_ int  // unique identifier
	count_ int  // Number of (group, next_hop) pairs
	pairs_ []HopPair
}

// Header
type Header struct {
	valid_ack_            bool // join reply ack
	multicast_group_addr_ int
	prev_hop_addr_        int
	valid_                bool // is this header in the packet and intinitialized ?
	join_query_           JoinQuery
	join_reply_           JoinReply
}

func newHeader() Header {
	var h Header
	h.valid_ = true
	h.multicast_group_addr_ = 0
	h.prev_hop_addr_ = -1

	h.join_query_.valid_ = false
	h.join_query_.seqno_ = 0
	h.join_query_.hop_count_ = 0
	h.join_query_.piggyback_ = false

	h.join_reply_.valid_ = false
	h.join_reply_.seqno_ = 0
	h.join_reply_.count_ = 0
	h.valid_ack_ = false

	return h
}

func (h *Header) Size() int {
	if h.valid_ {
		size := 0
		if h.join_query_.valid_ {
			size += JOIN_QUERY_SIZE
		}
		if h.join_reply_.valid_ {
			size += (JOIN_REPLY_SIZE + h.join_reply_.count_*8)
		}
		if h.valid_ack_ {
			size += 5
		}
		return size
	} else {
		return -1
	}
}
