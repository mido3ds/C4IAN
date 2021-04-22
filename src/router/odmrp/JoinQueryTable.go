package odmrp

import (
	"fmt"
	"net"
	"time"
)

// type SlotInterface interface{}
// type TableInterface interface {
// 	ForwardRefreshTime(int) Time
// 	Lookup(int, int, int, int) IsFound
// 	FindGroup(int) *SlotInterface
// 	FindPrevHop(int, int) int
// 	Insert(int, int, int, int, int) bool
// 	Remove(int, int) bool
// 	Print()
// }

type JoinQueryHashSlot struct {
	grp_addr             int // multicast group address
	jq_srcs              JoinSourceTable
	forward_refresh_time Time
	next                 *JoinQueryHashSlot
}

type JoinQueryTable struct {
	hash_table []JoinQueryHashSlot
}

func newJoinQueryTable() JoinQueryTable {
	var jt JoinQueryTable
	jt.hash_table = make([]JoinQueryHashSlot, JOIN_HASH_SIZE)
	for i := 0; i < len(jt.hash_table); i++ {
		// we shouldn't have multicast group addresses of 0
		jt.hash_table[i].grp_addr = 0
		jt.hash_table[i].forward_refresh_time = 0
		jt.hash_table[i].next = nil
	}

	return jt
}

func (jt *JoinQueryTable) ForwardRefreshTime(grp_addr int) Time {
	temp := jt.FindGroup(grp_addr)

	if temp != nil {
		return temp.forward_refresh_time
	}

	return 0
}

func (jt *JoinQueryTable) Lookup(grp_addr int, src_addr int, seqno int, hop_count int) IsFound {
	index := grp_addr % len(jt.hash_table)

	if jt.hash_table[index].grp_addr == grp_addr {
		return jt.hash_table[index].jq_srcs.Lookup(src_addr, seqno, hop_count)
	} else {
		temp := jt.hash_table[index].next
		for temp != nil {
			if temp.grp_addr == grp_addr {
				return temp.jq_srcs.Lookup(src_addr, seqno, hop_count)
			}
			temp = temp.next
		}
	}

	return NOT_FOUND
}

func (jt *JoinQueryTable) FindGroup(grp_addr int) *JoinQueryHashSlot {
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

func (jt *JoinQueryTable) FindPrevHop(grp_addr int, src_addr int) int {
	index := grp_addr % len(jt.hash_table)
	current_time := time.Now()

	if jt.hash_table[index].grp_addr == grp_addr {
		if Time(current_time.Second())-jt.ForwardRefreshTime(grp_addr) < FLAG_TIMEOUT {
			return jt.hash_table[index].jq_srcs.FindPrevHop(src_addr)
		} else {
			return 0 // prev hop not found
		}
	} else {
		temp := jt.hash_table[index].next
		for temp != nil {
			if temp.grp_addr == grp_addr {
				if Time(current_time.Second())-jt.ForwardRefreshTime(grp_addr) < FLAG_TIMEOUT {
					return temp.jq_srcs.FindPrevHop(src_addr)
				} else {
					return 0 // prev hop not found
				}
			}
			temp = temp.next
		}
	}

	return 0 // prev hop not found
}

func (jt *JoinQueryTable) Insert(grp_addr int, src_addr int, prev_hop int, seqno int, hop_count int) bool {
	index := grp_addr % len(jt.hash_table)

	if jt.hash_table[index].grp_addr == 0 || jt.hash_table[index].grp_addr == grp_addr {
		jt.hash_table[index].grp_addr = grp_addr
		jt.hash_table[index].jq_srcs.Insert(src_addr, prev_hop, seqno, hop_count)
	} else if jt.hash_table[index].next == nil {
		var jq_slot JoinQueryHashSlot
		jq_slot.grp_addr = grp_addr
		jq_slot.jq_srcs.Insert(src_addr, prev_hop, seqno, hop_count)
		jq_slot.next = nil
		jt.hash_table[index].next = &jq_slot
	} else {
		temp := jt.hash_table[index].next

		if temp.grp_addr == grp_addr {
			temp.jq_srcs.Insert(src_addr, prev_hop, seqno, hop_count)
			return true
		}

		for temp.next != nil {
			if temp.next.grp_addr == grp_addr {
				temp.next.jq_srcs.Insert(src_addr, prev_hop, seqno, hop_count)
				return true
			}
			temp = temp.next
		}

		var jq_slot JoinQueryHashSlot
		jq_slot.grp_addr = grp_addr
		jq_slot.jq_srcs.Insert(src_addr, prev_hop, seqno, hop_count)
		jq_slot.next = nil
		temp.next = &jq_slot
	}

	return true
}

func (jt *JoinQueryTable) Remove(grp_addr int, src_addr int) bool {
	index := grp_addr % len(jt.hash_table)

	if jt.hash_table[index].grp_addr == grp_addr {
		jt.hash_table[index].jq_srcs.Remove(src_addr)
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
				temp1.jq_srcs.Remove(src_addr)
				return true
			}
		}
	}

	return false
}

func (jt *JoinQueryTable) Print() {
	fmt.Printf("\nJoin Table\n")
	for i := 0; i < len(jt.hash_table); i++ {
		fmt.Printf("\n%d: ", i)
		if jt.hash_table[i].grp_addr != 0 {
			fmt.Printf(" %d ", jt.hash_table[i].grp_addr)
			jt.hash_table[i].jq_srcs.Print()
			if jt.hash_table[i].next != nil {
				temp := jt.hash_table[i].next

				for temp != nil {
					if temp.grp_addr != 0 {
						fmt.Printf(" %d ", temp.grp_addr)
						temp.jq_srcs.Print()
					}
					temp = temp.next
				}
			}
		}
	}
	fmt.Printf("\nEnd Join Table\n")
}
