package utils

type ListElement struct {
	list       *DoubleLinkedList
	next, prev *ListElement
	node       *Node
}

type DoubleLinkedList struct {
	root   ListElement
	length int
}
