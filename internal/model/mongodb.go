package model

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var mgoCli *mongo.Client
var (
	client     = GetMgoCli()
	db         *mongo.Database
	collection *mongo.Collection
)

func initEngine() {
	var err error
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// 连接到MongoDB
	mgoCli, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	// 检查连接
	err = mgoCli.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
}
func GetMgoCli() *mongo.Client {
	if mgoCli == nil {
		initEngine()
	}
	return mgoCli
}

// 初始化
func InitMongo() {
	//2.选择数据库 my_db
	db = client.Database("gift")

	//选择表 my_collection
	collection = db.Collection("player")
	collection = collection
}

func GetMongoCol() *mongo.Collection {

	return collection
}
