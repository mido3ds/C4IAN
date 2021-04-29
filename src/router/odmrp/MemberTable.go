package odmrp

import (
	"fmt"
	"log"
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
	grpIP    net.IP
	ageTimer *time.Timer
}

func newMemberTable() *MemberTable {
	return &MemberTable{
		m: &hashmap.HashMap{},
	}
}

// Get returns value associated with the given key, and whether the key existed or not
func (mt *MemberTable) Get(srcIP net.IP) (*memberEntry, bool) {
	v, ok := mt.m.Get(IPv4ToUInt32(srcIP))
	if !ok {
		return nil, false
	}

	return v.(*memberEntry), true
}

// Set the srcIP to a new sequence number
// Restart the timer attached to that src
func (mt *MemberTable) Set(srcIP net.IP, entry *memberEntry) {
	if entry == nil {
		log.Panic("you can't enter nil entry")
	}
	v, ok := mt.m.Get(IPv4ToUInt32(srcIP))
	// Stop the previous timer if it wasn't fired
	if ok {
		timer := v.(*memberEntry).ageTimer
		timer.Stop()
	}

	// Start new Timer
	fireFunc := fireMemberTableTimer(srcIP, mt)
	entry.ageTimer = time.AfterFunc(MTE_TIMEOUT, fireFunc)
	mt.m.Set(IPv4ToUInt32(srcIP), entry)
}

// Del silently fails if key doesn't exist
func (mt *MemberTable) Del(srcIP net.IP) {
	mt.m.Del(IPv4ToUInt32(srcIP))
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
		v := item.Value.(*memberEntry)
		s += fmt.Sprintf(" (srcIP=%#v, GrpIP=%#v)", UInt32ToIPv4(item.Key.(uint32)).String(), v.grpIP)
	}
	s += " }"

	return s
}

func fireMemberTableTimerHelper(srcIP net.IP, mt *MemberTable) {
	mt.Del(srcIP)
}

func fireMemberTableTimer(srcIP net.IP, mt *MemberTable) func() {
	return func() {
		fireMemberTableTimerHelper(srcIP, mt)
	}
}
