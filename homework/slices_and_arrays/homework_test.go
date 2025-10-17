package main

import (
	"reflect"
	"testing"

	"golang.org/x/exp/constraints"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type CircularQueue[T constraints.Signed] struct {
	values []T
	len    int
	front  int
	rear   int
}

func NewCircularQueue[T constraints.Signed](size int) CircularQueue[T] {
	return CircularQueue[T]{
		len:    0,
		values: make([]T, size),
		front:  -1,
		rear:   -1,
	}
}

func (q *CircularQueue[T]) Push(value T) bool {
	if q.Full() {
		return false
	}
	if q.len == 0 {
		q.front = 0
	}
	q.len++
	q.rear = (q.rear + 1) % cap(q.values)
	q.values[q.rear] = value
	return true
}

func (q *CircularQueue[T]) Pop() bool {
	if q.Empty() {
		return false
	}
	q.values[q.front] = 0
	q.front = (q.front + 1) % cap(q.values)
	if q.Empty() {
		q.front = -1
		q.rear = -1
	}
	q.len--
	return true
}

func (q *CircularQueue[T]) Front() T {
	if q.Empty() {
		return -1
	}
	return q.values[q.front]
}

func (q *CircularQueue[T]) Back() T {
	if q.Empty() {
		return -1
	}
	return q.values[q.rear]
}

func (q *CircularQueue[T]) Empty() bool {
	return q.len == 0
}

func (q *CircularQueue[T]) Full() bool {
	return cap(q.values) == q.len
}

func TestCircularQueue(t *testing.T) {
	const queueSize = 3
	queue := NewCircularQueue[int](queueSize)

	assert.True(t, queue.Empty())
	assert.False(t, queue.Full())

	assert.Equal(t, -1, queue.Front())
	assert.Equal(t, -1, queue.Back())
	assert.False(t, queue.Pop())

	assert.True(t, queue.Push(1))
	assert.True(t, queue.Push(2))
	assert.True(t, queue.Push(3))
	assert.False(t, queue.Push(4))

	assert.True(t, reflect.DeepEqual([]int{1, 2, 3}, queue.values))

	assert.False(t, queue.Empty())
	assert.True(t, queue.Full())

	assert.Equal(t, 1, queue.Front())
	assert.Equal(t, 3, queue.Back())

	assert.True(t, queue.Pop())
	assert.False(t, queue.Empty())
	assert.False(t, queue.Full())
	assert.True(t, queue.Push(4))

	assert.True(t, reflect.DeepEqual([]int{4, 2, 3}, queue.values))

	assert.Equal(t, 2, queue.Front())
	assert.Equal(t, 4, queue.Back())

	assert.True(t, queue.Pop())

	assert.True(t, reflect.DeepEqual([]int{4, 0, 3}, queue.values))

	assert.Equal(t, 3, queue.Front())
	assert.Equal(t, 4, queue.Back())

	assert.True(t, queue.Pop())

	assert.True(t, reflect.DeepEqual([]int{4, 0, 0}, queue.values))

	assert.Equal(t, 4, queue.Front())
	assert.Equal(t, 4, queue.Back())

	assert.True(t, queue.Push(5))

	assert.True(t, reflect.DeepEqual([]int{4, 5, 0}, queue.values))

	assert.False(t, queue.Empty())
	assert.False(t, queue.Full())

	assert.True(t, queue.Push(6))

	assert.True(t, reflect.DeepEqual([]int{4, 5, 6}, queue.values))

	assert.False(t, queue.Empty())
	assert.True(t, queue.Full())

	assert.True(t, queue.Pop())

	assert.True(t, reflect.DeepEqual([]int{0, 5, 6}, queue.values))

	assert.Equal(t, 5, queue.Front())
	assert.Equal(t, 6, queue.Back())

	assert.True(t, queue.Pop())

	assert.True(t, reflect.DeepEqual([]int{0, 0, 6}, queue.values))

	assert.Equal(t, 6, queue.Front())
	assert.Equal(t, 6, queue.Back())

	assert.True(t, queue.Pop())

	assert.True(t, reflect.DeepEqual([]int{0, 0, 0}, queue.values))

	assert.Equal(t, -1, queue.Front())
	assert.Equal(t, -1, queue.Back())

	assert.False(t, queue.Pop())

	assert.True(t, queue.Empty())
	assert.False(t, queue.Full())
}
