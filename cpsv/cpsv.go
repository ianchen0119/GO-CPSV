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

static int ckpt_read(char* sectionId, void* buffer, unsigned int offset, int dataSize){
	return cpsv_sync_read(sectionId, (unsigned char*)buffer, offset, dataSize);
}
*/
import "C"
import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"unsafe"
	"reflect"
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
	newReq.data = (*C.char)(unsafe.Pointer(&data[0]))
	newReq.offset = offset
	newReq.reqType = Sync
	newReq.size = size
	q.push(newReq)
}

// load data from ckpt
func Load(sectionId string, buffer *[]byte, offset uint32, dataSize int) int {
	cstr := C.CString(sectionId)
	defer C.free(unsafe.Pointer(cstr))
	return int(C.ckpt_read(cstr, unsafe.Pointer(buffer),
		C.uint(offset), C.int(dataSize)))
}

func GetSize(i interface{}) int {
	size := reflect.TypeOf(i).Size()
	return int(size)
}
