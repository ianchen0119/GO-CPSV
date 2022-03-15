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

func main() {
	cpsv.Start()
	v := &Vertex{X: 15, Y: 23, Z: 14}
	Len := cpsv.GetSize(Vertex{})
	fmt.Println(Len)

	wbuf := C.GoBytes(unsafe.Pointer(v), C.int(Len))

	fmt.Println(wbuf)
	cpsv.Store("d1", wbuf, int(Len), 0)

	fmt.Scanln()

	readData := cpsv.Load("d1", 0, Len)
	var bufV *Vertex = *(**Vertex)(unsafe.Pointer(&readData))
	fmt.Printf("X: %d, Y:%d, Z:%d\n", bufV.X, bufV.Y, bufV.Z)

	fmt.Scanln()
	cpsv.Destroy()
}
