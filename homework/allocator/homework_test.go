package main

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

func Defragment(memory []byte, pointers []unsafe.Pointer) {
	if len(pointers) == 0 {
		return
	}
	index := 0
	indexes := make([]int, 0, len(pointers))
	for _, p := range pointers {
		for index < len(memory) {
			if unsafe.Pointer(&memory[index]) == p {
				indexes = append(indexes, index)
				index++
				break
			}
			index++
		}
	}
	fragmentedIndex := 0
	for i, v := range memory {
		if fragmentedIndex == len(pointers) {
			return
		}
		if i != indexes[fragmentedIndex] {
			memory[i] = memory[indexes[fragmentedIndex]]
			pointers[fragmentedIndex] = unsafe.Pointer(&memory[i])
			memory[indexes[fragmentedIndex]] = v
			fragmentedIndex++
		} else if i == indexes[fragmentedIndex] {
			fragmentedIndex++
		}
	}
}

func TestDefragmentation(t *testing.T) {
	var fragmentedMemory = []byte{
		0xFF, 0x00, 0x00, 0x00,
		0x00, 0xFF, 0x00, 0x00,
		0x00, 0x00, 0xFF, 0x00,
		0x00, 0x00, 0x00, 0xFF,
	}

	var fragmentedPointers = []unsafe.Pointer{
		unsafe.Pointer(&fragmentedMemory[0]),
		unsafe.Pointer(&fragmentedMemory[5]),
		unsafe.Pointer(&fragmentedMemory[10]),
		unsafe.Pointer(&fragmentedMemory[15]),
	}

	var defragmentedPointers = []unsafe.Pointer{
		unsafe.Pointer(&fragmentedMemory[0]),
		unsafe.Pointer(&fragmentedMemory[1]),
		unsafe.Pointer(&fragmentedMemory[2]),
		unsafe.Pointer(&fragmentedMemory[3]),
	}

	var defragmentedMemory = []byte{
		0xFF, 0xFF, 0xFF, 0xFF,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
	}

	Defragment(fragmentedMemory, fragmentedPointers)
	assert.True(t, reflect.DeepEqual(defragmentedMemory, fragmentedMemory))
	assert.True(t, reflect.DeepEqual(defragmentedPointers, fragmentedPointers))
}
