package hw04lrucache

import (
	"github.com/stretchr/testify/require"
	"math/rand"
	"strconv"
	"sync"
	"testing"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("overflow capacity removes last added el", func(t *testing.T) {
		c := NewCache(3)

		c.Set("5", 5)
		c.Set("4", 4)
		c.Set("3", 3)

		c.Set("2", 2)

		val, ok := c.Get("5")
		require.False(t, ok)
		require.Nil(t, val)

		val, ok = c.Get("4")

		require.True(t, ok)
		require.Equal(t, val, 4)
	})

	t.Run("overflow capacity removes last used el", func(t *testing.T) {
		c := NewCache(3)

		c.Set("5", 5)
		c.Set("4", 4)
		c.Set("3", 3)
		c.Set("5", 50)

		val, ok := c.Get("5")
		require.True(t, ok)
		require.Equal(t, val, 50)

		c.Set("4", 40)
		c.Set("2", 20)

		val, ok = c.Get("3")

		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("cache clears", func(t *testing.T) {
		c := NewCache(3)

		c.Set("5", 5)
		c.Set("4", 4)
		c.Set("3", 3)

		c.Clear()
		val, ok := c.Get("5")
		require.False(t, ok)
		require.Nil(t, val)

		val, ok = c.Get("4")
		require.False(t, ok)
		require.Nil(t, val)

		val, ok = c.Get("4")
		require.False(t, ok)
		require.Nil(t, val)
	})
}

func TestCacheMultithreading(t *testing.T) {
	t.Run("it can work concurrently, and save cache capacity elements", func(t *testing.T) {
		c := NewCache(10)
		wg := &sync.WaitGroup{}
		wg.Add(2)

		go func() {
			defer wg.Done()
			for i := 0; i < 1_000_000; i++ {
				c.Set(Key(strconv.Itoa(i)), i)
			}
		}()

		go func() {
			defer wg.Done()
			for i := 0; i < 1_000_000; i++ {
				k := Key(strconv.Itoa(rand.Intn(1_000_000)))
				c.Get(k)
			}
		}()

		wg.Wait()
		storedAmount := 0

		for i := 0; i < 1_000_000; i++ {
			val, ok := c.Get(Key(strconv.Itoa(i)))

			if ok {
				require.Equal(t, i, val)
				storedAmount++
			}
		}

		require.Equal(t, storedAmount, 10)
	})

	t.Run("it can work concurrently, and can clear cache", func(t *testing.T) {
		threshold := 500_000
		c := NewCache(1_000_000)

		wg := &sync.WaitGroup{}
		wg.Add(3)

		savedThreshold := make(chan interface{}, 1)

		// Write to cache same amount of elements as its capacity
		go func() {
			defer wg.Done()
			for i := 0; i < 1_000_000; i++ {
				c.Set(Key(strconv.Itoa(i)), i)
				// When we get to threshold, i want to clear cache in another goroutine
				if i == threshold {
					savedThreshold <- true
				}
			}
		}()

		go func() {
			defer wg.Done()
			for i := 0; i < 1_000_000; i++ {
				k := Key(strconv.Itoa(rand.Intn(1_000_000)))
				c.Get(k)
			}
		}()

		go func() {
			_, ok := <-savedThreshold

			if ok {
				c.Clear()
				wg.Done()
			}
		}()

		wg.Wait()
		storedAmount := 0

		for i := 0; i < 1_000_000; i++ {
			val, ok := c.Get(Key(strconv.Itoa(i)))

			if ok {
				require.Equal(t, i, val)
				storedAmount++
			}
		}

		// As we do it concourently, we actualy dont know, how many elements left in cache
		// It cannot be more than 500_000
		require.Less(t, storedAmount, threshold)
	})
}
