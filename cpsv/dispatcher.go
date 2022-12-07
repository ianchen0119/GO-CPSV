package cpsv

/*
#cgo LDFLAGS: -L/usr/local/lib -lSaCkpt
#include "go-cpsv.h"

static int ckpt_write(char* sectionId, unsigned char* data, unsigned int offset, int dataSize){
	return cpsv_sync_write(sectionId, data, offset, dataSize, 1);
}

static int ckpt_non_fixed_write(char* sectionId, unsigned char* data, unsigned int offset, int dataSize){
	return cpsv_sync_write(sectionId, data, offset, dataSize, 0);
}
*/
import "C"
import (
	"unsafe"
)

func Worker(id int, jobs <-chan req, ckpt *CkptOps) {
	for req := range jobs {
		var status int
		cstr := C.CString(req.sectionId)
		cData := C.CBytes(req.data)
		defer C.free(unsafe.Pointer(cstr))
		defer C.free(cData)

		if req.reqType == Fixed {
			status = int(C.ckpt_write(cstr, (*C.uchar)(cData), C.uint(req.offset), C.int(req.size)))
		} else {
			status = int(C.ckpt_non_fixed_write(cstr, (*C.uchar)(cData), C.uint(req.offset), C.int(req.size)))
		}

		if status == -1 && req.resend > 0 {
			req.resend--
			ckpt.push(req)
		}
	}
}

func (ckpt *CkptOps) Dispatcher() {

	jobs := make(chan req, 100)
	for w := 1; w <= 10; w++ {
		go Worker(w, jobs, ckpt)
	}

	for {
		select {
		case <-ckpt.stopCh:
			// notify the main thread that the dispatcher is stopped
			defer close(ckpt.notifyCh)
			return
		case req, ok := <-ckpt.q:
			if ok {
				jobs <- req
			}
		}
	}
}
