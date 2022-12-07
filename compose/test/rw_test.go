package main

import "C"
import (
	"fmt"
	"sync"
	"time"

	"encoding/json"

	"github.com/ianchen0119/GO-CPSV/cpsv"
)

func testSyncMap() {
	var m sync.Map

	start := time.Now().UnixNano() / int64(time.Millisecond)

	for i := 0; i < 10000; i++ {
		data, _ := json.Marshal(i)
		m.Store(fmt.Sprintf("%d", i), data)
	}

	mid := time.Now().UnixNano() / int64(time.Millisecond)

	for i := 0; i < 10000; i++ {
		if _, ok := m.Load(i); ok {
			continue
		}
	}

	end := time.Now().UnixNano() / int64(time.Millisecond)

	fmt.Println("sync.Map", "W:", mid-start, "R:", end-mid, "Times:", 10000)
}

func testCPSVWrite(times int) {
	start := time.Now().UnixNano() / int64(time.Millisecond)

	for i := 0; i < times; i++ {
		data, _ := json.Marshal(i)
		cpsv.NonFixedStore(fmt.Sprintf("%d", i), data, len(data))
	}

	end := time.Now().UnixNano() / int64(time.Millisecond)

	fmt.Println("CPSV-W", "Start:", start, "End:", end, "Spent:", end-start, "Times:", times)
}

func testCPSVRead(times int) {
	start := time.Now().UnixNano() / int64(time.Millisecond)

	for i := 0; i < times; i++ {
		if _, err := cpsv.NonFixedLoad(fmt.Sprintf("%d", i)); err == nil {
			continue
		} else {
			fmt.Println(err)
		}
	}

	end := time.Now().UnixNano() / int64(time.Millisecond)

	fmt.Println("CPSV-R", "Start:", start, "End:", end, "Spent:", end-start, "Times:", times)
}

func main() {
	cpsv.Start("safCkpt=TEST2,safApp=safCkptService",
		cpsv.SetSectionNum(100000), cpsv.SetSectionSize(2000), cpsv.SetWorkerNum(10))

	testSyncMap()
	testCPSVWrite(100)
	time.Sleep(3 * time.Second)
	testCPSVRead(100)

	testCPSVWrite(1000)
	time.Sleep(3 * time.Second)
	testCPSVRead(1000)

	testCPSVWrite(10000)
	time.Sleep(3 * time.Second)
	testCPSVRead(10000)

	cpsv.Destroy()
}
