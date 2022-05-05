package linkedlist

func NewLinkedList[V comparable]() *LinkedList[V] {
	var v V
	l := &LinkedList[V]{
		head: NewNode(v, nil, nil),
		tail: NewNode(v, nil, nil),
	}
	l.head.next = l.tail
	l.tail.prev = l.head
	return l
}

type LinkedList[V comparable] struct {
	head *Node[V]
	tail *Node[V]
	size int
}

func (l *LinkedList[V]) Head() *Node[V] {
	return l.head.next
}

func (l *LinkedList[V]) Tail() *Node[V] {
	return l.tail.prev
}

func (l *LinkedList[V]) Append(value V) *Node[V] {
	node := NewNode(value, l.tail.prev, l.tail)
	l.tail.prev.next = node
	l.tail.prev = node
	l.size++
	return node
}

func (l *LinkedList[V]) Prepend(value V) *Node[V] {
	node := NewNode(value, l.head, l.head.next)
	l.head.next.prev = node
	l.head.next = node
	l.size++
	return node
}

func (l *LinkedList[V]) RemoveNode(node *Node[V]) bool {
	if node == nil || node == l.head || node == l.tail {
		return false
	}

	node.prev.next = node.next
	node.next.prev = node.prev
	l.size--
	return true
}

func (l *LinkedList[V]) IsNil() bool {
	return l.head.next == l.tail
}

func (l *LinkedList[V]) Iter() chan *Node[V] {
	ch := make(chan *Node[V])
	go func() {
		defer close(ch)
		node := l.head.next
		for node != l.tail {
			ch <- node
			node = node.next
		}
	}()

	return ch
}
func (l *LinkedList[V]) Size() int {
	return l.size
}

func NewNode[V comparable](value V, prev, next *Node[V]) *Node[V] {
	return &Node[V]{
		value: value,
		prev:  prev,
		next:  next,
	}
}

type Node[V comparable] struct {
	value V
	prev  *Node[V]
	next  *Node[V]
}

func (n *Node[V]) Value() V {
	return n.value
}

func (n *Node[V]) Next() *Node[V] {
	return n.next
}
func (n *Node[V]) Prev() *Node[V] {
	return n.prev
}
