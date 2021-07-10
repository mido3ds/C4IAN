package tables

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	. "github.com/mido3ds/C4IAN/src/router/zhls/zid"
	"github.com/starwander/goraph"
)

const topologyVertexAge = 3 * time.Second

type Topology struct {
	g    *goraph.Graph
	lock sync.RWMutex
	myIP net.IP
}

func NewTopology(myIP net.IP) *Topology {
	g := goraph.NewGraph()
	t := &Topology{g: g, myIP: myIP}
	// Start new Timer
	fireFunc := topologyFireTimer(ToNodeID(t.myIP), t)
	newTimer := time.AfterFunc(topologyVertexAge, fireFunc)
	t.g.AddVertexWithEdges(&myVertex{id: ToNodeID(myIP), outTo: make(map[NodeID]float64), inFrom: make(map[NodeID]float64), ageTimer: newTimer})
	return t
}

type myVertex struct {
	id       NodeID
	outTo    map[NodeID]float64
	inFrom   map[NodeID]float64
	ageTimer *time.Timer
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

func (t *Topology) Clear() {
	treeVertices := t.CalculateSinkTree(ToNodeID(t.myIP))

	t.lock.Lock()
	defer t.lock.Unlock()
	for vertexID, _ := range treeVertices {
		vertex, _ := t.g.GetVertex(vertexID.(NodeID))
		vertex.(*myVertex).ageTimer.Stop()
		t.g.DeleteVertex(vertexID.(NodeID))
	}

	fireFunc := topologyFireTimer(ToNodeID(t.myIP), t)
	newTimer := time.AfterFunc(topologyVertexAge, fireFunc)
	t.g.AddVertexWithEdges(&myVertex{id: ToNodeID(t.myIP), outTo: make(map[NodeID]float64), inFrom: make(map[NodeID]float64), ageTimer: newTimer})
}

func (t *Topology) Update(srcID NodeID, srcNeighbors *NeighborsTable) error {
	t.lock.Lock()
	defer t.lock.Unlock()

	// Ignore interzone LSR packets originating at my zone
	if zoneID, isZone := srcID.ToZoneID(); isZone && zoneID == MyZone().ID {
		return nil
	}

	vertex, notExist := t.g.GetVertex(srcID)
	if notExist == nil {
		t.removeOldInFromEdges(vertex.(*myVertex))
	}

	outToEdges := make(map[NodeID]float64)

	for n := range srcNeighbors.m.Iter() {
		nodeID := NodeID(n.Key.(uint64))

		// Ignore edges pointing to my zone
		if zoneID, isZone := nodeID.ToZoneID(); isZone && zoneID == MyZone().ID {
			continue
		}

		// Check if this neighbor vertex exists
		neighborVertex, notExist := t.g.GetVertex(nodeID)

		if notExist == nil {
			neighborVertex.(*myVertex).inFrom[srcID] = float64(n.Value.(*NeighborEntry).Cost)
			// Remove the old neighbor vertex
			t.g.DeleteVertex(nodeID)
			// Add the neighbor vertex with new inFrom edge
			t.g.AddVertexWithEdges(neighborVertex.(*myVertex))
		} else {
			neighborInFromEdges := make(map[NodeID]float64)
			neighborInFromEdges[srcID] = float64(n.Value.(*NeighborEntry).Cost)
			// Start new Timer
			fireFunc := topologyFireTimer(nodeID, t)
			newTimer := time.AfterFunc(topologyVertexAge, fireFunc)
			//log.Println("New msg from: ",srcID , "That added ", nodeID)
			t.g.AddVertexWithEdges(&myVertex{id: nodeID, outTo: make(map[NodeID]float64), inFrom: neighborInFromEdges, ageTimer: newTimer})
		}

		outToEdges[nodeID] = float64(n.Value.(*NeighborEntry).Cost)
	}

	vertex, notExist = t.g.GetVertex(srcID)
	if notExist == nil {
		vertex.(*myVertex).outTo = outToEdges
		vertex.(*myVertex).inFrom = t.validateInFromEdges(vertex.(*myVertex))
		vertex.(*myVertex).ageTimer.Stop()
		fireFunc := topologyFireTimer(srcID, t)
		newTimer := time.AfterFunc(topologyVertexAge, fireFunc)
		vertex.(*myVertex).ageTimer = newTimer
		//log.Println("New msg from: ", srcID)
		// Remove the old src vertex
		t.g.DeleteVertex(srcID)
		// Add the src vertex with new outTo edges
		return t.g.AddVertexWithEdges(vertex.(*myVertex))
	} else {
		// Start new Timer
		//log.Println(srcID, "is added to the topology")
		fireFunc := topologyFireTimer(srcID, t)
		newTimer := time.AfterFunc(topologyVertexAge, fireFunc)
		return t.g.AddVertexWithEdges(&myVertex{id: srcID, outTo: outToEdges, inFrom: make(map[NodeID]float64), ageTimer: newTimer})
	}
}

func (t *Topology) CalculateSinkTree(nodeID NodeID) map[goraph.ID]goraph.ID {
	t.lock.RLock()
	defer t.lock.RUnlock()
	_, parents, _ := t.g.Dijkstra(nodeID)
	return parents
}

func (t *Topology) GetGateways(srcNodeID NodeID) (gateways []NodeID) {
	t.lock.RLock()
	defer t.lock.RUnlock()
	if srcNodeID.isZone() {
		log.Panic("Trying to find gateways of zone vertex")
	}

	visited := make(map[NodeID]bool)

	var stack [](*myVertex)

	srcVertex, notExist := t.g.GetVertex(srcNodeID)
	if notExist != nil {
		log.Panic("Trying to find gateways of inexistent src vertex")
	}
	stack = append(stack, srcVertex.(*myVertex))

	for len(stack) != 0 {
		n := len(stack) - 1 // Top element
		currentVertex := stack[n]
		stack = stack[:n] // Pop

		if !visited[currentVertex.id] {
			visited[currentVertex.id] = true
		}

		markedAsGateway := false
		for vertexID := range currentVertex.outTo {
			vertex, notExist := t.g.GetVertex(vertexID)
			if notExist != nil {
				log.Panic("Incorrect topology")
			}
			if vertex.(*myVertex).id.isZone() {
				if !markedAsGateway {
					gateways = append(gateways, currentVertex.id)
					markedAsGateway = true
				}
			} else {
				if !visited[vertex.(*myVertex).id] {
					stack = append(stack, vertex.(*myVertex))
				}
			}
		}
	}
	return
}

func (t *Topology) validateInFromEdges(vertex *myVertex) map[NodeID]float64 {
	newInFrom := make(map[NodeID]float64)
	for from, cost := range vertex.inFrom {
		fromVertex, notExist := t.g.GetVertex(from)
		if notExist == nil {
			if _, ok := fromVertex.(*myVertex).outTo[vertex.id]; ok {
				newInFrom[from] = cost
			}
		}
	}
	return newInFrom
}

func (t *Topology) removeOldInFromEdges(vertex *myVertex) {
	for to, _ := range vertex.outTo {
		toVertex, toVertexNonExist := t.g.GetVertex(to)
		if toVertexNonExist == nil {
			delete(toVertex.(*myVertex).inFrom, vertex.id)
			t.g.DeleteVertex(toVertex.(*myVertex).id)
			t.g.AddVertexWithEdges(toVertex.(*myVertex))
		}
	}
}

func (t *Topology) GetNeighborZones(srcNodeID NodeID) (neighborZones []NodeID, isMaxIP bool) {
	t.lock.RLock()
	defer t.lock.RUnlock()
	if srcNodeID.isZone() {
		log.Panic("Trying to find neighbour zones of zone vertex")
	}

	isMaxIP = true
	visited := make(map[NodeID]bool)

	var stack [](*myVertex)

	srcVertex, notExist := t.g.GetVertex(srcNodeID)
	if notExist != nil {
		log.Panic("Trying to find neighbour zones  of inexistent src vertex")
	}
	stack = append(stack, srcVertex.(*myVertex))

	markedAsNeighbourZone := make(map[NodeID]bool)
	for len(stack) != 0 {
		n := len(stack) - 1 // Top element
		currentVertex := stack[n]
		stack = stack[:n] // Pop

		if !visited[currentVertex.id] {
			visited[currentVertex.id] = true
		}

		for vertexID := range currentVertex.outTo {
			vertex, notExist := t.g.GetVertex(vertexID)
			if notExist != nil {
				//log.Println(currentVertex.id, vertexID, currentVertex.outTo)
				//log.Println("Incorrect topology")
				continue
			}
			if vertex.(*myVertex).id.isZone() {
				_, marked := markedAsNeighbourZone[vertex.(*myVertex).id]
				if !marked {
					neighborZones = append(neighborZones, vertex.(*myVertex).id)
					markedAsNeighbourZone[vertex.(*myVertex).id] = true
				}
			} else {
				if vertexID > srcNodeID {
					isMaxIP = false
				}
				if !visited[vertex.(*myVertex).id] {
					stack = append(stack, vertex.(*myVertex))
				}
			}
		}
	}
	return
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

		for from, cost := range vertex.(*myVertex).inFrom {
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

func topologyFireTimerHelper(nodeID NodeID, t *Topology) {
	t.lock.Lock()
	defer t.lock.Unlock()

	vertex, notExist := t.g.GetVertex(nodeID)
	if notExist == nil {
		log.Println(nodeID, "is deleted from the topology")
		t.removeOldInFromEdges(vertex.(*myVertex))
		t.g.DeleteVertex(nodeID)
	}
}

func topologyFireTimer(nodeID NodeID, t *Topology) func() {
	return func() {
		topologyFireTimerHelper(nodeID, t)
	}
}

func (t *Topology) DisplaySinkTreeParents(sinkTreeParents map[goraph.ID]goraph.ID) {
	log.Println("----------- Sink Tree -------------")
	for dst, parent := range sinkTreeParents {
		if dst == nil {
			log.Println("Dst: ", dst, "Parent: ", parent.(NodeID))
			continue
		}
		if parent == nil {
			log.Println("Dst: ", dst.(NodeID), "Parent: ", parent)
			t.DisplayVertex(dst.(NodeID))
			continue
		}
		log.Println("Dst: ", dst.(NodeID), "Parent: ", parent.(NodeID))
		t.DisplayVertex(dst.(NodeID))
	}
	log.Println("-----------------------------------")
}
