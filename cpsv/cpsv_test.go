package cpsv

import (
	"testing"
	"unsafe"
	"time"
)

type Vertex struct {
	X int32
	Y int32
}

func TestStore(t *testing.T) {	
	Start("safCkpt=TEST1,safApp=safCkptService")
	v := &Vertex{X: 15, Y: 23}
	len := GetSize(Vertex{})
	wbuf := GoBytes(unsafe.Pointer(v), len)

	Store("d1", wbuf, int(len), 0)

	time.Sleep(3 * time.Second)

	readData, err := Load("d1", 0, len)

	if err == nil {
		var bufV *Vertex = *(**Vertex)(unsafe.Pointer(&readData))
		if bufV.X != v.X && bufV.Y != v.Y {
			t.Error("exception: readData is not expected")
		}
	} else {
		t.Error("got errors:", err)
	}

	Destroy()
}

func TestLoad(t *testing.T) {
	Start("safCkpt=TEST1,safApp=safCkptService")

	len := GetSize(Vertex{})
	_, err := Load("d2", 0, len)

	if err == nil {
		t.Error("exception: readData shouldn't exist")
	}

	Destroy()
}