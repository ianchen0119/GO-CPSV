package cpsv

import "C"

const (
	NonFixed int = 0
	Fixed    int = 1
)

type req struct {
	sectionId string
	data      []byte
	size      int
	offset    int
	reqType   int
	resend    int
}

type eventQ chan req
