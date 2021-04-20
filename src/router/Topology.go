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

type myVertex struct {
	id     uint32
	outTo  map[uint32]float64
	inFrom map[uint32]float64
}

func (vertex *myVertex) ID() goraph.ID {
	return vertex.id
}

type myEdge struct {
	from   uint32
	to     uint32
	weight float64
}

func (edge *myEdge) Get() (goraph.ID, goraph.ID, float64) {
	return edge.from, edge.to, edge.weight
}

func (vertex *myVertex) Edges() (edges []goraph.Edge) {
	edges = make([]goraph.Edge, len(vertex.outTo)+len(vertex.inFrom))
	i := 0
	for to, weight := range vertex.outTo {
		edges[i] = &myEdge{vertex.id, to, weight}
		i++
	}
	for from, weight := range vertex.inFrom {
		edges[i] = &myEdge{from, vertex.id, weight}
		i++
	}
	return
}

func (t *Topology) Update(srcIP net.IP, srcNeighbors *NeighborsTable) error {
	outToEdges := make(map[uint32]float64)

	for n := range srcNeighbors.m.Iter() {
		outToEdges[n.Key.(uint32)] = float64(n.Value.(*NeighborEntry).cost)
	}

	// remove the src vertex
	t.g.DeleteVertex(IPv4ToUInt32(srcIP))

	// add the src vertex with new edges
	return t.g.AddVertexWithEdges(&myVertex{id: IPv4ToUInt32(srcIP), outTo:  outToEdges})
}

func (t *Topology) CalculateSinkTree(myIP net.IP) map[goraph.ID]goraph.ID {
	_, parents, _ := t.g.Dijkstra(IPv4ToUInt32(myIP))
	return parents
}
