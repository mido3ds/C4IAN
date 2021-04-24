package tables

import (
	"fmt"
	"log"
	"net"

	. "github.com/mido3ds/C4IAN/src/router/ip"
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
		// Check if this neighbor vertex exists
		neighborVertex, notExist := t.g.GetVertex(n.Key.(uint32))

		if notExist == nil {
			neighborVertex.(*myVertex).inFrom[IPv4ToUInt32(srcIP)] = float64(n.Value.(*NeighborEntry).Cost)
			// Remove the old neighbor vertex
			t.g.DeleteVertex(n.Key.(uint32))
			// Add the neighbor vertex with new inFrom edge
			t.g.AddVertexWithEdges(neighborVertex.(*myVertex))
		} else {
			neighborInFromEdges := make(map[uint32]float64)
			neighborInFromEdges[IPv4ToUInt32(srcIP)] = float64(n.Value.(*NeighborEntry).Cost)
			t.g.AddVertexWithEdges(&myVertex{id: n.Key.(uint32), outTo: make(map[uint32]float64), inFrom: neighborInFromEdges})
		}

		outToEdges[n.Key.(uint32)] = float64(n.Value.(*NeighborEntry).Cost)
	}

	vertex, notExist := t.g.GetVertex(IPv4ToUInt32(srcIP))
	if notExist == nil {
		vertex.(*myVertex).outTo = outToEdges
		// Remove the old src vertex
		t.g.DeleteVertex(IPv4ToUInt32(srcIP))
		// Add the src vertex with new outTo edges
		return t.g.AddVertexWithEdges(vertex.(*myVertex))
	} else {
		return t.g.AddVertexWithEdges(&myVertex{id: IPv4ToUInt32(srcIP), outTo: outToEdges, inFrom: make(map[uint32]float64)})
	}
}

func (t *Topology) CalculateSinkTree(myIP net.IP) map[goraph.ID]goraph.ID {
	_, parents, _ := t.g.Dijkstra(IPv4ToUInt32(myIP))
	return parents
}

func (t *Topology) DisplayVertex(id goraph.ID) {
	vertex,_ := t.g.GetVertex(id)
	s := ""
	if vertex == nil {
		s += fmt.Sprintf("Vertex: %v isn't existd\n", id)
	} else {
		s += fmt.Sprintf("Vertex: %v, ", UInt32ToIPv4(vertex.(*myVertex).id))
		for to, cost := range vertex.(*myVertex).outTo {
			s += fmt.Sprintf("Edge to %v with cost %v, ",UInt32ToIPv4(to), cost)
		}

		for from, cost := range vertex.(*myVertex).outTo {
			s += fmt.Sprintf("Edge from %v with cost %v, ", UInt32ToIPv4(from), cost)
		}
	}

	log.Println(s)
}

func (t *Topology) DisplayEdge(fromID goraph.ID, toID goraph.ID) {
	edge,_ := t.g.GetEdge(fromID, toID)
	s := ""
	t.g.GetEdge(fromID, toID)
	if edge == nil {
		s += fmt.Sprintf("Edge from %v to %v isn't existd\n", fromID, toID)
	} else {
		s += fmt.Sprintf("Edge from %v to %v with cost %v\n", edge.(*myEdge).from, edge.(*myEdge).to, edge.(*myEdge).weight)
	}

	log.Println(s)
}
