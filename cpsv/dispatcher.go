package cpsv

/*
#cgo LDFLAGS: -L/usr/local/lib -lSaCkpt
#include "go-cpsv.h"

static int ckpt_write(char* sectionId, unsigned char* data, unsigned int offset, int dataSize){
	return cpsv_sync_write(sectionId, data, offset, dataSize);
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
			cstr := C.CString(req.sectionId)
			defer C.free(unsafe.Pointer(cstr))
			status := int(C.ckpt_write(cstr, (*C.uchar)(unsafe.Pointer(req.data)), C.uint(req.offset), C.int(req.size)))
			if status == -1 {
				q.push(req)
			}
		}
	}
}
