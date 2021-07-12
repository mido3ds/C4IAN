package odmrp

import (
	"fmt"
	"log"
	"net"

	"github.com/cornelk/hashmap"
	. "github.com/mido3ds/C4IAN/src/router/constants"
	. "github.com/mido3ds/C4IAN/src/router/ip"
	. "github.com/mido3ds/C4IAN/src/router/tables"
)

// forwardingTable is lock-free thread-safe hash table
// for multicast forwarding
// key: 4 bytes dest IPv4, value: *forwardingEntry
// routing table from destination to destIP
/*
------------------------------------
destIP | nextHop | Cost
------------------------------------
*/
type forwardingTable struct {
	m      *hashmap.HashMap
	timers *TimersQueue
}

type forwardingEntry struct {
	nextHop net.HardwareAddr
	cost    uint16
	timer   *Timer
}

func newForwardingTable(timers *TimersQueue) *forwardingTable {
	return &forwardingTable{
		m:      &hashmap.HashMap{},
		timers: timers,
	}
}

// get returns value associated with the given key, and whether the key existed or not
func (r *forwardingTable) get(destIP net.IP) (*forwardingEntry, bool) {
	v, ok := r.m.Get(IPv4ToUInt32(destIP))
	if !ok {
		return nil, false
	}

	return v.(*forwardingEntry), true
}

// set the destIP to a new sequence number
// Restart the timer attached to that dest
// return true if inserted/refreshed successfully
func (r *forwardingTable) set(destIP net.IP, entry *forwardingEntry) bool {
	if entry == nil {
		log.Panic("you can't enter nil entry")
	}
	v, ok := r.m.Get(IPv4ToUInt32(destIP))
	// Stop the previous timer if it wasn't fired
	if ok {
		// if less cost refresh
		val := v.(*forwardingEntry)
		if val.cost < entry.cost {
			return false
		}
		val.timer.Stop()
	}

	// Start new Timer
	entry.timer = r.timers.Add(ForwardTableTimeout, func() {
		r.del(destIP)
	})
	r.m.Set(IPv4ToUInt32(destIP), entry)
	return true
}

// del silently fails if key doesn't exist
func (r *forwardingTable) del(destIP net.IP) {
	r.m.Del(IPv4ToUInt32(destIP))
}

func (r *forwardingTable) len() int {
	return r.m.Len()
}

// clear ForwardingTable
func (r *forwardingTable) clear() {
	r.m = &hashmap.HashMap{}
}

func (r *forwardingTable) String() string {
	s := "&ForwardingTable{"
	for item := range r.m.Iter() {
		v := item.Value.(*forwardingEntry)
		s += fmt.Sprintf(" (destIP=%#v, nexthop=%#v, cost=%d)", UInt32ToIPv4(item.Key.(uint32)).String(), v.nextHop.String(), v.cost)
	}
	s += " }"

	return s
}
