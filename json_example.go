package main

import "C"
import (
	"fmt"
	// "unsafe"
	"encoding/json"

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
	jsonLoc, _ := json.Marshal(loc)

	len := len(jsonLoc)
	fmt.Println(jsonLoc)

	cpsv.NonFixedStore("d1", jsonLoc, int(len))

	fmt.Scanln()

	readData, err := cpsv.NonFixedLoad("d1")

	if err == nil {
		fmt.Println(readData)
		var bufL Location
		json.Unmarshal(readData, &bufL)
		fmt.Printf("X: %d, Y:%d, Z: %d\n", bufL.X, bufL.Y, bufL.Z)
	} else {
		fmt.Println("got errors:", err)
	}

	fmt.Scanln()

	cpsv.Destroy()
}
