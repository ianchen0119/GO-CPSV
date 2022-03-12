package cpsv

import "C"

const (
	Async int = 0
	Sync  int = 1
)

type req struct {
	sectionId string
	data      *C.char
	offset    int
	reqType   int
}

type EventQ struct {
	queue chan req
}
