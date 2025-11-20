package main

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

func Trace(stacks [][]uintptr) []uintptr {
	var result []uintptr
	pointersMap := make(map[uintptr]struct{})
	var zeroValue uintptr
	for _, stack := range stacks {
		for _, p := range stack {
			if p != zeroValue {
				getPointer(p, pointersMap, &result)
			}
		}
	}
	return result
}

func getPointer(pointer uintptr, pointersMap map[uintptr]struct{}, pointers *[]uintptr) {
	if _, ok := pointersMap[pointer]; ok {
		return
	}
	var zeroValue uintptr
	if pointer != zeroValue {
		*pointers = append(*pointers, pointer)
		p := *(*uintptr)(unsafe.Pointer(pointer))
		pointersMap[pointer] = struct{}{}
		if p != zeroValue {
			getPointer(p, pointersMap, pointers)
		}
	}
}

func TestTrace(t *testing.T) {
	var heapObjects = []int{
		0x00, 0x00, 0x00, 0x00, 0x00,
	}

	var heapPointer1 *int = &heapObjects[1]
	var heapPointer2 *int = &heapObjects[2]
	var heapPointer3 *int = nil
	var heapPointer4 **int = &heapPointer3

	p1 := uintptr(unsafe.Pointer(&heapPointer1))
	p2 := uintptr(unsafe.Pointer(&heapObjects[0]))
	p3 := uintptr(unsafe.Pointer(&heapPointer2))
	p4 := uintptr(unsafe.Pointer(&heapObjects[1]))
	p5 := uintptr(unsafe.Pointer(&heapObjects[2]))
	p6 := uintptr(unsafe.Pointer(&heapPointer4))
	p7 := uintptr(unsafe.Pointer(&heapPointer3))
	p8 := uintptr(unsafe.Pointer(&heapObjects[3]))

	var stacks = [][]uintptr{
		{
			p1, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, p2,
			0x00, 0x00, 0x00, 0x00,
		},
		{
			p3, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, p4,
			0x00, 0x00, 0x00, p5,
			p6, 0x00, 0x00, 0x00,
		},
		{
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, p8,
		},
	}

	pointers := Trace(stacks)
	expectedPointers := []uintptr{
		p1,
		p4,
		p2,
		p3,
		p5,
		p6,
		p7,
		p8,
	}

	assert.True(t, reflect.DeepEqual(expectedPointers, pointers))
}
