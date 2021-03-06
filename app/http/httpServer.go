package http

import (
	"MongoPlayer/internal/model"
	"MongoPlayer/internal/router"
	"github.com/gin-gonic/gin"
)

// http 服务
func HttpServer() {
	// 读取配置端口
	httpPort, _ := model.GetAppPort()
	//由于是外部调用包，所以必须含包名 gin. 作为前缀
	//Default 返回带有已连接 Logger 和 Recovery 中间件的 Engine 实例。
	r := gin.Default()
	// 默认是 0.0.0.0:8080 端口，内部使用了 http.ListenAndServe(address, engine)
	// Engine 结构体中内嵌了 RouterGroup 结构体，即继承了 RouterGroup（其有成员方法 GET、POST、DELETE、PUT、ANY 等）
	router.Router(r)
	r.Run("0.0.0.0:" + httpPort) // listen and serve on 0.0.0.0:8000
}
