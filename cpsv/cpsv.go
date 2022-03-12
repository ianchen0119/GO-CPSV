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

static int ckpt_read(void* buffer, unsigned int offset, int dataSize){
	return cpsv_sync_read((unsigned char*)buffer, offset, dataSize);
}

static int ckpt_write(void* data, unsigned int offset){
	return cpsv_sync_write((char*) data, offset);
}
*/
import "C"
import (
	"fmt"
	"os"
	"os/signal"
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
	}()

	go Dispatcher()
}

func Store(sectionId string, data []byte, offset int) {
	var newReq req
	newReq.sectionId = sectionId
	newReq.data = (*C.char)(unsafe.Pointer(&data[0]))
	newReq.offset = offset
	newReq.reqType = Sync
	q.push(newReq)
}

// load data from ckpt
func Load(buffer *[]byte, offset uint32, dataSize int) {
	C.ckpt_read(unsafe.Pointer(buffer),
		C.uint(offset), C.int(dataSize))
}
