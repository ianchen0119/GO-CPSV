package cpsv

/*
#cgo LDFLAGS: -L/usr/local/lib -lSaCkpt
#include "go-cpsv.h"
#include <stdint.h>

static void ckpt_init_with_section(char* newName, int sections, int sectionSize){
	cpsv_ckpt_init_with_section_number(newName, sections, sectionSize);
}

static void ckpt_destroy(){
	cpsv_ckpt_destroy();
}

static unsigned char* ckpt_read(char* sectionId, unsigned int offset, int dataSize){
	return cpsv_sync_read(sectionId, offset, dataSize, 1, (void*) 0);
}

static unsigned char* ckpt_non_fixed_read(char* sectionId, int* dataSizePtr){
	return cpsv_sync_read(sectionId, 0, 4, 0, dataSizePtr);
}
*/
import "C"

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"syscall"
	"time"
	"unsafe"
)

var _ storageAPI = (*CkptOps)(nil)

type CkptOps struct {
	startTime   time.Time
	q           chan req
	sectionNum  int
	suctionSize int
	workerNum   int
	resendMax   int
	stopCh      chan struct{}
	notifyCh    chan struct{}
}

type Option func(*CkptOps)

func SetResendMax(max int) Option {
	return func(ckpt *CkptOps) {
		ckpt.resendMax = max
	}
}

func SetSectionNum(num int) Option {
	return func(ckpt *CkptOps) {
		ckpt.sectionNum = num
	}
}

func SetSectionSize(size int) Option {
	return func(ckpt *CkptOps) {
		ckpt.suctionSize = size
	}
}

func SetWorkerNum(num int) Option {
	return func(ckpt *CkptOps) {
		ckpt.workerNum = num
	}
}

func start(ckptName string, ops ...func(*CkptOps)) *CkptOps {
	fmt.Println("Starting GO CPSV...")
	cStr := C.CString(ckptName)
	defer C.free(unsafe.Pointer(cStr))

	cpsv := &CkptOps{
		startTime:   time.Now(),
		q:           eventQInit(),
		sectionNum:  100000,
		suctionSize: 20000,
		stopCh:      make(chan struct{}),
		notifyCh:    make(chan struct{}),
		workerNum:   3,
	}

	for _, op := range ops {
		op(cpsv)
	}

	C.ckpt_init_with_section(cStr, C.int(cpsv.sectionNum), C.int(cpsv.suctionSize))

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		fmt.Println("Signal:")
		fmt.Println(sig)
		C.ckpt_destroy()
		os.Exit(0)
	}()

	go cpsv.Dispatcher()

	return cpsv
}

func (ckpt *CkptOps) destroy() {
	// notify dispatcher to stop
	close(ckpt.stopCh)
	// wait for dispatcher to stop
	<-ckpt.notifyCh
	C.ckpt_destroy()
}

func (ckpt *CkptOps) store(sectionId string, data []byte, size int, offset int) {
	var newReq req

	newReq.sectionId = sectionId
	newReq.data = data
	newReq.offset = offset
	newReq.reqType = Fixed
	newReq.size = size
	newReq.resend = ckpt.resendMax
	ckpt.push(newReq)
}

func (ckpt *CkptOps) nonFixedStore(sectionId string, data []byte, size int) {
	var newReq req

	newReq.sectionId = sectionId
	newReq.data = data
	newReq.offset = 4
	newReq.reqType = NonFixed
	newReq.size = size
	newReq.resend = ckpt.resendMax
	ckpt.push(newReq)
}

// load data from ckpt
func (ckpt *CkptOps) load(sectionId string, offset uint32, dataSize int) ([]byte, error) {
	defer func() {
		r := recover()
		if r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()
	cStr := C.CString(sectionId)
	data := C.ckpt_read(cStr,
		C.uint(offset), C.int(dataSize))
	defer C.free(unsafe.Pointer(cStr))
	if data != nil {
		defer C.free(unsafe.Pointer(data))
		return C.GoBytes(unsafe.Pointer(data), C.int(dataSize)), nil
	}
	return make([]byte, dataSize), errors.New("No data found")
}

func (ckpt *CkptOps) nonFixedLoad(sectionId string) ([]byte, error) {
	defer func() {
		r := recover()
		if r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()
	dataSize := 0
	dataSizePtr := (*C.int)(unsafe.Pointer(&dataSize))
	cStr := C.CString(sectionId)
	data := C.ckpt_non_fixed_read(cStr, dataSizePtr)
	defer C.free(unsafe.Pointer(cStr))
	if data != nil {
		defer C.free(unsafe.Pointer(data))
		return C.GoBytes(unsafe.Pointer(data), C.int(dataSize)), nil
	}
	return make([]byte, dataSize), errors.New("No data found")
}

func (ckpt *CkptOps) getSize(i interface{}) int {
	size := reflect.TypeOf(i).Size()
	return int(size)
}

func (ckpt *CkptOps) goBytes(unsafePtr unsafe.Pointer, length int) []byte {
	return C.GoBytes(unsafePtr, C.int(length))
}
