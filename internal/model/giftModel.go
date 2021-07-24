package model

import (
	"MongoPlayer/utils"
	"encoding/json"
	redigo "github.com/gomodule/redigo/redis"
	"log"
	"strconv"
	"time"
)

const GiftType = "_type"

const (
	//1-指定用户一次性消耗，2-不指定用户限制兑换次数，3-不限用户不限次数兑换
	CodeTypeOne   = 1 //指定用户一次性消耗
	CodeTypeTwo   = 2 //不指定用户限制兑换次数
	CodeTypeThree = 3 //不限用户不限次数兑换

	CodeTypeOneStr   = "1" //指定用户一次性消耗
	CodeTypeTwoStr   = "2" //不指定用户限制兑换次数
	CodeTypeThreeStr = "3" //不限用户不限次数兑换
)

// 礼品码信息
type CreateGiftModels struct {
	Code       string `form:"code" binding:""`
	CodeType   int    `form:"codeType" binding:"required"`
	DrawCount  int    `form:"drawCount" binding:"required"`
	Des        string `form:"des" binding:"required"`
	ValidTime  int64  `form:"validTime" binding:"required"`
	Content    string `form:"content" binding:"required"`
	CreateUser string `form:"createUser" binding:"required"`
	CostCount  int    `form:"costCount" binding:""`
	UserId     int    `form:"userId" binding:"required"`
}

// 事物回滚，删除礼品码的领取用户历史
func RollbackGiftCostHistoryRedis(code, uid, codeTypeStr string) (err error) {
	conn := RedisPool.Get()
	defer conn.Close()
	codeHistory := code + "_history"
	codeType := code + "_type"
	res, err := redigo.Int64(conn.Do("HDEl", codeHistory, uid))
	if err != nil {
		log.Println("--RollbackGiftCostHistoryRedis", res, err)
	}
	res, err = redigo.Int64(conn.Do("HINCRBY", code, "CostCount", -1))
	if err != nil {
		log.Println("--RollbackGiftCostHistoryRedis", res, err)
	}
	// 类型2，限制次数  decr
	if codeTypeStr == CodeTypeTwoStr {
		res, err = redigo.Int64(conn.Do("DECR", codeType))
		if err != nil {
			log.Println("--RollbackGiftCostHistoryRedis -decr ", codeType, res, err)
			return
		}
	}
	return
}

// 获取玩家礼品领取记录
func getPlayerGiftHistory(code, uid string) (historyData string, err error) {
	conn := RedisPool.Get()
	defer conn.Close()
	code += "_history"
	historyData, err = redigo.String(conn.Do("HGET", code, uid))
	if err != nil {
		log.Println("--getPlayerGiftHistory", historyData, err)
	}
	return
}

// 领取限制次数的礼品码
func getGiftTypeTwo(code string, drawCount int) (isOk bool, err error) {
	conn := RedisPool.Get()
	defer conn.Close()
	codeType := code + "_type"
	count := 1
RETRY:
	count += 1
	lock, err := Lock()
	if !lock {
		// 取消设置
		if count > 100 {
			return
		}
		// 重试
		goto RETRY
	}
	// 获取领取次数
	costCount, err := redigo.Int(conn.Do("GET", codeType))
	if err != nil {
		Unlock("lock_value")
		log.Println("--getGiftTypeTwo", costCount, err)
		return
	}
	if costCount >= drawCount {
		Unlock("lock_value")
		log.Println("--getGiftTypeTwo,领取失败，次数已使用完", costCount, err)
		return
	}
	res, err := redigo.Int64(conn.Do("INCR", codeType))
	if err != nil {
		Unlock("lock_value")
		log.Println("--getGiftTypeTwo -add ", codeType, res, err)
		return
	}
	res, err = redigo.Int64(conn.Do("HINCRBY", code, "CostCount", 1))
	if err != nil {
		Unlock("lock_value")
		log.Println("--saveGiftCostCountRedis", res, err)
	}
	isOk = true
	Unlock("lock_value")
	return
}

// 保存礼品码的领取用户历史
func saveGiftCostHistoryRedis(code, uid string) (err error) {
	conn := RedisPool.Get()
	defer conn.Close()
	code += "_history"
	nowTime := time.Now().Unix()
	res, err := redigo.Int64(conn.Do("HSET", code, uid, nowTime))
	if err != nil {
		log.Println("--saveGiftCostHistoryRedis", res, err)
	}
	return
}

