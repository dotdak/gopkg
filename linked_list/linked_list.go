package linkedlist

import (
	"fmt"
)

func NewNode(val int) Node {
	return &NodeImpl{
		val:  val,
		next: nil,
	}
}

type Node interface {
	Val() int
	Next() Node
	SetNext(n Node)
}

type NodeImpl struct {
	next Node
	val  int
}

func (n *NodeImpl) Val() int {
	return n.val
}

func (n *NodeImpl) Next() Node {
	return n.next
}

func (n *NodeImpl) SetNext(node Node) {
	n.next = node
}

func End(head Node) Node {
	if head == nil {
		return nil
	}
	for head.Next() != nil {
		head = head.Next()
	}
	return head
}

func Get(head Node, index int) Node {
	if head == nil {
		return nil
	}
	inx := 0
	for head != nil {
		if inx == index {
			return head
		}
		inx += 1
		head = head.Next()
	}
	return nil
}

func Insert(head Node, index, val int) bool {
	p := Get(head, index)
	if p == nil {
		return false
	}

	new := &NodeImpl{
		val:  val,
		next: p.Next(),
	}
	p.SetNext(new)
	return true
}

func Append(head Node, val ...int) Node {
	for _, v := range val {
		if head == nil {
			head = NewNode(v)
		}

		p := End(head)
		p.SetNext(NewNode(v))
	}

	return head
}

func Print(head Node) {
	for head != nil {
		fmt.Print("->", head.Val())
		head = head.Next()
	}
}

func main() {
	a := NewNode(19)
	Append(a, 10, 4, 1, 3, 45, 56)
	Print(a)
}
