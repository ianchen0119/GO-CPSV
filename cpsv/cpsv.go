package cpsv

/*
#cgo LDFLAGS: -L/usr/local/lib -lSaCkpt
#include "go-cpsv.h"
#include <stdint.h>

static void ckpt_init(char* ckptName){
	cpsv_ckpt_init(ckptName);
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
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"syscall"
	"unsafe"
	"errors"
)

func Start(ckptName string) {
	fmt.Println("Starting GO CPSV...")
	cStr := C.CString(ckptName)
	defer C.free(unsafe.Pointer(cStr))
	
	eventQInit()
	C.ckpt_init(cStr)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		fmt.Println("Signal:")
		fmt.Println(sig)
		C.ckpt_destroy()
		os.Exit(0)
	}()

	go Dispatcher()
}

func Destroy() {
	C.ckpt_destroy()
}

func Store(sectionId string, data []byte, size int, offset int) {
	var newReq req

	newReq.sectionId = sectionId
	newReq.data = data
	newReq.offset = 0
	newReq.reqType = Fixed
	newReq.size = size
	newReq.resend = 3
	q.push(newReq)
}

func NonFixedStore(sectionId string, data []byte, size int) {
	var newReq req

	newReq.sectionId = sectionId
	newReq.data = data
	newReq.offset = 4
	newReq.reqType = NonFixed
	newReq.size = size
	newReq.resend = 3
	q.push(newReq)
}

// load data from ckpt
func Load(sectionId string, offset uint32, dataSize int) ([]byte, error) {
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

func NonFixedLoad(sectionId string) ([]byte, error) {
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

func GetSize(i interface{}) int {
	size := reflect.TypeOf(i).Size()
	return int(size)
}

func GoBytes(unsafePtr unsafe.Pointer, length int) []byte {
	return C.GoBytes(unsafePtr, C.int(length))
}