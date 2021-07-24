package ctrl

import (
	"MongoPlayer/internal/handler"
	"MongoPlayer/internal/model"
	"MongoPlayer/internal/myerr"
	"MongoPlayer/protoFile/protoFile"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

const (
	//1-指定用户一次性消耗，2-不指定用户限制兑换次数，3-不限用户不限次数兑换
	CodeTypeOne   = 1 //指定用户一次性消耗
	CodeTypeTwo   = 2 //不指定用户限制兑换次数
	CodeTypeThree = 3 //不限用户不限次数兑换
)

// ping
func PingFunc(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "ping",
	})
	return
}

// 创建礼品码
func CreateGift(c *gin.Context) {
	// 绑定参数数据
	var formData model.CreateGiftModels
	if err := c.ShouldBind(&formData); err != nil {
		log.Println("--CreateGift err", err)
		msg := "参数错误，管理后台-创建礼品码"
		myerr.ResponseErr(c, msg)
		return
	}
	// 参数校验
	if !(formData.CodeType == CodeTypeOne || formData.CodeType == CodeTypeTwo || formData.CodeType == CodeTypeThree) {
		msg := "参数错误,礼品码类型：1-指定用户一次性消耗，2-不指定用户限制兑换次数，3-不限用户不限次数兑换"
		log.Println("--CreateGift err", msg)
		myerr.ResponseErr(c, msg)
		return
	}
	// 创建礼品码
	code, err := handler.CreateGiftHandler(formData)
	if err != nil {
		msg := "创建礼品码失败，redis失败，请确认 "
		log.Println("--CreateGift err", msg)
		myerr.ResponseErr(c, msg)
		return
	}
	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "管理后台-创建礼品码 ok",
		"data": code,
	})
	return
}

// 查询礼品码
func GetGift(c *gin.Context) {
	code := c.Query("code")
	if len(code) != 8 {
		msg := "参数错误，不是礼品码"
		log.Println("--GetGift err", msg)
		myerr.ResponseErr(c, msg)
		return
	}
	// 	获取数据
	giftData, err := handler.GetGiftHandler(code)
	if err != nil {
		msg := "查询礼品码失败，redis失败，请确认 "
		log.Println("--GetGift err", msg)
		myerr.ResponseErr(c, msg)
		return
	}
	if len(giftData) < 1 {
		msg := "查询礼品码失败，礼品码不存在，请确认 "
		log.Println("--GetGift err", msg)
		myerr.ResponseErr(c, msg)
		return
	}
	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "查询礼品码成功",
		"data": giftData,
	})
	return
}

// 验证礼品码
func CheckCode(c *gin.Context) {
	code := c.PostForm("code")
	uid := c.PostForm("uid")
	resData := protoFile.GeneralReward{}
	if code == "" || uid == "" {
		log.Println("--CheckCode-code", code, uid)
		resData.Code = 1
		resData.Msg = "验证礼品码 。参数错误"
		c.ProtoBuf(http.StatusOK, &resData)
		return
	}
	// 从redis 读取数据 redis
	codeType, content, msg, err := handler.CheckCodeHandler(uid, code)
	if err != nil {
		resData.Code = 1
		resData.Msg = msg
		c.ProtoBuf(http.StatusOK, &resData)
		return
	}
	if len(content) < 1 {
		resData.Code = 1
		resData.Msg = msg
		c.ProtoBuf(http.StatusOK, &resData)
		return
	}

	// 保存用户数据
	resData, err = handler.SavePlayerGiftHandler(uid, content)
	if err != nil {
		// 保存mongo失败
		// 事物回滚
		err = handler.Rollback(code, uid, codeType)
	}
	// 返回
	c.ProtoBuf(http.StatusOK, &resData)
	return
}

// 新增注册与登录接口
func PlayerLogin(c *gin.Context) {
	uniCode := c.PostForm("uniCode")
	msg, playerData, _ := handler.PlayerLoginHandler(uniCode)
	c.JSON(200, gin.H{
		"code": 0,
		"msg":  msg,
		"data": playerData,
	})
	return
}
