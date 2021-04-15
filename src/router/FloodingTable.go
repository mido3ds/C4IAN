package main

import (
	"fmt"
	"time"

	"github.com/cornelk/hashmap"
)


// key: 4 bytes IPv4, value: *ForwardingEntry
type FloodingTable struct {
	m hashmap.HashMap
}

type FloodingEntry struct {
	seqNumber uint32
	ageTimer time.timer
}

func NewFloodingTable() *FloodingTable {
	return &FloodingTable{}
}

// Get returns value associated with the given key 
// and whether the key existed or not
func (f *FloodingTable) Get(srcIP []byte) (seq, bool) {
	v, ok := f.m.Get(ipToUInt32(srcIP))
	if !ok {
		return nil, false
	}
	return v.(*FloodingEntry).seqNumber, true
}

func (f *FloodingTable) Set(srcIP []byte, seq uint32) {
	if seqNumber == nil {
		panic(fmt.Errorf("you can't enter nil entry"))
	}

	// Stop the previous timer from firing
	entry, ok := f.m.Get(ipToUInt32(srcIP))
	
	// Reset Timer
	fireFunc := fireTimer(srcIP, f)
	t := time.AfterFunc(60*time.Second, fireFunc)
	entry := &FloodingEntry{ seqNumber: seq, ageTimer: t }
	f.m.Set(ipToUInt32(srcIP), entry)
}

// Del silently fails if key doesn't exist
func (f *FloodingTable) Del(srcIP []byte) {
	f.m.Del(ipToUInt32(srcIP))
}

func (f *FloodingTable) Len() int {
	return f.m.Len()
}

func fireTimer(srcIP []byte, f *FloodingTable) {
	f.m.Set(ipToUInt32(srcIP), -1)
}

