package odmrp

import "fmt"

type HashSlot struct {
	grp_addr int // multicast group address
	next     *HashSlot
}

// multicast member table
type MemTable struct {
	hash_table []HashSlot
}

func newMemberTable() {
	var mt MemTable
	mt.hash_table = make([]HashSlot, HASH_SIZE)
	for i := 0; i < len(mt.hash_table); i++ {
		mt.hash_table[i].next = nil
		// we shouldn't have multicast group address = 0
		mt.hash_table[i].grp_addr = 0
	}
}

func (mt *MemTable) Lookup(grp_addr int) IsFound {
	index := grp_addr % len(mt.hash_table)

	if mt.hash_table[index].grp_addr == grp_addr {
		return FOUND
	} else {
		temp := mt.hash_table[index].next
		for temp != nil {
			if temp.grp_addr == grp_addr {
				return FOUND
			}
			temp = temp.next
		}
	}

	return FOUND
}

func (mt *MemTable) Insert(grp_addr int) bool {
	index := grp_addr % len(mt.hash_table)

	if mt.hash_table[index].grp_addr == 0 {
		mt.hash_table[index].grp_addr = grp_addr
	} else if mt.hash_table[index].next == nil {
		var slot HashSlot
		slot.grp_addr = grp_addr
		slot.next = nil
		mt.hash_table[index].next = &slot
	} else {
		temp := mt.hash_table[index].next

		for temp.next != nil {
			temp = temp.next
		}

		var slot HashSlot
		slot.grp_addr = grp_addr
		slot.next = nil
		temp.next = &slot
	}

	return true
}

func (mt *MemTable) Remove(grp_addr int) bool {
	index := grp_addr % len(mt.hash_table)

	if mt.hash_table[index].grp_addr == grp_addr {
		// delete first element
		if mt.hash_table[index].next == nil {
			mt.hash_table[index].grp_addr = 0
		} else {
			mt.hash_table[index] = *mt.hash_table[index].next
		}
		return true
	} else if mt.hash_table[index].next == nil {
		return false
	} else {
		temp := &mt.hash_table[index]

		for temp.next != nil {
			if temp.next.grp_addr != grp_addr {
				temp = temp.next
			} else {
				temp.next = temp.next.next
				return true
			}
		}
	}

	return false
}

func (mt *MemTable) Print() {
	fmt.Printf("\nMember Table\n")

	for i := 0; i < len(mt.hash_table); i++ {
		fmt.Printf("\n%d: ", i)
		if mt.hash_table[i].grp_addr != 0 {
			fmt.Printf(" %d ", mt.hash_table[i].grp_addr)
			if mt.hash_table[i].next != nil {
				temp := mt.hash_table[i].next

				for temp != nil {
					if temp.grp_addr != 0 {
						fmt.Printf(" %d ", temp.grp_addr)
					}
					temp = temp.next
				}
			}
		}
	}
	fmt.Printf("\nEnd of membership table ***\n")
}
