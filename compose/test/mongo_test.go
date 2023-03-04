package main

import (
	"github.com/ianchen0119/GO-CPSV/cpsv"
)

func main() {
	defer cpsv.Destroy()
	cpsv.SetStorageProvider("test", "mongo")
	cpsv.NonFixedStore("test", []byte("This is Mongo Client"), 0)

	// load data
	data, _ := cpsv.NonFixedLoad("test")
	println(string(data))
}
