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
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection
var client *mongo.Client

func initMongo() {
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://db:27017"))
	if err != nil {
		fmt.Printf("failed to connect to mongo: %v", err)
	}

	collection = client.Database("cpsv-test").Collection("test")
}

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

func testMongoWrite(times int) {
	start := time.Now().UnixNano() / int64(time.Millisecond)

	for i := 0; i < times; i++ {
		data, _ := json.Marshal(i)
		if _, err := collection.InsertOne(context.Background(), bson.D{{"value", data}}); err != nil {
			log.Fatal(err)
		}
	}

	end := time.Now().UnixNano() / int64(time.Millisecond)

	fmt.Println("Mongo-W", "Start:", start, "End:", end, "Spent:", end-start, "Times:", times)
}

func testMongoRead(times int) {
	start := time.Now().UnixNano() / int64(time.Millisecond)

	for i := 0; i < times; i++ {
		_ = collection.FindOne(context.Background(), bson.M{})
	}

	end := time.Now().UnixNano() / int64(time.Millisecond)

	fmt.Println("Mongo-R", "Start:", start, "End:", end, "Spent:", end-start, "Times:", times)
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

func testCPSVFRead(times int) {
	start := time.Now().UnixNano() / int64(time.Millisecond)

	for i := 0; i < times; i++ {
		if _, err := cpsv.Load(fmt.Sprintf("%d", i), 0, 4); err == nil {
			continue
		} else {
			fmt.Println(err)
		}
	}

	end := time.Now().UnixNano() / int64(time.Millisecond)

	fmt.Println("CPSV-Fixed-R", "Start:", start, "End:", end, "Spent:", end-start, "Times:", times)
}

func testCPSVFWrite(times int) {
	start := time.Now().UnixNano() / int64(time.Millisecond)

	for i := 0; i < times; i++ {
		data, _ := json.Marshal(i)
		cpsv.Store(fmt.Sprintf("%d", i), data, len(data), 0)
	}

	end := time.Now().UnixNano() / int64(time.Millisecond)

	fmt.Println("CPSV-Fixed-W", "Start:", start, "End:", end, "Spent:", end-start, "Times:", times)
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	// cpsv.Start("safCkpt=TEST2,safApp=safCkptService",
	// 	cpsv.SetSectionNum(100000), cpsv.SetSectionSize(2000),
	// 	cpsv.SetWorkerNum(10), cpsv.SetLifeCycleHooks(beforeUpdate, success, fail))

	cpsv.Start("safCkpt=TEST2,safApp=safCkptService",
		cpsv.SetSectionNum(100000), cpsv.SetSectionSize(2000),
		cpsv.SetWorkerNum(20))

	time.Sleep(3 * time.Second)

	initMongo()

	// testSyncMap()
	// testCPSVWrite(100)
	testCPSVFWrite(100)
	testMongoWrite(100)
	time.Sleep(3 * time.Second)
	// testCPSVRead(100)
	testCPSVFRead(100)
	testMongoRead(100)

	// testCPSVWrite(1000)
	testCPSVFWrite(1000)
	testMongoWrite(1000)
	time.Sleep(3 * time.Second)
	// testCPSVRead(1000)
	testCPSVFRead(1000)
	testMongoRead(1000)

	// testCPSVWrite(8000)
	testCPSVFWrite(8000)
	testMongoWrite(8000)
	time.Sleep(3 * time.Second)
	// testCPSVRead(8000)
	testCPSVFRead(8000)
	testMongoRead(8000)

	// testCPSVWrite(10000)
	testCPSVFWrite(10000)
	testMongoWrite(10000)
	time.Sleep(3 * time.Second)
	// testCPSVRead(10000)
	testCPSVFRead(10000)
	testMongoRead(10000)

	// testCPSVWrite(15000)
	testCPSVFWrite(15000)
	testMongoWrite(15000)
	time.Sleep(3 * time.Second)
	// testCPSVRead(15000)
	testCPSVFRead(15000)
	testMongoRead(15000)

	// testCPSVWrite(20000)
	testCPSVWrite(20000)
	testMongoWrite(20000)
	time.Sleep(3 * time.Second)
	// testCPSVRead(20000)
	testCPSVFRead(20000)
	testMongoRead(20000)

	cpsv.Destroy()
}
