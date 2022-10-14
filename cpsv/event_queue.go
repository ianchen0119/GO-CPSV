package cpsv

func eventQInit() chan req {
	return make(chan req, 100)
}

func (ckpt *CkptOps) push(req req) {
	ckpt.q <- req
}
