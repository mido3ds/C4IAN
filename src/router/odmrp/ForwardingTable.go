package odmrp

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/cornelk/hashmap"
	. "github.com/mido3ds/C4IAN/src/router/ip"
)

const FWRD_ENTRY_TIMEOUT = 960 * time.Millisecond

// ForwardingTable is lock-free thread-safe hash table
// for multicast forwarding
// key: 4 bytes dest IPv4, value: *forwardingEntry
// routing table from destination to destIP
/*
------------------------------------
destIP | nextHop | Cost
------------------------------------
*/
type ForwardingTable struct {
	m *hashmap.HashMap
}

type forwardingEntry struct {
	nextHop  net.HardwareAddr
	cost     uint16
	ageTimer *time.Timer
}

func newRoutingTable() *ForwardingTable {
	return &ForwardingTable{
		m: &hashmap.HashMap{},
	}
}

// Get returns value associated with the given key, and whether the key existed or not
func (r *ForwardingTable) Get(destIP net.IP) (*forwardingEntry, bool) {
	v, ok := r.m.Get(IPv4ToUInt32(destIP))
	if !ok {
		return nil, false
	}

	return v.(*forwardingEntry), true
}

// Set the destIP to a new sequence number
// Restart the timer attached to that dest
// return true if inserted/refreshed successfully
func (r *ForwardingTable) Set(destIP net.IP, entry *forwardingEntry) bool {
	if entry == nil {
		log.Panic("you can't enter nil entry")
	}
	v, ok := r.m.Get(IPv4ToUInt32(destIP))
	// Stop the previous timer if it wasn't fired
	if ok {
		// if less cost refresh
		if v.(*forwardingEntry).cost >= entry.cost {
			entry.cost = v.(*forwardingEntry).cost
			entry.nextHop = v.(*forwardingEntry).nextHop
			timer := v.(*forwardingEntry).ageTimer
			timer.Stop()
		} else {
			return false
		}
	}

	// Start new Timer
	fireFunc := fireRoutingTableTimer(destIP, r)
	entry.ageTimer = time.AfterFunc(FWRD_ENTRY_TIMEOUT, fireFunc)
	r.m.Set(IPv4ToUInt32(destIP), entry)
	return true
}

// Del silently fails if key doesn't exist
func (r *ForwardingTable) Del(destIP net.IP) {
	r.m.Del(IPv4ToUInt32(destIP))
}

func (r *ForwardingTable) Len() int {
	return r.m.Len()
}

// Clear ForwardingTable
func (r *ForwardingTable) Clear() {
	r.m = &hashmap.HashMap{}
}

func (r *ForwardingTable) String() string {
	s := "&ForwardingTable{"
	for item := range r.m.Iter() {
		v := item.Value.(*forwardingEntry)
		s += fmt.Sprintf(" (destIP=%#v, nexthop=%#v, cost=%d)", UInt32ToIPv4(item.Key.(uint32)).String(), v.nextHop.String(), v.cost)
	}
	s += " }"

	return s
}

func fireRoutingTableTimerHelper(destIP net.IP, r *ForwardingTable) {
	// r.Del(destIP)
}

func fireRoutingTableTimer(destIP net.IP, r *ForwardingTable) func() {
	return func() {
		fireRoutingTableTimerHelper(destIP, r)
	}
}
