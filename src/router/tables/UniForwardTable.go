package tables

import (
	"fmt"
	"log"
	"net"

	"github.com/cornelk/hashmap"
)

// UniForwardTable is lock-free thread-safe hash table
// for forwarding to single hop (not for multicast)
// key: 4 bytes IPv4, value: *UniForwardingEntry
type UniForwardTable struct {
	m *hashmap.HashMap
}

type UniForwardingEntry struct {
	NextHopMAC net.HardwareAddr
	DestZoneID uint32 // TODO: use ZoneID
}

func NewUniForwardTable() *UniForwardTable {
	return &UniForwardTable{
		m: &hashmap.HashMap{},
	}
}

// Get returns value associated with the given key, and whether the key existed or not
func (f *UniForwardTable) Get(nodeID NodeID) (*UniForwardingEntry, bool) {
	v, ok := f.m.Get(uint64(nodeID))
	if !ok {
		return nil, false
	}
	return v.(*UniForwardingEntry), true
}

func (f *UniForwardTable) Set(nodeID NodeID, entry *UniForwardingEntry) {
	if entry == nil {
		log.Panic("you can't enter nil entry")
	}
	f.m.Set(uint64(nodeID), entry)
}

// Del silently fails if key doesn't exist
func (f *UniForwardTable) Del(nodeID NodeID) {
	f.m.Del(uint64(nodeID))
}

func (f *UniForwardTable) Len() int {
	return f.m.Len()
}

func (f *UniForwardTable) Clear() {
	// Create a new hashmap as the underlying hashmap lacks a clear function :)
	f.m = &hashmap.HashMap{}
}

func (f *UniForwardTable) String() string {
	s := "&ForwardTable{"
	for item := range f.m.Iter() {
		s += fmt.Sprintf(" (ip=%#v,mac=%#v)\n", fmt.Sprint(item.Key), item.Value.(*UniForwardingEntry).NextHopMAC.String())

	}
	s += " }"
	return s
}
