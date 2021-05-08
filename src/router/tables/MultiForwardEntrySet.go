package tables

import (
	"fmt"
	"net"
	"time"

	"github.com/cornelk/hashmap"
	. "github.com/mido3ds/C4IAN/src/router/ip"
)

const FWRD_SET_ENTRY_TIMEOUT = 1 * time.Second

// MultiForwardEntrySet is lock-free thread-safe set
// for multicast forwarding
// key: 8 bytes IPv4, value: *NextHopEntry
type MultiForwardEntrySet struct {
	Items *hashmap.HashMap
}

type NextHopEntry struct {
	NextHop  net.HardwareAddr // Can be deleted if memory is so critical
	ageTimer *time.Timer
}

func NewMultiForwardEntrySet() *MultiForwardEntrySet {
	return &MultiForwardEntrySet{
		Items: &hashmap.HashMap{},
	}
}

// Get returns value associated with the given key, and whether the key existed or not
func (f *MultiForwardEntrySet) Get(nextHop net.HardwareAddr) (*NextHopEntry, bool) {
	v, ok := f.Items.Get(HwAddrToUInt64(nextHop))
	if !ok {
		return nil, false
	}

	return v.(*NextHopEntry), true
}

func (s *MultiForwardEntrySet) Set(nextHop net.HardwareAddr) {
	nextHopKey := HwAddrToUInt64(nextHop)
	v, ok := s.Items.Get(nextHopKey)
	// Stop the previous timer if it wasn't fired
	if ok {
		timer := v.(*NextHopEntry).ageTimer
		timer.Stop()
	}

	// Start new Timer
	fireFunc := fireSetTimer(nextHopKey, s)
	entry := &NextHopEntry{
		NextHop:  nextHop,
		ageTimer: time.AfterFunc(FWRD_SET_ENTRY_TIMEOUT, fireFunc),
	}
	s.Items.Set(nextHopKey, entry)
}

// Del silently fails if key doesn't exist
func (f *MultiForwardEntrySet) Del(NextHop net.HardwareAddr) {
	f.Items.Del(HwAddrToUInt64(NextHop))
}

func (f *MultiForwardEntrySet) Len() int {
	return f.Items.Len()
}

// Clear MultiForwardEntrySet
func (s *MultiForwardEntrySet) Clear() {
	s.Items = &hashmap.HashMap{}
}

func (s *MultiForwardEntrySet) String() string {
	str := "&MultiForwardEntrySet{"
	for item := range s.Items.Iter() {
		str += fmt.Sprintf("%#v, ", item.Value.(*NextHopEntry).NextHop.String())
	}
	str += " }"

	return str
}

func fireSetTimer(nextHopKey uint64, s *MultiForwardEntrySet) func() {
	return func() {
		// s.Items.Del(nextHopKey)
	}
}
