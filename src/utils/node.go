package utils

import "sync"

type Node struct {
	id         int
	weight     int64
	best_nodes []int
	edges      map[int]int64
	lock       sync.RWMutex
}

func newNode(id int) *Node {
	var new Node
	new.id = id
	new.weight = 0
	new.best_nodes = []int{-1}
	new.edges = make(map[int]int64)
	return &new
}

func (n *Node) BestContains(id int) bool {
	for _, node := range n.best_nodes {
		if node == id {
			return true
		}
	}
	return false
}

func (n *Node) AddEdge(dest int, weight int64) {
	n.lock.Lock()
	defer n.lock.Unlock()

	if n.edges == nil {
		n.edges = map[int]int64{}
	}
	n.edges[dest] = weight
}

func (n *Node) GetEdge(Destination int) (weight int64, ok bool) {
	if n.edges == nil {
		return 0, false
	}
	weight, ok = n.edges[Destination]
	return weight, ok
}
