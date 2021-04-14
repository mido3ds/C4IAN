package utils

type Node struct {
	id         int
	weight     int64
	best_nodes []int
	edges      map[int]int64
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
