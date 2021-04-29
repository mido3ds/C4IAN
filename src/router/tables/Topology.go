package tables

import (
	"fmt"
	"log"
	"sync"

	"github.com/starwander/goraph"
)

type Topology struct {
	g    *goraph.Graph
	lock sync.RWMutex
}

func NewTopology() *Topology {
	g := goraph.NewGraph()
	return &Topology{g: g}
}

type myVertex struct {
	id     NodeID
	outTo  map[NodeID]float64
	inFrom map[NodeID]float64
}

func (vertex *myVertex) ID() goraph.ID {
	return vertex.id
}

type myEdge struct {
	from   NodeID
	to     NodeID
	weight float64
}

func (edge *myEdge) Get() (goraph.ID, goraph.ID, float64) {
	return edge.from, edge.to, edge.weight
}

func (vertex *myVertex) Edges() (edges []goraph.Edge) {
	for to, weight := range vertex.outTo {
		edges = append(edges, &myEdge{vertex.id, to, weight})
	}
	for from, weight := range vertex.inFrom {
		edges = append(edges, &myEdge{from, vertex.id, weight})
	}
	return
}

func (t *Topology) Update(srcID NodeID, srcNeighbors *NeighborsTable) error {
	t.lock.Lock()
	defer t.lock.Unlock()

	outToEdges := make(map[NodeID]float64)

	for n := range srcNeighbors.m.Iter() {
		// Check if this neighbor vertex exists
		neighborVertex, notExist := t.g.GetVertex(NodeID(n.Key.(uint64)))

		if notExist == nil {
			neighborVertex.(*myVertex).inFrom[srcID] = float64(n.Value.(*NeighborEntry).Cost)
			// Remove the old neighbor vertex
			t.g.DeleteVertex(NodeID(n.Key.(uint64)))
			// Add the neighbor vertex with new inFrom edge
			t.g.AddVertexWithEdges(neighborVertex.(*myVertex))
		} else {
			neighborInFromEdges := make(map[NodeID]float64)
			neighborInFromEdges[srcID] = float64(n.Value.(*NeighborEntry).Cost)
			t.g.AddVertexWithEdges(&myVertex{id: NodeID(n.Key.(uint64)), outTo: make(map[NodeID]float64), inFrom: neighborInFromEdges})
		}

		outToEdges[NodeID(n.Key.(uint64))] = float64(n.Value.(*NeighborEntry).Cost)
	}

	vertex, notExist := t.g.GetVertex(srcID)
	if notExist == nil {
		vertex.(*myVertex).outTo = outToEdges
		// Remove the old src vertex
		t.g.DeleteVertex(srcID)
		// Add the src vertex with new outTo edges
		return t.g.AddVertexWithEdges(vertex.(*myVertex))
	} else {
		return t.g.AddVertexWithEdges(&myVertex{id: srcID, outTo: outToEdges, inFrom: make(map[NodeID]float64)})
	}
}

func (t *Topology) CalculateSinkTree(nodeID NodeID) map[goraph.ID]goraph.ID {
	t.lock.RLock()
	defer t.lock.RUnlock()
	_, parents, _ := t.g.Dijkstra(nodeID)
	return parents
}

func (t *Topology) DisplayVertex(id goraph.ID) {
	vertex, _ := t.g.GetVertex(id)
	s := ""
	if vertex == nil {
		s += fmt.Sprintf("Vertex: %v doesn't exist\n", id)
	} else {
		s += fmt.Sprintf("Vertex: %v, ", vertex.(*myVertex).id)
		for to, cost := range vertex.(*myVertex).outTo {
			s += fmt.Sprintf("Edge to %v with cost %v, ", to, cost)
		}

		for from, cost := range vertex.(*myVertex).outTo {
			s += fmt.Sprintf("Edge from %v with cost %v, ", from, cost)
		}
	}

	log.Println(s)
}

func (t *Topology) DisplayEdge(fromID goraph.ID, toID goraph.ID) {
	edge, _ := t.g.GetEdge(fromID, toID)
	s := ""
	t.g.GetEdge(fromID, toID)
	if edge == nil {
		s += fmt.Sprintf("Edge from %v to %v isn't existd\n", fromID, toID)
	} else {
		s += fmt.Sprintf("Edge from %v to %v with cost %v\n", edge.(*myEdge).from, edge.(*myEdge).to, edge.(*myEdge).weight)
	}

	log.Println(s)
}
