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
	"fmt"
	"unsafe"
)

func (ckpt *CkptOps) Dispatcher() {
	for {
		select {
		case <-ckpt.stopCh:
			// notify the main thread that the dispatcher is stopped
			defer close(ckpt.notifyCh)
			return
		case req, ok := <-ckpt.q:
			if ok {
				var status int
				fmt.Println("handle event from eventQ")
				cstr := C.CString(req.sectionId)
				cData := C.CBytes(req.data)
				defer C.free(unsafe.Pointer(cstr))
				defer C.free(cData)

				if req.reqType == Fixed {
					fmt.Println("Fixed")
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
	}
}
