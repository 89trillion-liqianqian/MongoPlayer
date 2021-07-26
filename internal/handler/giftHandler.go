package handler

import (
	"MongoPlayer/internal/model"
	"MongoPlayer/internal/service"
	"MongoPlayer/protoFile/protoFile"
	"log"
)

// 创建礼品码
func CreateGiftHandler(formData model.CreateGiftModels) (code string, err error) {
	codeType := formData.CodeType
	formData.CostCount = 0
	if codeType == model.CodeTypeOne {
		if formData.DrawCount != 1 {
			formData.DrawCount = 1
		}
	} else if codeType == model.CodeTypeThree {
		formData.DrawCount = 0
	}
	//code="SFDSHFUISD33"
	code = service.GetGiftCode()
	// 保存redis
	formData.Code = code
	err = model.SaveGiftRedis(formData)
	if codeType == model.CodeTypeTwo {
		// 限制次数的礼品码，增加set保存
		model.SaveGiftRedisType(code)
	}
	return
}

//  读取数据 redis
func GetGiftHandler(code string) (resData map[string]string, err error) {
	resData, err = model.GetGiftModel(code)
	return
}

// 领取礼品
func CheckCodeHandler(uid, code string) (codeType, content, msg string, err error) {
	codeType, content, msg, err = model.GetGiftReward(uid, code)
	return
}

//验证礼品码，保存用户数据
func SavePlayerGiftHandler(uid, content string) (reward protoFile.GeneralReward, err error) {
	// 保存用户数据
	reward, err = model.SavePlayerGiftModel(uid, content)
	return
}

//新增注册与登录接口
func PlayerLoginHandler(uniCode string) (msg string, playerData *model.PlayerModel, err error) {
	isNew, playerData, _ := model.GetLoginData(uniCode)
	log.Println("--uniCode", uniCode, isNew, playerData)
	msg = "新增注册与登录接口,老用户"
	if isNew {
		msg = "新增注册与登录接口,新用户"
	}
	return
}

// 事物回滚
func Rollback(code, uid, codeType string) (err error) {
	err = model.RollbackGiftCostHistoryRedis(code, uid, codeType)
	return
}
