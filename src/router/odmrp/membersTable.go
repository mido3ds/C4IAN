package odmrp

import (
	"fmt"
	"net"

	"github.com/cornelk/hashmap"
	. "github.com/mido3ds/C4IAN/src/router/constants"
	. "github.com/mido3ds/C4IAN/src/router/ip"
	. "github.com/mido3ds/C4IAN/src/router/tables"
)

// const MembersTableTimeout = 960 * time.Millisecond
// const MembersTableTimeout = 2 * time.Second

// membersTable is lock-free thread-safe hash table
// for multicast forwarding
// key: 4 bytes grpIP IPv4, value: *memberEntry
type membersTable struct {
	m      *hashmap.HashMap
	timers *TimersQueue
}

type memberEntry struct {
	timer *Timer
}

func newMembersTable(timers *TimersQueue) *membersTable {
	return &membersTable{
		m:      &hashmap.HashMap{},
		timers: timers,
	}
}

// get returns value associated with the given key, and whether the key existed or not
func (mt *membersTable) get(grpIP net.IP) bool {
	_, ok := mt.m.Get(IPv4ToUInt32(grpIP))
	return ok
}

// set the grpIP to a new sequence number
// Restart the timer attached to that src
func (mt *membersTable) set(grpIP net.IP) {
	v, ok := mt.m.Get(IPv4ToUInt32(grpIP))
	// Stop the previous timer if it wasn't fired
	if ok {
		v.(*memberEntry).timer.Stop()
	}
	// Start new Timer
	timer := mt.timers.Add(MembersTableTimeout, func() {
		mt.del(grpIP)
	})
	mt.m.Set(IPv4ToUInt32(grpIP), &memberEntry{timer: timer})
}

// del silently fails if key doesn't exist
func (mt *membersTable) del(grpIP net.IP) {
	mt.m.Del(IPv4ToUInt32(grpIP))
}

func (mt *membersTable) len() int {
	return mt.m.Len()
}

// clear MemberTable
func (mt *membersTable) clear() {
	mt.m = &hashmap.HashMap{}
}

func (mt *membersTable) String() string {
	s := "&MemberTable{"
	for item := range mt.m.Iter() {
		s += fmt.Sprintf(" (grpIP=%#v)", UInt32ToIPv4(item.Key.(uint32)).String())
	}
	s += " }"

	return s
}
