package odmrp

import (
	"fmt"
	"net"
	"time"

	"github.com/cornelk/hashmap"
	. "github.com/mido3ds/C4IAN/src/router/ip"
)

const MTE_TIMEOUT = 960 * time.Millisecond

// MemberTable is lock-free thread-safe hash table
// for multicast forwarding
// key: 4 bytes grpIP IPv4, value: *memberEntry
type MemberTable struct {
	m *hashmap.HashMap
}

type memberEntry struct {
	ageTimer *time.Timer
}

func newMemberTable() *MemberTable {
	return &MemberTable{
		m: &hashmap.HashMap{},
	}
}

// Get returns value associated with the given key, and whether the key existed or not
func (mt *MemberTable) Get(grpIP net.IP) bool {
	_, ok := mt.m.Get(IPv4ToUInt32(grpIP))
	return ok
}

// Set the grpIP to a new sequence number
// Restart the timer attached to that src
func (mt *MemberTable) Set(grpIP net.IP) {
	v, ok := mt.m.Get(IPv4ToUInt32(grpIP))
	// Stop the previous timer if it wasn't fired
	if ok {
		timer := v.(*memberEntry).ageTimer
		timer.Stop()
	}

	// Start new Timer
	fireFunc := fireMemberTableTimer(grpIP, mt)
	ageTimer := time.AfterFunc(MTE_TIMEOUT, fireFunc)
	mt.m.Set(IPv4ToUInt32(grpIP), &memberEntry{ageTimer: ageTimer})
}

// Del silently fails if key doesn't exist
func (mt *MemberTable) Del(grpIP net.IP) {
	mt.m.Del(IPv4ToUInt32(grpIP))
}

func (mt *MemberTable) Len() int {
	return mt.m.Len()
}

// Clear MemberTable
func (mt *MemberTable) Clear() {
	mt.m = &hashmap.HashMap{}
}

func (mt *MemberTable) String() string {
	s := "&MemberTable{"
	for item := range mt.m.Iter() {
		s += fmt.Sprintf(" (grpIP=%#v)", UInt32ToIPv4(item.Key.(uint32)).String())
	}
	s += " }"

	return s
}

func fireMemberTableTimerHelper(grpIP net.IP, mt *MemberTable) {
	// mt.Del(grpIP)
}

func fireMemberTableTimer(grpIP net.IP, mt *MemberTable) func() {
	return func() {
		fireMemberTableTimerHelper(grpIP, mt)
	}
}
