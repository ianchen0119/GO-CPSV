package cpsv

import "fmt"

var q EventQ

func eventQInit() {
	q.queue = make(chan req, 100)
}

func (q *EventQ) pull(req *req) {
	val, ok := <-q.queue
	if ok {
		*req = val
	} else {
		fmt.Println("No value was read.")
	}
}

func (q *EventQ) push(req req) {
	q.queue <- req
}
