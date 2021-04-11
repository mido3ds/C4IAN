package odmrp

type JoinQueryFound int

const (
	FOUND JoinQueryFound = iota
	FOUND_LONGER
	NOT_FOUND
)

type JoinQuerySlot struct {
	src       int
	prev_hop  int
	seqno     int
	hop_count int
	next      *JoinQuerySlot
}

type JoinQueryTable struct {
	hash_table []JoinQuerySlot
}

func newJoinQueryTable() JoinQueryTable {
	var jqt JoinQueryTable
	jqt.hash_table = make([]JoinQuerySlot, JQ_SRC_HASH_SIZE)
	for i := 0; i < len(jqt.hash_table); i++ {
		var new_jq_slot JoinQuerySlot
		new_jq_slot.next = nil
		new_jq_slot.src = -1
		new_jq_slot.seqno = 0
		jqt.hash_table[i] = new_jq_slot
	}
	return jqt
}

func (jqt *JoinQueryTable) Lookup(src_addr int, seqno int, hop_count int) JoinQueryFound {
	index := src_addr % len(jqt.hash_table)

	if jqt.hash_table[index].src == src_addr &&
		jqt.hash_table[index].seqno == seqno || jqt.hash_table[index].seqno > seqno {
		if jqt.hash_table[index].hop_count > hop_count {
			return FOUND_LONGER //TODO should we update to the new hop count and prev hop?
		}
		return FOUND
	} else {
		temp := jqt.hash_table[index].next
		for temp != nil {
			if temp.src == src_addr && (temp.seqno == seqno || temp.seqno > seqno) {
				if jqt.hash_table[index].hop_count > hop_count {
					return FOUND_LONGER
				}
				return FOUND
			}
			temp = temp.next
		}
	}
	return NOT_FOUND
}

func (jqt *JoinQueryTable) FindPrevHop(src_addr int) int {
	index := src_addr % len(jqt.hash_table)

	if jqt.hash_table[index].src == src_addr {
		return jqt.hash_table[index].prev_hop
	} else {
		temp := jqt.hash_table[index].next
		for temp != nil {
			if temp.src == src_addr {
				return temp.prev_hop
			}
			temp = temp.next
		}
	}

	return 0 // not found
}
