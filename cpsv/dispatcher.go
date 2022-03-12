package cpsv

/*
#cgo LDFLAGS: -L/usr/local/lib -lSaCkpt
#include "go-cpsv.h"

static int ckpt_write(char* data, unsigned int offset){
	return cpsv_sync_write(data, offset);
}
*/
import "C"
import (
	"fmt"
	"unsafe"
)

func Dispatcher() {
	for {
		req, ok := <-q.queue
		if ok {
			fmt.Println("handle event from eventQ")
			status := int(C.ckpt_write((*C.char)(unsafe.Pointer(&req.data)), C.uint(req.offset)))
			if status == -1 {
				q.push(req)
			}
		}
	}
}
