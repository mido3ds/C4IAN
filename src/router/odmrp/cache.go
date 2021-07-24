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

// key: srcIP
// value: (seqNo, grpIP, prevHop, cost)
type cacheEntry struct {
	seqNo   uint64
	grpIP   net.IP
	prevHop net.HardwareAddr
	cost    int8
	timer   *Timer
}

type cache struct {
	m      *hashmap.HashMap
	timers *TimersQueue
}

func newCache(timers *TimersQueue) *cache {
	return &cache{
		m:      &hashmap.HashMap{},
		timers: timers,
	}
}

// get returns value associated with the given key, and whether the key existed or not
func (f *cache) get(src net.IP) (*cacheEntry, bool) {
	v, ok := f.m.Get(IPv4ToUInt32(src))
	if !ok {
		return nil, false
	}

	return v.(*cacheEntry), true
}

func (f *cache) set(src net.IP, entry *cacheEntry) bool {
	if entry == nil {
		log.Panic("you can't enter nil entry")
	}
	v, ok := f.m.Get(IPv4ToUInt32(src))
	if ok {
		val := v.(*cacheEntry)
		// if it doesn't has less or equal cost cost take it
		if val.seqNo > entry.seqNo || (val.seqNo == entry.seqNo && val.cost < entry.cost) {
			return false
		}
		val.timer.Stop()
	}
	entry.timer = f.timers.Add(ODMRPCacheTimeout, func() {
		f.del(src)
	})
	f.m.Set(IPv4ToUInt32(src), entry)
	return true
}

// del silently fails if key doesn't exist
func (f *cache) del(src net.IP) {
	f.m.Del(IPv4ToUInt32(src))
}

func (f *cache) len() int {
	return f.m.Len()
}

// clear CacheTable
func (f *cache) clear() {
	f.m = &hashmap.HashMap{}
}

func (f *cache) String() string {
	s := "&CacheTable{"
	for item := range f.m.Iter() {
		v := item.Value.(*cacheEntry)
		s += fmt.Sprintf(" (srcIP=%#v, seq=%d, grpIP=%#v, prevhop=%#v, Cost=%d)", UInt32ToIPv4(item.Key.(uint32)).String(), v.seqNo, v.grpIP.String(), v.prevHop.String(), v.cost)
	}
	s += " }"

	return s
}
