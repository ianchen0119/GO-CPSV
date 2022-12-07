package main

import "C"
import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"unsafe"

	"github.com/ianchen0119/GO-CPSV/cpsv"
)

type Vertex struct {
	X int32
	Y int32
}

func main() {
	cpsv.Start("safCkpt=TEST1,safApp=safCkptService")
	v := &Vertex{X: 15, Y: 23}
	len := cpsv.GetSize(Vertex{})
	wbuf := C.GoBytes(unsafe.Pointer(v), C.int(len))

	cpsv.Store("data", wbuf, int(len), 0)

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		for s := range c {
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				fmt.Println("Program Exit...", s)
				cpsv.Destroy()
			case syscall.SIGUSR1:
				fmt.Println("usr1 signal", s)
			case syscall.SIGUSR2:
				fmt.Println("usr2 signal", s)
			default:
				fmt.Println("other signal", s)
			}
		}
	}()

	for {
	}
}
