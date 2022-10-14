package cpsv

import (
	"sync"
	"unsafe"
)

type CheckPoint struct {
	ops *CkptOps
	mu  sync.Mutex
}

var globalCkpt CheckPoint

func StartWithSectionConfig(ckptName string, sections int, secitonSize int) {
	globalCkpt.mu.Lock()
	defer globalCkpt.mu.Unlock()
	globalCkpt.ops = startWithSectionConfig(ckptName, sections, secitonSize)
}

func Start(ckptName string) {
	globalCkpt.mu.Lock()
	defer globalCkpt.mu.Unlock()
	globalCkpt.ops = start(ckptName)
}

func Destroy() {
	globalCkpt.mu.Lock()
	defer globalCkpt.mu.Unlock()
	globalCkpt.ops.destroy()
}

func Store(sectionId string, data []byte, size int, offset int) {
	globalCkpt.ops.store(sectionId, data, size, offset)
}

func NonFixedStore(sectionId string, data []byte, size int) {
	globalCkpt.ops.nonFixedStore(sectionId, data, size)
}

func Load(sectionId string, offset uint32, dataSize int) ([]byte, error) {
	return globalCkpt.ops.load(sectionId, offset, dataSize)
}

func NonFixedLoad(sectionId string) ([]byte, error) {
	return globalCkpt.ops.nonFixedLoad(sectionId)
}

func GetSize(i interface{}) int {
	return globalCkpt.ops.getSize(i)
}

func GoBytes(unsafePtr unsafe.Pointer, length int) []byte {
	return globalCkpt.ops.goBytes(unsafePtr, length)
}
