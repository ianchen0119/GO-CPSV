package main

import (
	"fmt"
	"time"
	"unsafe"
	"reflect"

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

func getSize(i interface{}) int {
	size := reflect.TypeOf(i).Size()
	return int(size)
}

func main() {
	cpsv.Start()
	v := &Vertex{X: 15, Y: 25, Z:44}
	Len := getSize(Vertex{})
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

	var readData []byte=make([]byte, Len)
	len := cpsv.Load(&readData, 0, Len)
	fmt.Println(len)
	fmt.Printf("%v\n", readData)
	var bufV *Vertex = *(**Vertex)(unsafe.Pointer(&readData))
	fmt.Printf("X: %d, Y:%d, Z:%d\n", bufV.X, bufV.Y, bufV.Z)
	
}
