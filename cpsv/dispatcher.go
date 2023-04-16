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
	"context"
	"fmt"
	"unsafe"
)

func Worker(id int, jobs <-chan req, ckpt *CkptOps) {
	for req := range jobs {
		ctx := context.Background()
		if ckpt.beforeUpdate != nil {
			if res := ckpt.beforeUpdate(context.WithValue(ctx, "req", req)); res != 0 {
				continue
			}
		}

		var status int
		cstr := C.CString(req.sectionId)
		cData := (*C.uchar)(unsafe.Pointer(&req.data[0]))
		defer C.free(unsafe.Pointer(cstr))

		if req.reqType == Fixed {
			status = int(C.ckpt_write(cstr, (*C.uchar)(cData), C.uint(req.offset), C.int(req.size)))
		} else {
			status = int(C.ckpt_non_fixed_write(cstr, (*C.uchar)(cData), C.uint(req.offset), C.int(req.size)))
		}

		if status == -1 && req.resend > 0 {
			req.resend--
			ckpt.push(req)
			if ckpt.ifError != nil {
				ckpt.ifError(context.WithValue(ctx, "req", req))
			}
		} else {
			if ckpt.afterUpdate != nil {
				ckpt.afterUpdate(context.WithValue(ctx, "req", req))
			}
		}
	}
}

func (ckpt *CkptOps) Dispatcher() {

	jobs := make(chan req, 100)
	fmt.Println("Dispatcher started, worker number: ", ckpt.workerNum)

	for w := 1; w <= ckpt.workerNum; w++ {
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