// 保存礼品码的领取次数数据
func saveGiftCostCountRedis(code string) (err error) {
	conn := RedisPool.Get()
	defer conn.Close()
	log.Println("----test", code, 1)
	res, err := redigo.Int64(conn.Do("HINCRBY", code, "CostCount", 1))
	if err != nil {
		log.Println("--saveGiftCostCountRedis", res, err)
	}
	return
}

// 获取礼品数据
func getGiftRedis(code string) (resMap map[string]string, err error) {
	conn := RedisPool.Get()
	defer conn.Close()

	resMap, err = redigo.StringMap(conn.Do("HGETAll", code))
	if err != nil {
		log.Println("--getGiftRedis", resMap, err)
	}
	return
}

// 保存礼品数据
func SaveGiftRedis(formData CreateGiftModels) (err error) {
	conn := RedisPool.Get()
	defer conn.Close()

	res, err := redigo.String(conn.Do("HMSET", redigo.Args{formData.Code}.AddFlat(formData)...))
	if err != nil {
		log.Println("--saveGiftRedis", res, err)
	}
	return
}

// 保存礼品数据
func SaveGiftRedisType(code string) (err error) {
	conn := RedisPool.Get()
	defer conn.Close()

	// 限制兑换次数
	code += GiftType
	res, err := redigo.String(conn.Do("set", code, 0))
	if err != nil {
		log.Println("--saveGiftRedisType", res, err)
	}
	return
}

// 查询礼品码
func GetGiftModel(code string) (resData map[string]string, err error) {
	resData, err = getGiftRedis(code)
	if len(resData) > 0 {
		// 获取历史记录
		historyData, _ := getGiftRedis(code + "_history")
		b, _ := json.Marshal(historyData)
		resData["historyData"] = string(b)
	}
	return
}

// 领取礼品
func GetGiftReward(uid, code string) (codeType, content, msg string, err error) {
	resData, err := getGiftRedis(code)
	if len(resData) < 1 {
		return
	}
	//  判有效性
	if _, ok := resData["validTime"]; ok {
		validTime := resData["validTime"]
		if !utils.CheckTime(validTime) {
			// 过期，领取失败，
			msg = "过期，领取失败"
			return
		}
	}
	//uid:=0
	codeType = CodeTypeOneStr // 1-指定用户一次性消耗，2-不指定用户限制兑换次数，3-不限用户不限次数兑换
	if _, ok := resData["CodeType"]; ok {
		codeType = resData["CodeType"]
	}
	costCount := resData["CostCount"]
	costCountInt, _ := strconv.Atoi(costCount)
	if codeType == CodeTypeOneStr {
		uidStr := resData["UserId"]
		if uid == uidStr && costCountInt < 1 {
			// 返礼品
			content = resData["Content"]
			//  消耗、保存
			err = saveGiftCostCountRedis(code)
			err = saveGiftCostHistoryRedis(code, uid)
			return
		} else {
			// 领取失败，
			msg = "失败，已领取"
			return
		}
	} else if codeType == CodeTypeTwoStr {
		historyData, _ := getPlayerGiftHistory(code, uid)
		if len(historyData) > 0 {
			msg = "玩家已经领取过该礼品码"
			return
		}
		//  限制次数
		drawCount := resData["DrawCount"]
		drawCountInt, _ := strconv.Atoi(drawCount)
		if costCountInt >= drawCountInt {
			msg = "礼品码，领取次数已完"
			return
		}
		// 领取礼品
		isOk, _ := getGiftTypeTwo(code, drawCountInt)
		if isOk {
			// 返礼品
			content = resData["Content"]
			// 保存历史
			err = saveGiftCostHistoryRedis(code, uid)
			return
		}
		msg = "领取失败,次数已使用完"
		return
	} else if codeType == CodeTypeThreeStr {
		historyData, _ := getPlayerGiftHistory(code, uid)
		if len(historyData) > 0 {
			msg = "玩家已经领取过该礼品码"
			return
		}
		// 返礼品
		content = resData["Content"]
		//  消耗、更新次数，保存历史
		err = saveGiftCostCountRedis(code)
		err = saveGiftCostHistoryRedis(code, uid)
		if err != nil {
			msg = "领取失败，请重试"
		}
		return
	}
	return
}
