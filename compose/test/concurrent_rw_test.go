package main

import "C"
import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"encoding/json"

	"github.com/ianchen0119/GO-CPSV/cpsv"
)

func testCPSVRW(times int) {
	wg := sync.WaitGroup{}

	start := time.Now().UnixNano() / int64(time.Millisecond)
	wg.Add(2)
	go func() {
		for i := 0; i < times; i++ {
			data, _ := json.Marshal(i)
			cpsv.Store(fmt.Sprintf("%d", i), data, len(data), 0)
		}
		wg.Done()
	}()
	go func() {
		for i := 0; i < times; i++ {
			if _, err := cpsv.Load(fmt.Sprintf("%d", i), 0, 4); err == nil {
				continue
			} else {
				fmt.Println(err)
			}
		}
		wg.Done()
	}()

	wg.Wait()

	end := time.Now().UnixNano() / int64(time.Millisecond)

	fmt.Println("CPSV-RW", "Start:", start, "End:", end, "Spent:", end-start, "Times:", times)
}

func success(ctx context.Context) {
	// get the value from context
	res, err := cpsv.GetResult(ctx)
	if err != nil {
		log.Fatalln("failed to get result", err)
	} else {
		log.Println("success", res.SecId, res.Data)
	}
}

func fail(ctx context.Context) {
	log.Println("fail")
}

func beforeUpdate(ctx context.Context) {
	log.Println("beforeUpdate")
}

func main() {
	cpsv.Start("safCkpt=TEST2,safApp=safCkptService",
		cpsv.SetSectionNum(100000), cpsv.SetSectionSize(2000),
		cpsv.SetWorkerNum(20))

	time.Sleep(3 * time.Second)

	testCPSVRW(3000)
	testCPSVRW(8000)
	testCPSVRW(15000)

	cpsv.Destroy()
}
