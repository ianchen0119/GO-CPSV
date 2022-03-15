package main

import "C"
import (
	"fmt"
	// "time"
	"unsafe"

	"github.com/ianchen0119/GO-CPSV/cpsv"
)

type Vertex struct {
	X int32
	Y int32
	Z int32
}

type binary struct {
	addr uintptr
	len  int
	cap  int
}

func test(val *[4096]byte){
	fmt.Println((*val)[0])
}

func main() {
	cpsv.Start()
	v := &Vertex{X: 15, Y: 23, Z: 16}
	Len := cpsv.GetSize(Vertex{})
	fmt.Println(Len)

	var wbuf [4096]byte
	byteSlice := C.GoBytes(unsafe.Pointer(v), C.int(Len))
	
	for i :=0;i<Len;i++ {
		wbuf[i] = byteSlice[i]
	}

	cpsv.Store("d1", wbuf, int(Len), 0)

	fmt.Scanln()

	readData := cpsv.Load("d1", 0, Len)
	var bufV *Vertex = *(**Vertex)(unsafe.Pointer(&readData))
	fmt.Scanln()

	fmt.Printf("X: %d, Y:%d, Z:%d\n", bufV.X, bufV.Y, bufV.Z)

	cpsv.Destroy()
}
