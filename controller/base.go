package controller

import (
	"MongoPlayer/model"
	"MongoPlayer/protoFile/protoFile"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

/**
http get post 处理
*/

func Routers(r *gin.Engine) {
	r.GET("/ping", pingFunc)
	r.GET("/getGift", GetGift)
	r.POST("/createGift", CreateGift)
	r.POST("/checkCode", CheckCode)
	r.POST("/playerLogin", PlayerLogin)
	return
}

// 测试
func pingFunc(c *gin.Context) {
	// 获取mongo数据
	model.TestGetMongoDataModel()
	reward := model.TestData()
	c.ProtoBuf(http.StatusOK, &reward)
	log.Println("--re", reward)
	//c.JSON(200, gin.H{
	//	"message": "ping22",
	//})
	return
}

// 创建礼品码

func CreateGift(c *gin.Context) {
	var formData model.CreateGiftModels
	if err := c.ShouldBind(&formData); err != nil {
		log.Println("--CreateGift err", err)
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "参数错误，管理后台-创建礼品码",
			"data": "",
		})
		return
	}
	// 保存redis
	code, err := model.CreateGiftModel(formData)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "创建礼品码失败，redis失败，请确认 ",
			"data": "",
		})
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
	log.Println("--GetGift-code", code)
	// 从redis 读取数据 redis
	giftData, err := model.GetGiftModel(code)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "查询礼品码失败，redis失败，请确认 ",
			"data": "",
		})
		return
	}
	if len(giftData) < 1 {
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "查询礼品码失败，礼品码不存在，请确认 ",
			"data": "",
		})
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
		//c.JSON(200, gin.H{
		//	"code": 1,
		//	"msg":  "验证礼品码 。参数错误",
		//	"data": "",
		//})
		return
	}
	log.Println("--CheckCode-code", code, uid)
	// 从redis 读取数据 redis
	content, msg, err := model.GetGiftReward(uid, code)
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
	resData, _ = model.SavePlayerGiftModel(uid, content)
	// 返回
	c.ProtoBuf(http.StatusOK, &resData)
	return
}

// 新增注册与登录接口

func PlayerLogin(c *gin.Context) {
	uniCode := c.PostForm("uniCode")
	isNew, playerData, _ := model.GetLoginData(uniCode)
	log.Println("--uniCode", uniCode, isNew, playerData)
	msg := "新增注册与登录接口,老用户"
	if isNew {
		msg = "新增注册与登录接口,新用户"
	}
	c.JSON(200, gin.H{
		"code": 0,
		"msg":  msg,
		"data": playerData,
	})
	return
}
