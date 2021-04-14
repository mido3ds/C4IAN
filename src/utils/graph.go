package utils

type Graph struct {
	best         int64
	visited_dest bool
	nodes        []Node
	visiting     *DoubleLinkedList
}

func newGraph() *Graph {
	return new(Graph)
}

func (g *Graph) FillDefaults(weight int64, best_node int) {
	for i := range g.nodes {
		g.nodes[i].best_nodes = []int{best_node}
		g.nodes[i].weight = weight
	}
}
