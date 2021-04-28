package odmrp

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/cornelk/hashmap"
	. "github.com/mido3ds/C4IAN/src/router/ip"
)

type cacheEntry struct {
	// SrcIP    net.IP // commented as it is the key no need to store it again
	SeqNo    uint64
	GrpIP    net.IP
	PrevHop  net.HardwareAddr
	ageTimer *time.Timer
}

type CacheTable struct {
	m *hashmap.HashMap
}

func newCacheTable() *CacheTable {
	return &CacheTable{
		m: &hashmap.HashMap{},
	}
}

// Get returns value associated with the given key, and whether the key existed or not
func (f *CacheTable) Get(src net.IP) (*cacheEntry, bool) {
	v, ok := f.m.Get(IPv4ToUInt32(src))
	if !ok {
		return nil, false
	}

	return v.(*cacheEntry), true
}

func (f *CacheTable) Set(src net.IP, entry *cacheEntry) {
	if entry == nil {
		log.Panic("you can't enter nil entry")
	}
	f.m.Set(IPv4ToUInt32(src), entry)
}

// Del silently fails if key doesn't exist
func (f *CacheTable) Del(src net.IP) {
	f.m.Del(IPv4ToUInt32(src))
}

func (f *CacheTable) Len() int {
	return f.m.Len()
}

// Clear CacheTable
func (f *CacheTable) Clear() {
	f.m = &hashmap.HashMap{}
}

func (f *CacheTable) String() string {
	s := "&CacheTable{"
	for item := range f.m.Iter() {
		v := item.Value.(*cacheEntry)
		s += fmt.Sprintf(" (srcIP=%#v,seq=%d, grpIP=%#v,prevhop=%#v)", UInt32ToIPv4(item.Key.(uint32)).String(), v.SeqNo, v.GrpIP, v.PrevHop)
	}
	s += " }"

	return s
}
