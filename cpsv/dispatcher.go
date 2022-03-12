package cpsv

/*
#cgo LDFLAGS: -L/usr/local/lib -lSaCkpt
#include "go-cpsv.h"

static void write(void* data, unsigned int offset){
	cpsv_sync_write((char*) data, offset);
}
*/

import "C"

func Dispatcher() {
	for true {
	}
}
