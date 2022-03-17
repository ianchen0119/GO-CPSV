package main

import "C"
import (
	"fmt"
	"unsafe"

	"github.com/ianchen0119/GO-CPSV/cpsv"
)

func main() {
	cpsv.Start("safCkpt=TEST1,safApp=safCkptService")
	len := cpsv.GetSize(Vertex{})

	fmt.Scanln()

	readData,err := cpsv.Load("d1", 0, len)
	if err == nil {
		var bufV *Vertex = *(**Vertex)(unsafe.Pointer(&readData))
		fmt.Scanln()

		fmt.Printf("X: %d, Y:%d\n", bufV.X, bufV.Y)
	} else {
		fmt.Println(err)
		fmt.Scanln()
	}
	cpsv.Destroy()
}
