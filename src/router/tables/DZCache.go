package tables

import (
	"fmt"
	"net"
	"time"

	"github.com/cornelk/hashmap"
	. "github.com/mido3ds/C4IAN/src/router/constants"
	. "github.com/mido3ds/C4IAN/src/router/ip"
	. "github.com/mido3ds/C4IAN/src/router/zhls/zid"
)

// FloodingTable is lock-free thread-safe hash table
// optimized for fastest read access
// key: 4 bytes IPv4, value: *FloodingEntry
type DZCache struct {
	m hashmap.HashMap
}

type DZEntry struct {
	zoneID   ZoneID
	ageTimer *time.Timer
}

func NewDZCache() *DZCache {
	return &DZCache{}
}

// Get returns value associated with the given key
// and whether the key existed or not
func (z *DZCache) Get(dstIP net.IP) (ZoneID, bool) {
	v, ok := z.m.Get(IPv4ToUInt32(dstIP))
	if !ok {
		return 0, false
	}
	return v.(*DZEntry).zoneID, true
}

// Set the srcIP to a new sequence number
// Restart the timer attached to that src
func (z *DZCache) Set(dstIP net.IP, zoneID ZoneID) {
	v, ok := z.m.Get(IPv4ToUInt32(dstIP))
	// Stop the previous timer if it wasn't fired
	if ok {
		timer := v.(*DZEntry).ageTimer
		timer.Stop()
	}

	// Start new Timer
	fireFunc := zoneCacheFireTimer(dstIP, z)
	newTimer := time.AfterFunc(DZCacheAge, fireFunc)

	z.m.Set(IPv4ToUInt32(dstIP), &DZEntry{zoneID: zoneID, ageTimer: newTimer})
}

// Del silently fails if key doesn't exist
func (z *DZCache) Del(dstIP net.IP) {
	z.m.Del(IPv4ToUInt32(dstIP))
}

func (z *DZCache) Len() int {
	return z.m.Len()
}

func (z *DZCache) String() string {
	s := "&DZCache{"
	for item := range z.m.Iter() {
		s += fmt.Sprintf(" (ip=%#v,zoneId=%#v)", UInt32ToIPv4(item.Key.(uint32)).String(), item.Value.(*DZEntry).zoneID)
	}
	s += " }"
	return s
}

func zoneCacheFireTimerHelper(dstIP net.IP, z *DZCache) {
	z.Del(dstIP)
}

func zoneCacheFireTimer(srcIP net.IP, z *DZCache) func() {
	return func() {
		zoneCacheFireTimerHelper(srcIP, z)
	}
}
