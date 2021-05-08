package odmrp

import (
	"fmt"
	"net"
	"time"

	"github.com/cornelk/hashmap"
	. "github.com/mido3ds/C4IAN/src/router/ip"
)

const mteTimeout = 960 * time.Millisecond

// membersTable is lock-free thread-safe hash table
// for multicast forwarding
// key: 4 bytes grpIP IPv4, value: *memberEntry
type membersTable struct {
	m *hashmap.HashMap
}

type memberEntry struct {
	ageTimer *time.Timer
}

func newMembersTable() *membersTable {
	return &membersTable{
		m: &hashmap.HashMap{},
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
		timer := v.(*memberEntry).ageTimer
		timer.Stop()
	}

	// Start new Timer
	fireFunc := fireMembersTableTimer(grpIP, mt)
	ageTimer := time.AfterFunc(mteTimeout, fireFunc)
	mt.m.Set(IPv4ToUInt32(grpIP), &memberEntry{ageTimer: ageTimer})
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

func fireMemberTableTimerHelper(grpIP net.IP, mt *membersTable) {
	// mt.Del(grpIP)
}

func fireMembersTableTimer(grpIP net.IP, mt *membersTable) func() {
	return func() {
		fireMemberTableTimerHelper(grpIP, mt)
	}
}
