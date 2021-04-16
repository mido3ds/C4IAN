package main

import (
	"time"

	"github.com/cornelk/hashmap"
)

const Age = 60

// key: 4 bytes IPv4, value: *ForwardingEntry
type FloodingTable struct {
	m hashmap.HashMap
}

type FloodingEntry struct {
	seqNumber uint32
	ageTimer *time.Timer
}

func NewFloodingTable() *FloodingTable {
	return &FloodingTable{}
}

// Get returns value associated with the given key 
// and whether the key existed or not
func (f *FloodingTable) Get(srcIP []byte) (uint32, bool) {
	v, ok := f.m.Get(IPv4ToUInt32(srcIP))
	if !ok {
		return 0, false
	}
	return v.(*FloodingEntry).seqNumber, true
}

// Set the srcIP to a new sequence number
// Restart the timer attached to that src
func (f *FloodingTable) Set(srcIP []byte, seq uint32) {
	v, ok := f.m.Get(IPv4ToUInt32(srcIP))
	// Stop the previous timer if it wasn't fired
	if ok {
		timer := v.(*FloodingEntry).ageTimer
		timer.Stop()
	} 

	// Start new Timer
	fireFunc := fireTimer(srcIP, f)
	newTimer := time.AfterFunc(Age * time.Second, fireFunc)
	entry := &FloodingEntry{ seqNumber: seq, ageTimer: newTimer }
	f.m.Set(IPv4ToUInt32(srcIP), entry)
}

// Del silently fails if key doesn't exist
func (f *FloodingTable) Del(srcIP []byte) {
	f.m.Del(IPv4ToUInt32(srcIP))
}

func (f *FloodingTable) Len() int {
	return f.m.Len()
}

func fireTimerHelper(srcIP []byte, f *FloodingTable) {
	f.Del(srcIP)
}

func fireTimer(srcIP []byte, f *FloodingTable) func() {
	return func() {
        fireTimerHelper(srcIP, f)
    }
}
