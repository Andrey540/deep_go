package main

import (
	"reflect"
	"runtime"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

type COWBuffer struct {
	data   []byte
	refs   *int
	closed bool
}

func NewCOWBuffer(data []byte) COWBuffer {
	result := COWBuffer{
		data: data,
		refs: new(int),
	}
	*result.refs = 1
	runtime.SetFinalizer(&result, func(b *COWBuffer) {
		b.Close()
	})
	return result
}

func (b *COWBuffer) Clone() COWBuffer {
	*b.refs++
	return *b
}

func (b *COWBuffer) Close() {
	if !b.closed {
		*b.refs--
		runtime.SetFinalizer(b, nil)
		b.closed = true
	}
}

func (b *COWBuffer) Update(index int, value byte) bool {
	if index < 0 || index > len(b.data)-1 {
		return false
	}
	if *b.refs == 1 {
		b.data[index] = value
		return true
	}

	*b.refs--

	newData := make([]byte, len(b.data))
	copy(newData, b.data)
	newData[index] = value

	b.data = newData
	b.refs = new(int)
	*b.refs = 1
	return true
}

func (b *COWBuffer) String() string {
	return *(*string)(unsafe.Pointer(&b.data))
}

func TestCOWBuffer(t *testing.T) {
	data := []byte{'a', 'b', 'c', 'd'}
	buffer := NewCOWBuffer(data)
	assert.Equal(t, *buffer.refs, 1)
	defer buffer.Close()

	copy1 := buffer.Clone()
	copy2 := buffer.Clone()

	assert.Equal(t, *buffer.refs, 3)
	assert.Equal(t, *copy1.refs, 3)
	assert.Equal(t, *copy2.refs, 3)

	assert.Equal(t, unsafe.SliceData(data), unsafe.SliceData(buffer.data))
	assert.Equal(t, unsafe.SliceData(buffer.data), unsafe.SliceData(copy1.data))
	assert.Equal(t, unsafe.SliceData(copy1.data), unsafe.SliceData(copy2.data))

	assert.True(t, (*byte)(unsafe.SliceData(data)) == unsafe.StringData(buffer.String()))
	assert.True(t, (*byte)(unsafe.StringData(buffer.String())) == unsafe.StringData(copy1.String()))
	assert.True(t, (*byte)(unsafe.StringData(copy1.String())) == unsafe.StringData(copy2.String()))

	assert.True(t, buffer.Update(0, 'g'))
	assert.False(t, buffer.Update(-1, 'g'))
	assert.False(t, buffer.Update(4, 'g'))

	assert.Equal(t, *buffer.refs, 1)
	assert.Equal(t, *copy1.refs, 2)
	assert.Equal(t, *copy2.refs, 2)

	assert.True(t, reflect.DeepEqual([]byte{'g', 'b', 'c', 'd'}, buffer.data))
	assert.True(t, reflect.DeepEqual([]byte{'a', 'b', 'c', 'd'}, copy1.data))
	assert.True(t, reflect.DeepEqual([]byte{'a', 'b', 'c', 'd'}, copy2.data))

	assert.NotEqual(t, unsafe.SliceData(buffer.data), unsafe.SliceData(copy1.data))
	assert.Equal(t, unsafe.SliceData(copy1.data), unsafe.SliceData(copy2.data))

	copy1.Close()
	assert.Equal(t, *copy1.refs, 1)

	previous := copy2.data
	copy2.Update(0, 'f')
	current := copy2.data

	// 1 reference - don't need to copy buffer during update
	assert.Equal(t, unsafe.SliceData(previous), unsafe.SliceData(current))

	copy2.Close()
	assert.Equal(t, *copy2.refs, 0)

	buffer.Close()
	assert.Equal(t, *buffer.refs, 0)
	buffer.Close()
	assert.Equal(t, *buffer.refs, 0)

	newBuffer := NewCOWBuffer(data)
	newBuffer1 := newBuffer.Clone()
	assert.Equal(t, *newBuffer1.refs, 2)
	newBuffer = COWBuffer{}
	runtime.GC()

	assert.Equal(t, *newBuffer1.refs, 1)
}
