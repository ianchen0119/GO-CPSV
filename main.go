package main

import (
	"fmt"
	"time"
	"unsafe"

	"github.com/ianchen0119/GO-CPSV/cpsv"
)

type Vertex struct {
	X int
	Y int
}

func main() {
	cpsv.Start()
	var p *Vertex
	v := Vertex{1, 2}
	p = &v
	p.X = 30
	fmt.Println(v)
	byteData := []byte(fmt.Sprintf("%v", v))
	cpsv.Store("data-1", byteData, 0)

	time.Sleep(8 * time.Second)

	var readData []byte
	cpsv.Load(&readData, 0, 50)
	var newV Vertex = *(*Vertex)(unsafe.Pointer(&readData))
	fmt.Println(newV.X)
}
