package odmrp

import (
	"fmt"
	"time"
)

type JoinReplyHashSlot struct {
	jr_srcs              JoinSourceTable
	grp_addr             int // multicast group address
	forward_refresh_time Time
	jr_refresh_time      Time
	next                 *JoinReplyHashSlot
}

type JoinReplyTable struct {
	hash_table []JoinReplyHashSlot
}

func newJoinReplyTable() JoinReplyTable {
	var jt JoinReplyTable
	jt.hash_table = make([]JoinReplyHashSlot, JOIN_HASH_SIZE)
	for i := 0; i < len(jt.hash_table); i++ {
		// we shouldn't have multicast group addresses of 0
		jt.hash_table[i].grp_addr = 0
		jt.hash_table[i].forward_refresh_time = 0
		jt.hash_table[i].next = nil
	}

	return jt
}

func (jt *JoinReplyTable) ForwardRefreshTime(grp_addr int) Time {
	temp := jt.FindGroup(grp_addr)

	if temp != nil {
		return temp.forward_refresh_time
	}

	return 0
}

func (jt *JoinReplyTable) JoinReplyRefreshTime(grp_addr int) Time {
	temp := jt.FindGroup(grp_addr)

	if temp != nil {
		return temp.jr_refresh_time
	}

	return 0
}

func (jt *JoinReplyTable) Lookup(grp_addr int, src_addr int, seqno int, hop_count int) IsFound {
	index := grp_addr % len(jt.hash_table)

	if jt.hash_table[index].grp_addr == grp_addr {
		return jt.hash_table[index].jr_srcs.Lookup(src_addr, seqno, hop_count)
	} else {
		temp := jt.hash_table[index].next
		for temp != nil {
			if temp.grp_addr == grp_addr {
				return temp.jr_srcs.Lookup(src_addr, seqno, hop_count)
			}
			temp = temp.next
		}
	}

	return NOT_FOUND
}

func (jt *JoinReplyTable) FindGroup(grp_addr int) *JoinReplyHashSlot {
	index := grp_addr % len(jt.hash_table)

	if jt.hash_table[index].grp_addr == grp_addr {
		return &jt.hash_table[index]
	} else {
		temp := jt.hash_table[index].next
		for temp != nil {
			if temp.grp_addr == grp_addr {
				return temp
			}
			temp = temp.next
		}
	}

	return nil
}

func (jt *JoinReplyTable) FindPrevHop(grp_addr int, src_addr int) int {
	index := grp_addr % len(jt.hash_table)
	current_time := time.Now()

	if jt.hash_table[index].grp_addr == grp_addr {
		if Time(current_time.Second())-jt.ForwardRefreshTime(grp_addr) < FLAG_TIMEOUT {
			return jt.hash_table[index].jr_srcs.FindPrevHop(src_addr)
		} else {
			return 0 // prev hop not found
		}
	} else {
		temp := jt.hash_table[index].next
		for temp != nil {
			if temp.grp_addr == grp_addr {
				if Time(current_time.Second())-jt.ForwardRefreshTime(grp_addr) < FLAG_TIMEOUT {
					return temp.jr_srcs.FindPrevHop(src_addr)
				} else {
					return 0 // prev hop not found
				}
			}
			temp = temp.next
		}
	}

	return 0 // prev hop not found
}

func (jt *JoinReplyTable) Insert(grp_addr int, src_addr int, prev_hop int, seqno int, hop_count int) bool {
	index := grp_addr % len(jt.hash_table)

	if jt.hash_table[index].grp_addr == 0 || jt.hash_table[index].grp_addr == grp_addr {
		jt.hash_table[index].grp_addr = grp_addr
		jt.hash_table[index].jr_srcs.Insert(src_addr, prev_hop, seqno, hop_count)
	} else if jt.hash_table[index].next == nil {
		var jq_slot JoinReplyHashSlot
		jq_slot.grp_addr = grp_addr
		jq_slot.jr_srcs.Insert(src_addr, prev_hop, seqno, hop_count)
		jq_slot.next = nil
		jt.hash_table[index].next = &jq_slot
	} else {
		temp := jt.hash_table[index].next

		if temp.grp_addr == grp_addr {
			temp.jr_srcs.Insert(src_addr, prev_hop, seqno, hop_count)
			return true
		}

		for temp.next != nil {
			if temp.next.grp_addr == grp_addr {
				temp.next.jr_srcs.Insert(src_addr, prev_hop, seqno, hop_count)
				return true
			}
			temp = temp.next
		}

		var jq_slot JoinReplyHashSlot
		jq_slot.grp_addr = grp_addr
		jq_slot.jr_srcs.Insert(src_addr, prev_hop, seqno, hop_count)
		jq_slot.next = nil
		temp.next = &jq_slot
	}

	return true
}

func (jt *JoinReplyTable) Remove(grp_addr int, src_addr int) bool {
	index := grp_addr % len(jt.hash_table)

	if jt.hash_table[index].grp_addr == grp_addr {
		jt.hash_table[index].jr_srcs.Remove(src_addr)
		if jt.hash_table[index].next == nil {
			jt.hash_table[index].grp_addr = 0
		} else {
			jt.hash_table[index] = *jt.hash_table[index].next
		}
		return true
	} else if jt.hash_table[index].next == nil {
		return false
	} else {
		temp := &jt.hash_table[index]

		for temp.next != nil {
			if temp.next.grp_addr != grp_addr {
				temp = temp.next
			} else {
				temp1 := temp.next
				temp.next = temp.next.next
				temp1.jr_srcs.Remove(src_addr)
				return true
			}
		}
	}

	return false
}

func (jt *JoinReplyTable) Print() {
	fmt.Printf("\nJoin Table\n")
	for i := 0; i < len(jt.hash_table); i++ {
		fmt.Printf("\n%d: ", i)
		if jt.hash_table[i].grp_addr != 0 {
			fmt.Printf(" %d ", jt.hash_table[i].grp_addr)
			jt.hash_table[i].jr_srcs.Print()
			if jt.hash_table[i].next != nil {
				temp := jt.hash_table[i].next

				for temp != nil {
					if temp.grp_addr != 0 {
						fmt.Printf(" %d ", temp.grp_addr)
						temp.jr_srcs.Print()
					}
					temp = temp.next
				}
			}
		}
	}
	fmt.Printf("\nEnd Join Table\n")
}
