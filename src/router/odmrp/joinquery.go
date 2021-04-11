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
