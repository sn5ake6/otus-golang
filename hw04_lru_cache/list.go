package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	First *ListItem
	Last  *ListItem
	len   int
}

func (l list) Len() int {
	return l.len
}

func (l list) Front() *ListItem {
	return l.First
}

func (l list) Back() *ListItem {
	return l.Last
}

func (l *list) PushFront(v interface{}) *ListItem {
	i := &ListItem{v, nil, nil}
	if l.First == nil {
		l.len++
		l.First = i
		l.Last = i
	} else {
		l.insertBefore(l.First, i)
	}

	return i
}

func (l *list) insertBefore(node *ListItem, i *ListItem) {
	l.len++

	i.Next = node
	if node.Prev == nil {
		i.Prev = nil
		l.First = i
	} else {
		i.Prev = node.Prev
		node.Prev.Next = i
	}
	node.Prev = i
}

func (l *list) PushBack(v interface{}) *ListItem {
	if l.Last == nil {
		return l.PushFront(v)
	}

	NewNode := &ListItem{v, nil, nil}
	l.insertAfter(l.Last, NewNode)

	return NewNode
}

func (l *list) insertAfter(node *ListItem, i *ListItem) {
	l.len++

	i.Prev = node
	if node.Next == nil {
		i.Next = nil
		l.Last = i
	} else {
		i.Next = node.Next
		node.Next.Prev = i
	}
	node.Next = i
}

func (l *list) Remove(i *ListItem) {
	if i == nil {
		return
	}

	l.len--

	if i.Prev == nil {
		l.First = i.Next
	} else {
		i.Prev.Next = i.Next
	}

	if i.Next == nil {
		l.Last = i.Prev
	} else {
		i.Next.Prev = i.Prev
	}
}

func (l *list) MoveToFront(i *ListItem) {
	if i == nil {
		return
	}

	if l.len > 1 {
		l.Remove(i)
		l.insertBefore(l.First, i)
	}
}

func NewList() List {
	return new(list)
}
