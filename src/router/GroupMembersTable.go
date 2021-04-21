package main

import (
	"fmt"
	"net"

	"github.com/cornelk/hashmap"
)

type GroupMembersEntry struct {
	grpIP net.IP
	dests []net.IP
}

type GroupMembersTable struct {
	m hashmap.HashMap
}

// Set the grpIP to a new destinations group
func (f *GroupMembersTable) Set(grpIP net.IP, dests []net.IP) {
	entry := &GroupMembersEntry{grpIP: grpIP, dests: dests}
	f.m.Set(IPv4ToUInt32(grpIP), entry)
}

// Get returns value associated with the given key
// and whether the key existed or not
func (g *GroupMembersTable) Get(grpIP net.IP) ([]net.IP, bool) {
	v, ok := g.m.Get(IPv4ToUInt32(grpIP))
	if !ok {
		return nil, false
	}
	return v.(*GroupMembersEntry).dests, true
}

func (f *GroupMembersTable) Del(grpIP net.IP) {
	f.m.Del(IPv4ToUInt32(grpIP))
}

func (f *GroupMembersTable) Len() int {
	return f.m.Len()
}

func (f *GroupMembersTable) String() string {
	s := "&GroupMembersTable{"
	for item := range f.m.Iter() {
		s += fmt.Sprintf(" (ip=%#v, dests=(", UInt32ToIPv4(item.Key.(uint32)).String())
		for _, dest := range item.Value.(*GroupMembersEntry).dests {
			s += fmt.Sprintf("%#v, ", dest.String())
		}
		s += "))"
	}
	s += " }"
	return s
}
