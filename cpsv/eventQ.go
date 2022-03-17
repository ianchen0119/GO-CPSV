package cpsv

var q eventQ

func eventQInit() {
	q.queue = make(chan req, 100)
}

func (q *eventQ) push(req req) {
	q.queue <- req
}
