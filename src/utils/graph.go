package utils

import (
	"errors"
)

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

func (g *Graph) InsertNode() *Node {
	for i, node := range g.nodes {
		if i != node.id {
			g.nodes[i] = Node{id: i}
			return &g.nodes[i]
		}
	}
	id := len(g.nodes)
	g.InsertNodes(Node{id: id})
	return &g.nodes[id]
}

func (g *Graph) InsertNodes(nodes ...Node) {
	for _, node := range nodes {
		node.best_nodes = []int{-1}
		if node.id >= len(g.nodes) {
			g.nodes = append(g.nodes, make([]Node, node.id-len(g.nodes)+1)...)
		}
		g.nodes[node.id] = node
	}
}

func (g *Graph) GetNode(id int) (node *Node, err error) {
	if id >= len(g.nodes) {
		return nil, errors.New("Node is not found")
	}
	return &g.nodes[id], nil
}
