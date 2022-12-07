package main

import "C"
import (
	"fmt"
	"time"

	"encoding/json"

	"github.com/ianchen0119/GO-CPSV/cpsv"
)

type Location struct {
	X int32
	Y int32
	Z int32
}

func main() {
	cpsv.Start("safCkpt=TEST2,safApp=safCkptService")
	loc := Location{X: 25, Y: 23, Z: 50}
	jsonLoc, _ := json.Marshal(loc)

	len := len(jsonLoc)
	fmt.Println(jsonLoc)

	cpsv.NonFixedStore("json-data", jsonLoc, int(len))

	time.Sleep(3 * time.Second)

	readData, err := cpsv.NonFixedLoad("json-data")

	if err == nil {
		fmt.Println(readData)
		var bufL Location
		json.Unmarshal(readData, &bufL)
		fmt.Printf("X: %d, Y:%d, Z: %d\n", bufL.X, bufL.Y, bufL.Z)
	} else {
		fmt.Println("got errors:", err)
	}

	cpsv.Destroy()
}
