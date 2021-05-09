package tables

import (
	"encoding/json"
	"fmt"
	"log"
	"net"

	"github.com/cornelk/hashmap"
	. "github.com/mido3ds/C4IAN/src/router/ip"
)

type GroupMembersEntry struct {
	dests []net.IP
}

type GroupMembersTable struct {
	m *hashmap.HashMap
}

// json string to GroupMembersTable
func NewGroupMembersTable(j string) *GroupMembersTable {
	grpTable := &GroupMembersTable{
		m: &hashmap.HashMap{},
	}
	var reading map[string][]string
	err := json.Unmarshal([]byte(j), &reading)
	if err != nil {
		log.Panic("Wrong Group Member Table Json Format")
	}
	dests := []net.IP{}
	for key, value := range reading {
		grpIP := net.ParseIP(key)
		if grpIP != nil {
			for _, dest := range value {
				destIP := net.ParseIP(dest)
				if destIP != nil {
					dests = append(dests, destIP)
				} else {
					log.Panicf("Wrong Unicast IP Address: %#v", dest)
				}
			}
			grpTable.Set(grpIP, dests)
		} else {
			log.Panicf("Wrong Multicast IP Address: %#v", key)
		}
		dests = []net.IP{} // clear dests array
	}
	return grpTable
}

// Set the grpIP to a new destinations group
func (f *GroupMembersTable) Set(grpIP net.IP, dests []net.IP) {
	if !grpIP.IsMulticast() {
		log.Panic("Wrong Group IP Is Not Multicast IP")
	}
	for _, dest := range dests {
		if !dest.IsGlobalUnicast() {
			log.Panic("Wrong Group IP Is Not Global Unicast IP")
		}
	}
	entry := &GroupMembersEntry{dests: dests}
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
