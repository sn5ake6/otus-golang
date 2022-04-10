package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
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

	t.Run("set and get only one element", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)
	})

	t.Run("clear when only one element", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		c.Clear()

		val, ok = c.Get("aaa")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("set more elements than capacity", func(t *testing.T) {
		c := NewCache(2)

		c.Set("one", 1)
		c.Set("two", 2)
		c.Set("three", 3)

		val, ok := c.Get("one")
		require.False(t, ok)
		require.Nil(t, val)

		val, ok = c.Get("two")
		require.True(t, ok)
		require.Equal(t, 2, val)

		val, ok = c.Get("three")
		require.True(t, ok)
		require.Equal(t, 3, val)
	})

	t.Run("LRU", func(t *testing.T) {
		c := NewCache(3)

		c.Set("one", 1)
		c.Set("two", 2)
		c.Set("three", 3)

		c.Get("one")
		c.Get("three")

		c.Set("four", 4)

		val, ok := c.Get("two")
		require.False(t, ok)
		require.Nil(t, val)

		val, ok = c.Get("one")
		require.True(t, ok)
		require.Equal(t, 1, val)

		val, ok = c.Get("three")
		require.True(t, ok)
		require.Equal(t, 3, val)

		val, ok = c.Get("four")
		require.True(t, ok)
		require.Equal(t, 4, val)
	})
}

func TestCacheMultithreading(t *testing.T) {
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
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}
