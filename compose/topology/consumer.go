package main

import "C"
import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
	"unsafe"

	"github.com/ianchen0119/GO-CPSV/cpsv"
)

type Vertex struct {
	X int32
	Y int32
}

func main() {
	cpsv.Start("safCkpt=TEST1,safApp=safCkptService")
	len := cpsv.GetSize(Vertex{})

	readData, err := cpsv.Load("data", 0, len)
	if err == nil {
		var bufV *Vertex = *(**Vertex)(unsafe.Pointer(&readData))
		fmt.Printf("X: %d, Y:%d\n", bufV.X, bufV.Y)
	} else {
		fmt.Println(err)
		time.Sleep(5 * time.Second)
	}

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
