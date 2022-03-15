package main

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
	v := &Vertex{X: 15, Y: 25, Z:44}
	Len := cpsv.GetSize(Vertex{})
	testBytes := &binary{
		addr: uintptr(unsafe.Pointer(v)),
		cap:  int(Len),
		len:  int(Len),
	}
	wbuf := *(*[]byte)(unsafe.Pointer(testBytes))
	fmt.Println(wbuf)
	cpsv.Store("d1", wbuf, int(Len), 0)

	fmt.Scanln()

	var readData []byte=make([]byte, Len)
	cpsv.Load("d1", &readData, 0, Len)
	fmt.Println(readData)
	var bufV *Vertex = *(**Vertex)(unsafe.Pointer(&readData))
	fmt.Printf("X: %d, Y:%d, Z:%d\n", bufV.X, bufV.Y, bufV.Z)
	
	fmt.Scanln()
	cpsv.Destroy()
}
