package flood

import (
	"fmt"
	"net"
	"time"

	"github.com/cornelk/hashmap"
	. "github.com/mido3ds/C4IAN/src/router/constants"
	. "github.com/mido3ds/C4IAN/src/router/ip"
)

// floodingTable is lock-free thread-safe hash table
// optimized for fastest read access
// key: 4 bytes IPv4, value: *FloodingEntry
type floodingTable struct {
	m hashmap.HashMap
}

type floodingEntry struct {
	seqNumber uint32
	ageTimer  *time.Timer
}

func newFloodingTable() *floodingTable {
	return &floodingTable{}
}

// get returns value associated with the given key
// and whether the key existed or not
func (f *floodingTable) get(srcIP net.IP) (uint32, bool) {
	v, ok := f.m.Get(IPv4ToUInt32(srcIP))
	if !ok {
		return 0, false
	}
	return v.(*floodingEntry).seqNumber, true
}

// isHighestSeqNum returns true if seq is highest sequence ever for srcIP
func (f *floodingTable) isHighestSeqNum(srcIP net.IP, seq uint32) bool {
	seq2, exist := f.get(srcIP)
	return !exist || greaterSeqNumber(seq, seq2)
}

// = MAX_UINT32/2
const halfUint32Range uint32 = 0xFFFF_FFFF >> 1

// greaterSeqNumber returns if (s1 > s2) taking considering overflow may happen
// see https://datatracker.ietf.org/doc/html/draft-gerla-manet-odmrp#section-6
func greaterSeqNumber(s1, s2 uint32) bool {
	return (s2 < s1 && (s1-s2) <= halfUint32Range) ||
		(s1 < s2 && (s2-s1) > halfUint32Range)
}

// set the srcIP to a new sequence number
// Restart the timer attached to that src
func (f *floodingTable) set(srcIP net.IP, seq uint32) {
	v, ok := f.m.Get(IPv4ToUInt32(srcIP))
	// Stop the previous timer if it wasn't fired
	if ok {
		timer := v.(*floodingEntry).ageTimer
		timer.Stop()
	}

	// TODO (low priority): use TimerQueue
	// Start new Timer
	fireFunc := fireTimer(srcIP, f)
	newTimer := time.AfterFunc(FloodingTableEntryAge*time.Second, fireFunc)
	entry := &floodingEntry{seqNumber: seq, ageTimer: newTimer}
	f.m.Set(IPv4ToUInt32(srcIP), entry)
}

// del silently fails if key doesn't exist
func (f *floodingTable) del(srcIP net.IP) {
	f.m.Del(IPv4ToUInt32(srcIP))
}

func (f *floodingTable) len() int {
	return f.m.Len()
}

func (f *floodingTable) String() string {
	s := "&FloodingTable{"
	for item := range f.m.Iter() {
		s += fmt.Sprintf(" (ip=%#v,seq=%#v)", UInt32ToIPv4(item.Key.(uint32)).String(), item.Value.(*floodingEntry).seqNumber)
	}
	s += " }"
	return s
}

func fireTimerHelper(srcIP net.IP, f *floodingTable) {
	f.del(srcIP)
}

func fireTimer(srcIP net.IP, f *floodingTable) func() {
	return func() {
		fireTimerHelper(srcIP, f)
	}
}
