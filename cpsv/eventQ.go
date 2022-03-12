package cpsv

var q EventQ

func eventQInit() {
	q.queue = make(chan req, 100)
}

func (q *EventQ) pull(req *req) {
	*req = <-q.queue
}

func (q *EventQ) push(req req) {
	q.queue <- req
}
