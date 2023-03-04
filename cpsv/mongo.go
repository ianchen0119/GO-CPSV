package cpsv

import "C"

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"time"
	"unsafe"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// implement mongo client, which implements storageAPI

type MongoOps struct {
	collection *mongo.Collection
	client     *mongo.Client
}

type ByteArray struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	Data []byte             `bson:"data,omitempty"`
}

var _ storageAPI = (*MongoOps)(nil)

func MongoStart(collName string) *MongoOps {
	var mongoEntity MongoOps
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mongoEntity.client, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://db:27017"))
	if err != nil {
		fmt.Printf("failed to connect to mongo: %v", err)
	}

	mongoEntity.collection = mongoEntity.client.Database("cpsv-test").Collection(collName)
	return &mongoEntity
}

func (mongoEntity *MongoOps) destroy() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mongoEntity.client.Disconnect(ctx)
}

// store byte array to mongo
func (mongoEntity *MongoOps) store(sectionId string, data []byte, size int, offset int) {
	// delete old data
	if _, err := mongoEntity.collection.DeleteOne(context.Background(), bson.M{"sectionId": sectionId}); err != nil {
		log.Fatal(err)
	}
	// insert new data
	target := &ByteArray{
		Data: data,
	}
	if _, err := mongoEntity.collection.InsertOne(context.Background(), bson.M{"sectionId": sectionId, "value": target}); err != nil {
		log.Fatal(err)
	}
}

func (mongoEntity *MongoOps) nonFixedStore(sectionId string, data []byte, size int) {
	mongoEntity.store(sectionId, data, size, 0)
}

func (mongoEntity *MongoOps) load(sectionId string, offset uint32, dataSize int) ([]byte, error) {
	var result map[string]interface{}
	mongoEntity.collection.FindOne(context.Background(), bson.M{"sectionId": sectionId}).Decode(&result)
	return result["value"].(map[string]interface{})["data"].(primitive.Binary).Data, nil
}

func (mongoEntity *MongoOps) nonFixedLoad(sectionId string) ([]byte, error) {
	return mongoEntity.load(sectionId, 0, 0)
}

func (mongoEntity *MongoOps) getSize(i interface{}) int {
	size := reflect.TypeOf(i).Size()
	return int(size)
}

func (mongoEntity *MongoOps) goBytes(unsafePtr unsafe.Pointer, length int) []byte {
	return C.GoBytes(unsafePtr, C.int(length))
}
