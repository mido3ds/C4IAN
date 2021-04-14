package utils

import (
	"errors"
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

func (g *Graph) BestPath(src, dest int) BestPath {
	var path []int
	for curr := g.nodes[dest]; curr.id != src; curr = g.nodes[curr.best_nodes[0]] {
		path = append(path, curr.id)
	}
	path = append(path, src)
	// reverse path
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return BestPath{g.nodes[dest].weight, path}
}

func (g *Graph) ShortestPath(src, dest int) (BestPath, error) {
	// Intialize graph
	g.Initialize(src)
	old_curr := -1
	var curr *Node
	for g.visiting.length > 0 {
		curr = g.visiting.PopOrdered()
		if old_curr == curr.id {
			continue
		}
		old_curr = curr.id
		if curr.weight >= g.best {
			continue
		}
		for n, dist := range curr.edges {
			if curr.weight+dist < g.nodes[n].weight {
				g.nodes[n].weight = curr.weight + dist
				g.nodes[n].best_nodes[0] = curr.id
				// if destination update best
				if n == dest {
					g.best = curr.weight + dist
					g.visited_dest = true
					continue
				}
				g.visiting.PushOrdered(&g.nodes[n])
			}
		}
	}
	if g.visited_dest == false {
		return BestPath{}, errors.New("No path found")
	}
	return g.BestPath(src, dest), nil
}

func (g *Graph) ShortestPaths(src, dest int) (BestPaths, error) {
	//Setup graph
	g.Initialize(src)
	old_curr := -1
	var curr *Node
	for g.visiting.length > 0 {
		curr = g.visiting.PopOrdered()
		if old_curr == curr.id {
			continue
		}
		old_curr = curr.id
		if curr.weight > g.best {
			continue
		}
		for n, dist := range curr.edges {
			if (curr.weight+dist < g.nodes[n].weight) ||
				(curr.weight+dist == g.nodes[n].weight && g.nodes[n].BestContains(curr.id) == false) {
				if curr.weight+dist == g.nodes[n].weight {
					g.nodes[n].best_nodes = append(g.nodes[n].best_nodes, curr.id)
				} else {
					g.nodes[n].weight = curr.weight + dist
					g.nodes[n].best_nodes = []int{curr.id}
				}
				if n == dest {
					g.visited_dest = true
					g.best = curr.weight + dist
					continue
				}
				g.visiting.PushOrdered(&g.nodes[n])
			}
		}
	}
	if g.visited_dest == false {
		return BestPaths{}, errors.New("No path has been found")
	}
	return g.BestPaths(src, dest), nil
}

func (g *Graph) VisitPath(src, dest, curr_node int) [][]int {
	if curr_node == src {
		return [][]int{{curr_node}}
	}
	paths := [][]int{}
	for _, Node := range g.nodes[curr_node].best_nodes {
		shortest_paths := g.VisitPath(src, dest, Node)
		for i := range shortest_paths {
			paths = append(paths, append([]int{curr_node}, shortest_paths[i]...))
		}
	}
	return paths
}

func (g *Graph) BestPaths(src, dest int) BestPaths {
	paths := g.VisitPath(src, dest, dest)
	best := BestPaths{}
	for indexPaths := range paths {
		for i, j := 0, len(paths[indexPaths])-1; i < j; i, j = i+1, j-1 {
			paths[indexPaths][i], paths[indexPaths][j] = paths[indexPaths][j], paths[indexPaths][i]
		}
		best = append(best, BestPath{g.nodes[dest].weight, paths[indexPaths]})
	}

	return best
}
