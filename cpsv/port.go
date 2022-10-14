package cpsv

import "unsafe"

/* for all of storage solutions */
type genericAPI interface {
	destroy()
	load(string, uint32, int) ([]byte, error)
	store(string, []byte, int, int)
}

/* for openSAF checkpoint service */
type customAPI interface {
	nonFixedStore(string, []byte, int)
	nonFixedLoad(string) ([]byte, error)
	getSize(interface{}) int
	goBytes(unsafe.Pointer, int) []byte
}

type storageAPI interface {
	genericAPI
	customAPI
}
