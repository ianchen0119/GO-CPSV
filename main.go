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
	cpsv.Start("safCkpt=TEST1,safApp=safCkptService")
	v := &Vertex{X: 25, Y: 23}
	len := cpsv.GetSize(Vertex{})
	wbuf := cpsv.GoBytes(unsafe.Pointer(v), len)

	cpsv.NonFixedStore("d1", wbuf, int(len))

	fmt.Scanln()

	readData, err := cpsv.NonFixedLoad("d1")

	if err == nil {
		var bufV *Vertex = *(**Vertex)(unsafe.Pointer(&readData))
		fmt.Printf("X: %d, Y:%d\n", bufV.X, bufV.Y)
	} else {
		fmt.Println("got errors:", err)
	}

	fmt.Scanln()

	cpsv.Destroy()
}
