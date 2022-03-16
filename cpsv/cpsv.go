package cpsv

/*
#cgo LDFLAGS: -L/usr/local/lib -lSaCkpt
#include "go-cpsv.h"
#include <stdint.h>

static void ckpt_init(){
	cpsv_ckpt_init();
}
static void ckpt_destroy(){
	cpsv_ckpt_destroy();
}

static unsigned char* ckpt_read(char* sectionId, unsigned int offset, int dataSize){
	return cpsv_sync_read(sectionId, offset, dataSize);
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
)

func Start() {
	fmt.Println("Starting GO CPSV...")
	eventQInit()
	C.ckpt_init()

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
	newReq.offset = offset
	newReq.reqType = Sync
	newReq.size = size
	q.push(newReq)
}

// load data from ckpt
func Load(sectionId string, offset uint32, dataSize int) []byte {
	cstr := C.CString(sectionId)
	data := C.ckpt_read(cstr,
		C.uint(offset), C.int(dataSize))
	defer C.free(unsafe.Pointer(cstr))
	if data != nil && *(*C.uchar)(data) != 0 {
		defer C.free(unsafe.Pointer(data))
		return C.GoBytes(unsafe.Pointer(data), C.int(dataSize))
	}
	return make([]byte, dataSize)
}

func GetSize(i interface{}) int {
	size := reflect.TypeOf(i).Size()
	return int(size)
}
