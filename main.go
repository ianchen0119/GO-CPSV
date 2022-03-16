package main

import "C"
import (
	"fmt"
	"unsafe"

	"github.com/ianchen0119/GO-CPSV/cpsv"
)

type Vertex struct {
	X int32
	Y int32
}

func main() {
	cpsv.Start()
	v := &Vertex{X: 15, Y: 23}
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

	fmt.Printf("X: %d, Y:%d\n", bufV.X, bufV.Y)

	cpsv.Destroy()
}
