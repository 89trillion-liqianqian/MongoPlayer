package model

import (
	"MongoPlayer/protoFile/protoFile"
	"MongoPlayer/utils"
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type TimePorint struct {
	StartTime int64 `bson:"startTime"` //开始时间
	EndTime   int64 `bson:"endTime"`   //结束时间
}

type LogRecord struct {
	JobName string     `bson:"jobName"` //任务名
	Command string     `bson:"command"` //shell命令
	Err     string     `bson:"err"`     //脚本错误
	Content string     `bson:"content"` //脚本输出
	Tp      TimePorint //执行时间
}

//查询实体

type FindByJobName struct {
	JobName string `bson:"jobName"` //任务名
}

// 玩家的游戏数据

type PlayerModel struct {
	UniCode string // 客户端传递唯一识别码
	UID     string // 唯一UID
	Gold    int64  // 金币
	Diamond int64  // 钻石
}

// 增加玩家的金币钻石
func addPlayerData(uid string, gold, diamond int) (err error) {
	//update:=`{"$inc":{"gold":2}}`
	update := map[string]map[string]int{
		"$inc": {"gold": gold, "diamond": diamond},
	}
	uResult, err := collection.UpdateMany(context.TODO(), map[string]string{"uid": uid}, update)
	if err != nil {
		log.Println("--err addPlayerData", err, uResult.MatchedCount)
	}
	return
}

// 获取用户信息，根据 玩家uid
func getPlayerInfoByUid(uid string) (resData *PlayerModel, err error) {
	resData = new(PlayerModel)
	cursor, err := collection.Find(context.TODO(), map[string]string{"uid": uid})
	if err != nil {
		log.Println("--getPlayerInfo ", err)
		return
	}
	//延迟关闭游标
	defer func() {
		if err = cursor.Close(context.TODO()); err != nil {
			log.Fatal(err)
		}
	}()

	// 获取查询结果
	var playerData []PlayerModel
	err = cursor.All(context.TODO(), &playerData)
	if err != nil {
		return
	}
	for _, result := range playerData {
		if result.UID == uid {
			resData = &result
			return
		}
	}
	return
}

// 获取用户信息，根据 客户端传递唯一识别码
func getPlayerInfo(uniCode string) (resData *PlayerModel, err error) {
	resData = new(PlayerModel)
	cursor, err := collection.Find(context.TODO(), map[string]string{"unicode": uniCode})
	if err != nil {
		log.Println("--getPlayerInfo ", err)
		return
	}
	//延迟关闭游标
	defer func() {
		if err = cursor.Close(context.TODO()); err != nil {
			log.Fatal(err)
		}
	}()

	// 获取查询结果
	var playerData []PlayerModel
	err = cursor.All(context.TODO(), &playerData)
	if err != nil {
		return
	}
	for _, result := range playerData {
		if result.UniCode == uniCode {
			resData = &result
			return
		}
	}
	return
}

// 创建玩家
func createPlayer(uniCode string) (saveData *PlayerModel, err error) {
	// 生成 UID
	uid := utils.GetUID()
	saveData = new(PlayerModel)
	saveData.UID = uid
	saveData.UniCode = uniCode
	//_, err = collection.InsertOne(context.TODO(), &map[string]string{"name":"qqqq2"})
	_, err = collection.InsertOne(context.TODO(), saveData)
	if err != nil {
		return
	}
	return
}

// 客户端登陆注册

func GetLoginData(uniCode string) (isNew bool, playerData *PlayerModel, err error) {
	playerData, err = getPlayerInfo(uniCode)
	if playerData.UID == "" {
		isNew = true
		// 注册玩家
		playerData, err = createPlayer(uniCode)
	}
	return
}

// 测试 protobuf

func TestData() (reward protoFile.GeneralReward) {
	reward = protoFile.GeneralReward{}
	reward.Code = 200
	reward.Msg = "新增注册与登录接口"
	return
}

// 测试 连接mongo

func TestGetMongoDataModel() {
	var (
		err     error
		iResult *mongo.InsertOneResult
		id      primitive.ObjectID
		lr      *LogRecord
		cursor  *mongo.Cursor
	)
	//cond := FindByJobName{JobName: "job10"}
	//按照jobName字段进行过滤jobName="job10",翻页参数0-2
	if cursor, err = collection.Find(context.TODO(), map[string]string{"name": "qqqq2"}); err != nil {
		log.Println("--find", err)
		return
	}
	//延迟关闭游标
	defer func() {
		if err = cursor.Close(context.TODO()); err != nil {
			log.Fatal(err)
		}
	}()
	//这里的结果遍历可以使用另外一种更方便的方式：
	//var results []LogRecord
	var results []map[string]string
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Println("--All", err)
	}
	for _, result := range results {
		log.Println("--result=", result)
	}
	if iResult, err = collection.InsertOne(context.TODO(), &map[string]string{"name": "qqqq2"}); err != nil {
		log.Println("-2-InsertOne", err)
		return
	}
	id = iResult.InsertedID.(primitive.ObjectID)
	log.Println("--2自增ID", id.Hex())
	// 错误用法
	if iResult, err = collection.InsertOne(context.TODO(), lr); err != nil {
		log.Println("InsertOne", err)
		return
	}
	//_id:默认生成一个全局唯一ID
	id = iResult.InsertedID.(primitive.ObjectID)
	log.Println("自增ID", id.Hex())
}

//验证礼品码，保存用户数据

func SavePlayerGiftModel(uid, content string) (reward protoFile.GeneralReward, err error) {
	oldPlayer, _ := getPlayerInfoByUid(uid)
	var contentMap map[string]int
	json.Unmarshal([]byte(content), &contentMap)
	log.Println("--old", content, oldPlayer, contentMap)
	if len(contentMap) < 1 {
		return
	}
	balanceMap := make(map[uint32]uint64) // 变化前
	balanceMap[1] = uint64(oldPlayer.Gold)
	balanceMap[2] = uint64(oldPlayer.Diamond)

	changesMap := make(map[uint32]uint64) //变化量

	counterMap := make(map[uint32]uint64) //变化后
	counterMap[1] = uint64(oldPlayer.Gold)
	counterMap[2] = uint64(oldPlayer.Diamond)
	//添加玩家
	// 1- 金币，2-钻石
	gold := 0
	diamond := 0
	if _, ok := contentMap["1"]; ok {
		gold = contentMap["1"]
		changesMap[1] = uint64(gold)
		counterMap[1] += uint64(gold)
	}
	if _, ok := contentMap["2"]; ok {
		diamond = contentMap["2"]
		changesMap[2] = uint64(diamond)
		counterMap[2] += uint64(diamond)
	}
	err = addPlayerData(uid, gold, diamond)
	reward = protoFile.GeneralReward{}
	reward.Code = 200
	reward.Msg = "验证礼品码"
	reward.Balance = balanceMap
	reward.Changes = changesMap
	reward.Counter = counterMap
	return
}
