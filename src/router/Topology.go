package main

import (
	"net"

	"github.com/starwander/goraph"
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

func (t *Topology) CalculateSinkTree(myIP net.IP) map[goraph.ID]goraph.ID {
	_, parents, _ := t.g.Dijkstra(IPv4ToUInt32(myIP))
	return parents
}
