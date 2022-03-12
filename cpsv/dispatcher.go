package cpsv

/*
#cgo LDFLAGS: -L/usr/local/lib -lSaCkpt
#include "go-cpsv.h"

static int ckpt_write(void* data, unsigned int offset){
	return cpsv_sync_write((char*) data, offset);
}
*/
import "C"
import "unsafe"

func Dispatcher() {
	for true {
		var req *req
		q.pull(req)
		status := int(C.ckpt_write(unsafe.Pointer(&req.data), C.uint(req.offset)))
		if status == -1 {
			q.push(*req)
		}
	}
}
