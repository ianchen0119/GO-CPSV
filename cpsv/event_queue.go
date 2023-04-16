package cpsv

func eventQInit(num int) chan req {
	return make(chan req, num)
}

func (ckpt *CkptOps) push(req req) {
	ckpt.q <- req
}
