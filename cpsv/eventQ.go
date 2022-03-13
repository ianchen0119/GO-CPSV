package cpsv

import "fmt"

var q eventQ

func eventQInit() {
	q.queue = make(chan req, 100)
}

func (q *eventQ) pull(req *req) {
	val, ok := <-q.queue
	if ok {
		*req = val
	} else {
		fmt.Println("No value was read.")
	}
}

func (q *eventQ) push(req req) {
	q.queue <- req
}
