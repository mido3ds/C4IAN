package odmrp

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/cornelk/hashmap"
	. "github.com/mido3ds/C4IAN/src/router/ip"
)

const RTE_TIMEOUT = 960 * time.Microsecond

// RoutingTable is lock-free thread-safe hash table
// for multicast forwarding
// key: 4 bytes src IPv4, value: *routingEntry
// routing table from destination to srcIP
type RoutingTable struct {
	m *hashmap.HashMap
}

type routingEntry struct {
	nextHop  net.HardwareAddr
	ageTimer *time.Timer
}

func newRoutingTable() *RoutingTable {
	return &RoutingTable{
		m: &hashmap.HashMap{},
	}
}

// Get returns value associated with the given key, and whether the key existed or not
func (r *RoutingTable) Get(src net.IP) (*routingEntry, bool) {
	v, ok := r.m.Get(IPv4ToUInt32(src))
	if !ok {
		return nil, false
	}

	return v.(*routingEntry), true
}

// Set the srcIP to a new sequence number
// Restart the timer attached to that src
func (r *RoutingTable) Set(srcIP net.IP, entry *routingEntry) {
	if entry == nil {
		log.Panic("you can't enter nil entry")
	}
	v, ok := r.m.Get(IPv4ToUInt32(srcIP))
	// Stop the previous timer if it wasn't fired
	if ok {
		timer := v.(*routingEntry).ageTimer
		timer.Stop()
	}

	// Start new Timer
	fireFunc := fireRoutingTableTimer(srcIP, r)
	entry.ageTimer = time.AfterFunc(RTE_TIMEOUT, fireFunc)
	r.m.Set(IPv4ToUInt32(srcIP), entry)
}

// Del silently fails if key doesn't exist
func (r *RoutingTable) Del(srcIP net.IP) {
	r.m.Del(IPv4ToUInt32(srcIP))
}

func (r *RoutingTable) Len() int {
	return r.m.Len()
}

// Clear RoutingTable
func (r *RoutingTable) Clear() {
	r.m = &hashmap.HashMap{}
}

func (r *RoutingTable) String() string {
	s := "&RoutingTable{"
	for item := range r.m.Iter() {
		v := item.Value.(*routingEntry)
		s += fmt.Sprintf(" (srcIP=%#v, nexhop=%#v)", UInt32ToIPv4(item.Key.(uint32)).String(), v.nextHop)
	}
	s += " }"

	return s
}

func fireRoutingTableTimerHelper(srcIP net.IP, r *RoutingTable) {
	r.Del(srcIP)
}

func fireRoutingTableTimer(srcIP net.IP, r *RoutingTable) func() {
	return func() {
		fireRoutingTableTimerHelper(srcIP, r)
	}
}
