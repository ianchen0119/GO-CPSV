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
	len := cpsv.GetSize(Vertex{})

	fmt.Scanln()

	readData := cpsv.Load("d1", 0, len)
	var bufV *Vertex = *(**Vertex)(unsafe.Pointer(&readData))
	fmt.Printf("X: %d, Y:%d\n", bufV.X, bufV.Y)

	fmt.Scanln()
	cpsv.Destroy()
}
