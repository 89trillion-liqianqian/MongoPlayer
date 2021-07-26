package main

import (
	"MongoPlayer/app/http"
	"MongoPlayer/internal/model"
)

// 入口
func main() {
	// 加载配置
	filepath := "../config/app.ini"
	model.GetAppIni(filepath)
	// 初始化 redis
	model.Init()
	// 初始化 mongo
	model.InitMongo()
	// 启动http server
	http.HttpServer()
}
