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
	front *ListItem
	back  *ListItem
	len   int
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	li := &ListItem{Value: v, Next: l.front}
	l.len++

	if l.front != nil {
		l.front.Prev = li
	}

	l.front = li

	if l.back == nil {
		l.back = li
	}

	return li
}

func (l *list) PushBack(v interface{}) *ListItem {
	li := &ListItem{Value: v, Prev: l.back}

	l.len++

	if l.back != nil {
		l.back.Next = li
	}

	l.back = li

	if l.front == nil {
		l.front = li
	}

	return li
}

func (l *list) Remove(v *ListItem) {
	l.len--

	if v.Prev != nil {
		v.Prev.Next = v.Next
	} else {
		l.front = v.Next
	}

	if v.Next != nil {
		v.Next.Prev = v.Prev
	} else {
		l.back = v.Prev
	}
}

func (l *list) MoveToFront(v *ListItem) {
	if v.Prev == nil {
		return
	}

	v.Prev.Next = v.Next

	if v.Next != nil {
		v.Next.Prev = v.Prev
	} else {
		l.back = v.Prev
	}

	if l.front != nil {
		l.front.Prev = v
	}

	v.Next = l.front
	v.Prev = nil

	l.front = v
}

func NewList() List {
	return new(list)
}
