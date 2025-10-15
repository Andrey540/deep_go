package main

import (
	"fmt"
	"runtime"
	"time"
	"unsafe"
)

type Data struct {
	name string
}

func NewData(name string) *Data {
	data := &Data{
		name: name,
	}

	fmt.Println("created", name)
	runtime.SetFinalizer(data, func(ptr *Data) {
		fmt.Println("finalizer called on addr", unsafe.Pointer(ptr), "value is", ptr.name)
	})

	return data
}

func main() {
	data := NewData("__data__")
	_ = data

	runtime.GC()
	time.Sleep(time.Second)
}
