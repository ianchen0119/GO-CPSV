package main

import "C"
import (
	"fmt"
	"unsafe"

	"github.com/ianchen0119/GO-CPSV/cpsv"
)

type Location struct {
	X int32
	Y int32
	Z int32
}

func main() {
	cpsv.Start("safCkpt=TEST1,safApp=safCkptService")
	v := &Location{X: 25, Y: 23, Z: 50}
	len := cpsv.GetSize(Location{})
	wbuf := cpsv.GoBytes(unsafe.Pointer(v), len)

	cpsv.NonFixedStore("d1", wbuf, int(len))

	fmt.Scanln()

	readData, err := cpsv.NonFixedLoad("d1")

	if err == nil {
		var bufV *Location = *(**Location)(unsafe.Pointer(&readData))
		fmt.Printf("X: %d, Y:%d, Z: %d\n", bufV.X, bufV.Y, bufV.Z)
	} else {
		fmt.Println("got errors:", err)
	}

	fmt.Scanln()

	cpsv.Destroy()
}
