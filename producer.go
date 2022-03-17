package main

import "C"
import (
	"fmt"
	"unsafe"

	"github.com/ianchen0119/GO-CPSV/cpsv"
)

func main() {
	cpsv.Start("safCkpt=TEST1,safApp=safCkptService")
	v := &Vertex{X: 15, Y: 23}
	len := cpsv.GetSize(Vertex{})
	wbuf := C.GoBytes(unsafe.Pointer(v), C.int(len))

	cpsv.Store("d1", wbuf, int(len), 0)

	fmt.Scanln()

	readData, err := cpsv.Load("d1", 0, len)
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
