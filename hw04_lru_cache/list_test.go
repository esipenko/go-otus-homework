package hw04lrucache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func getListEntries(l List) []int {
	elems := make([]int, 0, l.Len())

	for i := l.Front(); i != nil; i = i.Next {
		elems = append(elems, i.Value.(int))
	}

	return elems
}

func TestPush(t *testing.T) {
	t.Run("push front", func(t *testing.T) {
		l := NewList()
		l.PushFront(10)
		l.PushFront(20)

		require.Equal(t, 2, l.Len())
		require.Equal(t, []int{20, 10}, getListEntries(l))
	})

	t.Run("push back", func(t *testing.T) {
		l := NewList()

		l.PushBack(10)
		l.PushBack(20)

		require.Equal(t, 2, l.Len())

		require.Equal(t, []int{10, 20}, getListEntries(l))
	})
}

func TestRemove(t *testing.T) {
	t.Run("remove from len 1", func(t *testing.T) {
		l := NewList()

		l.PushBack(10)
		li := l.Front()

		l.Remove(li)

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Back())
		require.Nil(t, l.Front())
	})

	t.Run("remove from front", func(t *testing.T) {
		l := NewList()

		l.PushBack(10)
		l.PushBack(20)

		li := l.Front()

		l.Remove(li)

		require.Equal(t, 1, l.Len())
		require.Equal(t, l.Back(), l.Front())
	})

	t.Run("remove from back", func(t *testing.T) {
		l := NewList()

		l.PushBack(10)
		l.PushBack(20)

		li := l.Back()

		l.Remove(li)

		require.Equal(t, 1, l.Len())
		require.Equal(t, l.Back(), l.Front())
	})

	t.Run("remove from middle", func(t *testing.T) {
		l := NewList()

		l.PushBack(10)
		l.PushBack(20)
		l.PushBack(30)

		li := l.Front().Next

		l.Remove(li)

		require.Equal(t, 2, l.Len())
		require.Equal(t, []int{10, 30}, getListEntries(l))
	})
}

func TestMoveFront(t *testing.T) {
	t.Run("move front front", func(t *testing.T) {
		l := NewList()

		l.PushBack(10)
		l.PushBack(20)
		l.PushBack(30)

		li := l.Front()
		l.MoveToFront(li)

		require.Equal(t, []int{10, 20, 30}, getListEntries(l))
	})
	t.Run("move front middle", func(t *testing.T) {
		l := NewList()

		l.PushBack(10)
		l.PushBack(20)
		l.PushBack(30)

		li := l.Front().Next

		l.MoveToFront(li)
		require.Equal(t, []int{20, 10, 30}, getListEntries(l))
	})

	t.Run("move front back", func(t *testing.T) {
		l := NewList()

		l.PushBack(10)
		l.PushBack(20)
		l.PushBack(30)

		li := l.Back()

		l.MoveToFront(li)

		require.Equal(t, []int{30, 10, 20}, getListEntries(l))
		require.Equal(t, l.Back().Value, 20)
	})
}

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, 3, l.Len())

		middle := l.Front().Next // 20
		l.Remove(middle)         // [10, 30]
		require.Equal(t, 2, l.Len())

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, 7, l.Len())
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]

		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, getListEntries(l))
	})
}
