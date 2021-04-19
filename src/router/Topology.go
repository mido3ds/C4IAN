package main

import (
	"net"

	"github.com/starwander/goraph"
	"github.com/cornelk/hashmap"
)

type Topology struct {
	g *goraph.Graph
}

func NewTopology() *Topology {
	g := goraph.NewGraph()
	return &Topology{g: g}
}

func (t *Topology) Update(srcIP net.IP, srcNeighbors *NeighborsTable) {
	// if srcIP node already exist, it'll be ignored silently
	t.g.AddVertex(IPv4ToUInt32(srcIP), nil)
	
	for n := range srcNeighbors.m.Iter() {
		t.g.AddVertex(n.Key.(uint32), nil)

		t.g.AddEdge(IPv4ToUInt32(srcIP),
					n.Key.(uint32),
					float64(n.Value.(*NeighborEntry).cost), 
					nil)
	}
}

func (t *Topology) Dijkstra(myIP net.IP) *hashmap.HashMap {
	nextHopTable := &hashmap.HashMap{}
	
	_, prev, _ := t.g.Dijkstra(IPv4ToUInt32(myIP))

	for dst, prevNode := range prev {
		// dst is the same as the src node
		if dst == IPv4ToUInt32(myIP) {
			continue
		}

		// dst is one of the src neighbors
		if prevNode == IPv4ToUInt32(myIP) {
			nextHopTable.Set(dst.(uint32), dst.(uint32))
			continue
		}
		
		nextHop := prevNode
		// iterate till reaching one of the src neighbors
		// or of the nodes that we have already known its nextHop
		for prevNode != IPv4ToUInt32(myIP) {
			prevNodeNextHop, exist := nextHopTable.Get(prevNode.(uint32))
			if exist {
				nextHop = prevNodeNextHop 
				break
			}
			nextHop = prevNode
			prevNode = prev[prevNode]
		}

		nextHopTable.Set(dst.(uint32), nextHop.(uint32))
	}

	return nextHopTable
}