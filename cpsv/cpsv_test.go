package cpsv

import (
	"testing"
	"fmt"
	"unsafe"
	"time"
)

type Vertex struct {
	X int32
	Y int32
}

func TestStore_1(t *testing.T) {	
	Start("safCkpt=TEST1,safApp=safCkptService")
	v := &Vertex{X: 15, Y: 23}
	len := GetSize(Vertex{})
	wbuf := GoBytes(unsafe.Pointer(v), len)

	Store("d1", wbuf, int(len), 0)

	time.Sleep(3 * time.Second)

	var newY int32 = 20
	newYByte := GoBytes(unsafe.Pointer(&newY), len)
	Store("d1", newYByte, int(len), 4)

	time.Sleep(3 * time.Second)

	readData, err := Load("d1", 0, len)

	if err == nil {
		var bufV *Vertex = *(**Vertex)(unsafe.Pointer(&readData))
		if bufV.X != v.X && bufV.Y != newY {
			t.Error("exception: readData is not expected")
			fmt.Printf("X: %d, Y: %d", bufV.X, bufV.Y)
		}
	} else {
		t.Error("got errors:", err)
	}

	Destroy()
}

func TestStore_2(t *testing.T) {	
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

func TestLoad_1(t *testing.T) {
	Start("safCkpt=TEST1,safApp=safCkptService")

	len := GetSize(Vertex{})
	_, err := Load("d2", 0, len)

	if err == nil {
		t.Error("exception: readData shouldn't exist")
	}

	Destroy()
}

func TestNonFixedLoad_1(t *testing.T) {
	Start("safCkpt=TEST1,safApp=safCkptService")
	v := &Vertex{X: 4, Y: 23}
	len := GetSize(Vertex{})
	wbuf := GoBytes(unsafe.Pointer(v), len)

	Store("d1", wbuf, int(len), 0)

	time.Sleep(3 * time.Second)

	readData, err := NonFixedLoad("d1")

	if err == nil {
		var expectedY *int32 = *(**int32)(unsafe.Pointer(&readData))
		if *expectedY != v.Y {
			t.Error("exception: readData is not expected")
		}
	} else {
		t.Error("got errors:", err)
	}

	Destroy()
}

func TestNonFixedStore_1(t *testing.T) {
	Start("safCkpt=TEST1,safApp=safCkptService")
	v := &Vertex{X: 4, Y: 23}
	len := GetSize(Vertex{})
	wbuf := GoBytes(unsafe.Pointer(v), len)

	NonFixedStore("d1", wbuf, int(len))

	time.Sleep(3 * time.Second)

	readData, err := NonFixedLoad("d1")

	if err == nil {
		var bufV *Vertex = *(**Vertex)(unsafe.Pointer(&readData))
		if bufV.X != v.X || bufV.Y != v.Y {
			t.Error("exception: readData is not expected")
			fmt.Printf("X: %d, Y: %d", bufV.X, bufV.Y)
		}
	} else {
		t.Error("got errors:", err)
	}

	Destroy()
}