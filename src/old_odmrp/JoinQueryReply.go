package odmrp

import "fmt"

type JoinSourceSlot struct {
	src       int
	prev_hop  int
	seqno     int
	hop_count int
	next      *JoinSourceSlot
}

type JoinSourceTable struct {
	hash_table []JoinSourceSlot
}

func newJoinSourceTable() JoinSourceTable {
	var jst JoinSourceTable
	jst.hash_table = make([]JoinSourceSlot, JOIN_SRC_HASH_SIZE)
	for i := 0; i < len(jst.hash_table); i++ {
		jst.hash_table[i].src = -1
		jst.hash_table[i].seqno = 0
		jst.hash_table[i].next = nil
	}

	return jst
}

func (jst *JoinSourceTable) Lookup(src_addr int, seqno int, hop_count int) IsFound {
	index := src_addr % len(jst.hash_table)

	if jst.hash_table[index].src == src_addr &&
		jst.hash_table[index].seqno == seqno || jst.hash_table[index].seqno > seqno {
		if jst.hash_table[index].hop_count > hop_count {
			return FOUND_LONGER //TODO should we update to the new hop count and prev hop?
		}
		return FOUND
	} else {
		temp := jst.hash_table[index].next
		for temp != nil {
			if temp.src == src_addr && (temp.seqno == seqno || temp.seqno > seqno) {
				if jst.hash_table[index].hop_count > hop_count {
					return FOUND_LONGER
				}
				return FOUND
			}
			temp = temp.next
		}
	}
	return NOT_FOUND
}

func (jst *JoinSourceTable) FindPrevHop(src_addr int) int {
	index := src_addr % len(jst.hash_table)

	if jst.hash_table[index].src == src_addr {
		return jst.hash_table[index].prev_hop
	} else {
		temp := jst.hash_table[index].next
		for temp != nil {
			if temp.src == src_addr {
				return temp.prev_hop
			}
			temp = temp.next
		}
	}

	return 0 // not found
}

func (jst *JoinSourceTable) Insert(src_addr int, prev_hop int, seqno, hop_count int) bool {
	index := src_addr % len(jst.hash_table)

	if jst.hash_table[index].src == -1 || jst.hash_table[index].src == src_addr {
		jst.hash_table[index].src = src_addr
		jst.hash_table[index].prev_hop = prev_hop
		jst.hash_table[index].hop_count = hop_count
		jst.hash_table[index].seqno = seqno
	} else if jst.hash_table[index].next == nil {
		var new_jq_slot JoinSourceSlot
		new_jq_slot.src = src_addr
		new_jq_slot.prev_hop = prev_hop
		new_jq_slot.hop_count = hop_count
		new_jq_slot.seqno = seqno
		new_jq_slot.next = nil
		jst.hash_table[index].next = &new_jq_slot
	} else {
		temp := jst.hash_table[index].next

		if jst.hash_table[index].next.src == src_addr {
			jst.hash_table[index].next.seqno = seqno
			jst.hash_table[index].next.prev_hop = prev_hop
			jst.hash_table[index].next.hop_count = hop_count
			return true
		}

		for temp.next != nil {
			if temp.next.src == src_addr {
				temp.next.seqno = seqno
				temp.next.prev_hop = prev_hop
				temp.next.hop_count = hop_count
				return true
			}
			temp = temp.next
		}

		var new_jq_slot JoinSourceSlot
		new_jq_slot.src = src_addr
		new_jq_slot.seqno = seqno
		new_jq_slot.prev_hop = prev_hop
		new_jq_slot.hop_count = hop_count
		new_jq_slot.next = nil
		temp.next = &new_jq_slot
	}

	return true
}

func (jst *JoinSourceTable) Clear() {
	var t JoinSourceTable
	t = newJoinSourceTable()
	jst = &t
}

func (jst *JoinSourceTable) Remove(src_addr int) bool {
	// get index
	index := src_addr % len(jst.hash_table)

	if jst.hash_table[index].src == src_addr {
		if jst.hash_table[index].next == nil {
			jst.hash_table[index].src = -1
		} else {
			jst.hash_table[index] = *jst.hash_table[index].next
		}
		return true
	} else if jst.hash_table[index].next == nil {
		return false
	} else {
		temp := &jst.hash_table[index]

		for temp.next != nil {
			if temp.next.src != src_addr {
				temp = temp.next
			} else {
				temp.next = temp.next.next
				return true
			}
		}
	}

	return false
}

func (jst *JoinSourceTable) Print() {
	fmt.Printf("\nJoin Source Table\n")

	for i := 0; i < len(jst.hash_table); i++ {
		fmt.Printf("\n%d: ", i)

		if jst.hash_table[i].src != 0 {
			fmt.Printf("%d ", jst.hash_table[i].src)
			fmt.Printf(":%d ", jst.hash_table[i].prev_hop)
			fmt.Printf(":%d ", jst.hash_table[i].hop_count)
			fmt.Printf(":%d ", jst.hash_table[i].seqno)

			if jst.hash_table[i].next != nil {
				temp := jst.hash_table[i].next

				for temp != nil {
					if temp.src != 0 {
						fmt.Printf(":%d ", temp.src)
						fmt.Printf(":%d ", temp.prev_hop)
						fmt.Printf(":%d ", temp.hop_count)
						fmt.Printf(":%d ", temp.seqno)
					}
					temp = temp.next
				}
			}
		}
	}
	fmt.Printf("\nEnd Join Source Table\n")
}
