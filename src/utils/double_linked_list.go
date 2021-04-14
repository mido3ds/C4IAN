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

func newDoubleLinkedList() *DoubleLinkedList {
	var l DoubleLinkedList
	return l.Clear()
}

func (l *DoubleLinkedList) Clear() *DoubleLinkedList {
	l.root.prev = &l.root
	l.root.next = &l.root
	l.length = 0
	return l
}

func (l *DoubleLinkedList) Front() *ListElement {
	if l.length > 0 {
		return l.root.next
	}
	return nil
}

func (l *DoubleLinkedList) Back() *ListElement {
	if l.length > 0 {
		return l.root.prev
	}
	return nil
}
