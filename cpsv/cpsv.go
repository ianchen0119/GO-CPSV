package cpsv

/*
#cgo LDFLAGS: -L/usr/local/lib -lSaCkpt
#include "go-cpsv.h"

static void ckpt_init(){
	cpsv_ckpt_init();
}
static void ckpt_destroy(){
	cpsv_ckpt_destroy();
}

static void read(void* buffer, uint32_t offset, int dataSize){
	cpsv_sync_read((unsigned char*)buffer, offset, dataSize);
}

static void write(voidr* data, uint32_t offset){
	cpsv_sync_write((char*) data, offset);
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
	fmt.Println("Starting Go CPSV...")
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
	newReq.data = C.CBytes(data)
	newReq.offset = offset
	newReq.reqType = Sync
	q.push(newReq)
}

// load data from ckpt
func Load(buffer []byte, offset uint32, dataSize int) {
	C.read(unsafe.Pointer(&buffer[0])),
		C.uint(offset), C.int(dataSize))
}
