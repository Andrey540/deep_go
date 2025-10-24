package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homewrok_test.go

type treeNode[T comparable, V any] struct {
	key   T
	value V
	left  *treeNode[T, V]
	right *treeNode[T, V]
}

type SortFunc[T comparable] func(first, second T) int

type OrderedMap[T comparable, V any] struct {
	root   *treeNode[T, V]
	size   int
	sortFn SortFunc[T]
}

func NewOrderedMap[T comparable, V any](sortFn SortFunc[T]) OrderedMap[T, V] {
	return OrderedMap[T, V]{
		sortFn: sortFn,
	}
}

func (m *OrderedMap[T, V]) Insert(key T, value V) {
	m.size++
	if m.root == nil {
		m.root = &treeNode[T, V]{
			key: key,
		}
		return
	}
	m.insertImpl(m.root, key, value)
}

func (m *OrderedMap[T, V]) Erase(key T) {
	m.eraseImpl(&m.root, key)
}

func (m *OrderedMap[T, V]) Contains(key T) bool {
	return m.containsImpl(m.root, key)
}

func (m *OrderedMap[T, V]) Size() int {
	return m.size
}

func (m *OrderedMap[T, V]) ForEach(action func(T, V)) {
	m.forEachImpl(m.root, action)
}

func (m *OrderedMap[T, V]) eraseImpl(node **treeNode[T, V], key T) {
	if *node == nil {
		return
	}

	if (*node).key == key {
		m.size--
		if (*node).left == nil && (*node).right == nil {
			*node = nil
			return
		}
		if (*node).left != nil && (*node).right == nil {
			*node = (*node).left
			return
		}
		if (*node).right != nil && (*node).left == nil {
			*node = (*node).right
			return
		}

		minRight := (*node).right
		for minRight.left != nil {
			minRight = minRight.left
		}
		(*node).key = minRight.key
		(*node).value = minRight.value
		*minRight = *(minRight).right
		return
	}
	m.eraseImpl(&(*node).left, key)
	m.eraseImpl(&(*node).right, key)
}

func (m *OrderedMap[T, V]) containsImpl(node *treeNode[T, V], key T) bool {
	if node == nil {
		return false
	}
	if node.key == key {
		return true
	}
	return m.containsImpl(node.left, key) || m.containsImpl(node.right, key)
}

func (m *OrderedMap[T, V]) forEachImpl(node *treeNode[T, V], action func(T, V)) {
	if node == nil {
		return
	}
	m.forEachImpl(node.left, action)
	action(node.key, node.value)
	m.forEachImpl(node.right, action)
}

func (m *OrderedMap[T, V]) insertImpl(node *treeNode[T, V], key T, value V) {
	if m.sortFn(node.key, key) > 0 {
		if node.left == nil {
			node.left = &treeNode[T, V]{
				key:   key,
				value: value,
			}
		} else {
			m.insertImpl(node.left, key, value)
		}
	} else {
		if node.right == nil {
			node.right = &treeNode[T, V]{
				key:   key,
				value: value,
			}
		} else {
			m.insertImpl(node.right, key, value)
		}
	}
}

func TestOrderedMap(t *testing.T) {
	data := NewOrderedMap[int, int](func(first, second int) int {
		return first - second
	})
	assert.Zero(t, data.Size())

	data.Insert(10, 10)
	data.Insert(5, 5)
	data.Insert(15, 15)
	data.Insert(2, 2)
	data.Insert(4, 4)
	data.Insert(12, 12)
	data.Insert(14, 14)

	assert.Equal(t, 7, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(3))
	assert.False(t, data.Contains(13))

	var keys []int
	expectedKeys := []int{2, 4, 5, 10, 12, 14, 15}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))

	data.Erase(10)
	data.Erase(15)
	data.Erase(14)
	data.Erase(2)

	assert.Equal(t, 3, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(2))
	assert.False(t, data.Contains(14))

	keys = nil
	expectedKeys = []int{4, 5, 12}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))
}
