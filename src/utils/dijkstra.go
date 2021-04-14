package utils

import (
	"math"
)

type BestPath struct {
	weight int64
	Path   []int
}

type BestPaths []BestPath

func (g *Graph) Initialize(src int) {
	g.visiting = newDoubleLinkedList()
	g.visited_dest = false
	g.FillDefaults(int64(math.MaxInt64)-2, -1)
	g.best = int64(math.MaxInt64)
	g.nodes[src].weight = 0
	g.visiting.PushOrdered(&g.nodes[src])
}
