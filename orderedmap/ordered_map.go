package orderedmap

import linkedlist "github.com/dotdak/gopkg/linked_list"

func NewOrderedMap[K, V comparable]() *OrderedMap[K, V] {
	return &OrderedMap[K, V]{
		data: make(map[K]*Value[K, V]),
		list: linkedlist.NewLinkedList[K](),
	}
}

type Value[K, V comparable] struct {
	value   V
	address *linkedlist.Node[K]
}

type OrderedMap[K, V comparable] struct {
	data map[K]*Value[K, V]
	list *linkedlist.LinkedList[K]
}

func (m *OrderedMap[K, V]) Add(k K, v V) {
	if value, ok := m.data[k]; ok {
		m.list.RemoveNode(value.address)
	}

	m.data[k] = &Value[K, V]{
		value:   v,
		address: m.list.Append(k),
	}
}

func (m *OrderedMap[K, V]) Pop(k K) (v V) {
	if value, ok := m.data[k]; ok {
		m.list.RemoveNode(value.address)
		delete(m.data, k)
		return value.value
	}

	return
}

func (m *OrderedMap[K, V]) PopFirst() (k K, v V) {
	if m.IsEmpty() {
		return
	}
	head := m.list.Head()
	m.list.RemoveNode(head)
	value := m.data[head.Value()]
	delete(m.data, head.Value())
	return head.Value(), value.value
}

func (m *OrderedMap[K, V]) PopLast() (k K, v V) {
	if m.IsEmpty() {
		return
	}
	tail := m.list.Tail()
	m.list.RemoveNode(tail)
	value := m.data[tail.Value()]
	delete(m.data, tail.Value())
	return tail.Value(), value.value
}

func (m *OrderedMap[K, V]) Keys() []K {
	out := make([]K, 0, m.list.Size())
	for i := range m.list.Iter() {
		out = append(out, i.Value())
	}

	return out
}

func (m *OrderedMap[K, V]) IsEmpty() bool {
	return len(m.data) == 0 && m.list.IsNil()
}
