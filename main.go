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
	loc := &Location{X: 25, Y: 23, Z: 50}
	len := cpsv.GetSize(Location{})
	wbuf := cpsv.GoBytes(unsafe.Pointer(loc), len)

	cpsv.NonFixedStore("d1", wbuf, int(len))

	fmt.Scanln()

	readData, err := cpsv.NonFixedLoad("d1")

	if err == nil {
		var bufL *Location = *(**Location)(unsafe.Pointer(&readData))
		fmt.Printf("X: %d, Y:%d, Z: %d\n", bufL.X, bufL.Y, bufL.Z)
	} else {
		fmt.Println("got errors:", err)
	}

	fmt.Scanln()

	cpsv.Destroy()
}
