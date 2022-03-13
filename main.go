package main

import (
	"fmt"
	"time"
	"unsafe"

	"github.com/ianchen0119/GO-CPSV/cpsv"
)

type Vertex struct {
	X int32
	Y int32
}

type binary struct {
	addr uintptr
	len  int
	cap  int
}

func main() {
	cpsv.Start()
	v := &Vertex{X: 1, Y: 2}
	Len := unsafe.Sizeof(*v)
	testBytes := &binary{
		addr: uintptr(unsafe.Pointer(v)),
		cap:  int(Len),
		len:  int(Len),
	}
	fmt.Printf("Len: %d\n", int(Len))
	wbuf := *(*[]byte)(unsafe.Pointer(testBytes))
	fmt.Println(wbuf)
	cpsv.Store("data-1", wbuf, int(Len), 0)

	time.Sleep(2 * time.Second)

	var readData []byte
	//bufPtr := (*[]byte)(unsafe.Pointer(&readData))
	cpsv.Load(&readData, 0, int(Len))
	fmt.Printf("%v\n", readData)
}
