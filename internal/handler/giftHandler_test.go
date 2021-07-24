package handler

import (
	"MongoPlayer/internal/model"
	"log"
	"testing"
)

var Code = ""

// 测试创建礼品码
func TestCreateGiftHandler(t *testing.T) {
	var formData model.CreateGiftModels
	code, err := CreateGiftHandler(formData)
	if err != nil {
		log.Println("err 测试创建礼品码")
	}
	if len(code) == 8 {
		log.Println("ok 测试创建礼品码")
	}
	Code = code
	return
}

// 测试获取礼品码
func TestGetGiftHandler(t *testing.T) {
	resData, err := GetGiftHandler(Code)
	if err != nil {
		log.Println("err 测试获取礼品码")
	}
	if len(resData) < 1 {
		log.Println("ok 测试获取礼品码")
	}
	return
}
