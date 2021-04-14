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

func (l *DoubleLinkedList) InsertNode(n *Node, prev_element *ListElement) *ListElement {
	return l.InsertElement(&ListElement{node: n}, prev_element)
}

func (l *DoubleLinkedList) InsertElement(element_to_insert, prev_element *ListElement) *ListElement {
	tmp := prev_element.next
	element_to_insert.list = l
	element_to_insert.prev = prev_element
	prev_element.next = element_to_insert
	element_to_insert.next = tmp
	tmp.prev = element_to_insert
	l.length += 1
	return element_to_insert
}

func (l *DoubleLinkedList) RemoveElement(e *ListElement) *ListElement {
	e.prev.next = e.next
	e.next.prev = e.prev
	e.prev = nil
	e.next = nil
	e.list = nil
	l.length -= 1
	return e
}

func (l *DoubleLinkedList) PushFront(n *Node) *ListElement {
	if l.root.next == nil {
		l.Clear()
	}
	return l.InsertNode(n, &l.root)
}

func (l *DoubleLinkedList) PushBack(n *Node) *ListElement {
	if l.root.next == nil {
		l.Clear()
	}
	return l.InsertNode(n, l.Back())
}

func (l *DoubleLinkedList) PopFront() *Node {
	tmp := l.Front()
	if tmp.list == l {
		l.RemoveElement(tmp)
	}
	return tmp.node
}

func (l *DoubleLinkedList) PopBack() *Node {
	tmp := l.Back()
	if tmp.list == l {
		l.RemoveElement(tmp)
	}
	return tmp.node
}

func (l *DoubleLinkedList) PushOrdered(n *Node) *ListElement {
	if l.root.next == nil {
		l.Clear()
	}
	if l.length == 0 {
		return l.PushFront(n)
	}
	back := l.Back()
	if back.node.weight < n.weight {
		return l.InsertNode(n, l.root.prev)
	}
	curr := l.Front()
	for curr.node.weight < n.weight && curr.node.id != n.id {
		curr = curr.next
	}
	if curr.node.id == n.id {
		return curr
	}
	return l.InsertNode(n, curr.prev)
}

func (l *DoubleLinkedList) PopOrdered() *Node {
	return l.PopBack()
}
