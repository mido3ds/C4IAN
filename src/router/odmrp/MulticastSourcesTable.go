package odmrp

import (
	"fmt"
	"net"
	"time"

	"github.com/cornelk/hashmap"
	. "github.com/mido3ds/C4IAN/src/router/ip"
)

const MSE_TIMEOUT = 2 * time.Second // TODO use another values

// MulticastSourcesTable is lock-free thread-safe hash table
// for multicast forwarding
// key: 4 bytes grpIP IPv4, value: *multicastSourceEntry
type MulticastSourcesTable struct {
	m *hashmap.HashMap
}

type multicastSourceEntry struct {
	ageTimer *time.Timer
}

func newMulticastSourcesTable() *MulticastSourcesTable {
	return &MulticastSourcesTable{
		m: &hashmap.HashMap{},
	}
}

// Get returns value associated with the given key, and whether the key existed or not
func (mt *MulticastSourcesTable) Get(srcIP net.IP) bool {
	_, ok := mt.m.Get(IPv4ToUInt32(srcIP))
	return ok
}

// Set the srcIP to a new sequence number
// Restart the timer attached to that src
func (mt *MulticastSourcesTable) Set(srcIP net.IP) {
	v, ok := mt.m.Get(IPv4ToUInt32(srcIP))
	// Stop the previous timer if it wasn't fired
	if ok {
		timer := v.(*multicastSourceEntry).ageTimer
		timer.Stop()
	}

	// Start new Timer
	fireFunc := fireMulticastSourcesTableTimer(srcIP, mt)
	entry := &multicastSourceEntry{ageTimer: time.AfterFunc(MSE_TIMEOUT, fireFunc)}
	mt.m.Set(IPv4ToUInt32(srcIP), entry)
}

// Del silently fails if key doesn't exist
func (mt *MulticastSourcesTable) Del(srcIP net.IP) {
	mt.m.Del(IPv4ToUInt32(srcIP))
}

func (mt *MulticastSourcesTable) Len() int {
	return mt.m.Len()
}

// Clear MulticastSourcesTable
func (mt *MulticastSourcesTable) Clear() {
	mt.m = &hashmap.HashMap{}
}

func (mt *MulticastSourcesTable) String() string {
	s := "&MulticastSourcesTable{"
	for item := range mt.m.Iter() {
		s += fmt.Sprintf(" (srcIP=%#v)", UInt32ToIPv4(item.Key.(uint32)).String())
	}
	s += " }"

	return s
}

func fireMulticastSourcesTableTimerHelper(srcIP net.IP, mt *MulticastSourcesTable) {
	// mt.Del(srcIP)
}

func fireMulticastSourcesTableTimer(srcIP net.IP, mt *MulticastSourcesTable) func() {
	return func() {
		fireMulticastSourcesTableTimerHelper(srcIP, mt)
	}
}
