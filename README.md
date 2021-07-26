## 1.整体框架

Mongo与Protobuf使用,基于三

（1）【客户端调用, http】新增注册与登录接口：客户端传递唯一识别码（一个任意字符串）至服务器，服务器通过该识别码判断是否存在该玩家：不存在则注册新用户，生成唯一UID；存在则返回用户登陆数据（唯一UID、金币数、钻石数）。玩家信息储存在mongo数据库中

（2）【客户端调用, http】验证礼品码接口修改：按照管理员所添加的金币与钻石奖励数目，发放奖励存储至数据库。编译protobuf文件，将返回信息封装为protobuf对象以 **[]byte** 作为接口返回值返回给客户端。客户端接收到的是二进制序列，可以编写单测函数通过protobuf的decode方法解析，自测内容正确性。

## 2.目录结构

```
目录：
liqianqian@liqianqian MongoPlayer % pwd
/Users/liqianqian/go/src/MongoPlayer
项目结构分析：
liqianqian@liqianqian MongoPlayer % tree
.
├── README.md
├── app
│   ├── http
│   │   └── httpServer.go				#http 启动
│   └── main.go				#入口
├── config				#配置文件
│   └── app.ini
├── go.mod
├── go.sum
├── internal
│   ├── ctrl
│   │   └── giftCtrl.go				#礼品码，登陆注册控制器
│   ├── handler
│   │   ├── giftHandler_test.go				#单元测试
│   │   └── gigtHandler.go				#礼品码，注册登陆逻辑
│   ├── model				#models 模型层
│   │   ├── giftModel.go
│   │   ├── loginModel.go
│   │   ├── mongodb.go
│   │   └── redis.go
│   ├── myerr				#错误返回
│   │   └── err.go
│   └── router				#路由
│       └── router.go
│   └── service
│       └── service.go
├── locust				#压测
│   ├── __pycache__
│   │   ├── load.cpython-37.pyc
│   │   └── locust.cpython-37.pyc
│   ├── load.py				#压测脚本
│   └── report_1626961311.935244.html
├── protoFile				#消息协议
│   ├── generalReward.proto
│   └── protoFile
│       └── generalReward.pb.go
├── test				#测试demo
│   └── test.go
└── utils				#工具
│    └── tool.go
└── 登陆和验证礼品码流程图.jpg
14 directories, 22 files

```

## 3.逻辑代码分层

|    层     | 文件夹                           | 主要职责                 | 调用关系                  | 其它说明     |
| :-------: | :------------------------------- | ------------------------ | ------------------------- | ------------ |
|  应用层   | /app/http/shttpServer.go         | http 服务器启动          | 调用路由层                | 不可同层调用 |
|  路由层   | /internal/router/router.go       | 路由转发                 | 被应用层调用，调用控制层  | 不可同层调用 |
|  控制层   | /internal/ctrl/giftCtrl,go       | 礼品码管理，玩家注册登陆 | 被路由层调用，调用handler | 不可同层调用 |
| handler层 | /internal/handler/giftHandler.go | 处理具体业务             | 被控制层调用              | 不可同层调   |
|   model   | /internal/model                  | reids,mongodb 数据处理   | 被控制层调用              |              |
| 压力测试  | Locust/load.py                   | 进行压力测试             | 无调用关系                | 不可同层调用 |

## 4.存储设计

礼品码信息：

| 内容         | 数据库 | Key        | 类型 | 说明 |
| ------------ | ------ | ---------- | ---- | ---- |
| 礼品码       | redis  | Code       | Hash |      |
| 礼品码类型   | redis  | CodeType   | Hash |      |
| 可领取次数   | redis  | DrawCount  | Hash |      |
| 有效期时间戳 | redis  | ValidTime  | Hash |      |
| 奖品内容     | redis  | Content    | Hash |      |
| 管理员       | redis  | CreateUser | Hash |      |
| 已领取次数   | redis  | CostCount  | Hash |      |
| 指定玩家     | redis  | UserId     | Hash |      |
|              |        |            |      |      |

礼品码：限制次数类型的存储

| 内容   | 数据库 | Key       | 类型   | 说明 |      |
| ------ | ------ | --------- | ------ | ---- | ---- |
| 礼品码 | redis  | Code+type | String |      |      |
|        |        |           |        |      |      |

玩家数据：

| 内容           | 数据库   | Key     | 说明    |      |
| -------------- | -------- | ------- | ------- | ---- |
| UID            | Mongoldb | UID     | 唯一UID |      |
| 金币           | Mongoldb | Gold    |         |      |
| 钻石           | Mongoldb | Diamond |         |      |
| 客户端唯一标示 | Mongoldb | UniCode |         |      |
|                |          |         |         |      |



## 5.接口设计供客户端调用的接口

5.1玩家登陆

请求方法

http post 

接口地址：

127.0.0.1:8000/playerLogin

请求参数：

```
{
	"uniCode":uniCode001,			//	客户端唯一表示
}
```

json

请求响应

```
{
	"code": 0,
	"data": {
		"UniCode": "uniCode001",												//	客户端唯一表示
		"UID": "b329b2f0-962e-489d-ba1a-9cb662efcdc4", 	//UID
		"Gold": 0,																			// 金币
		"Diamond": 0																	  // 钻石
	},
	"msg": "新增注册与登录接口,新用户"
}
```

响应状态码

| 状态码 | 说明     |
| ------ | -------- |
| 0      | 创建成功 |
| 1      | 创建失败 |

5.2验证礼品码：

请求方法

http post 

接口地址：

127.0.0.1:8000/checkCode

请求参数：

```
{
		"code":90KKHauh,			//	礼品码
		"uid":8a601a2f-e101-437a-baa4-af37783c38f7,			    //	用户ID
}
```

json

请求响应

```
{
	"code": 0,
	"msg": "验证礼品码成功",
	"changes": {1:1000,2:10000},  // 变化量，礼品内容 1-金币，2-钻石
	"balance": {1:2,2:4},  // 变化前，礼品内容 1-金币，2-钻石
	"counter": {1:1002,2:10004},  // 变化后，礼品内容 1-金币，2-钻石
}
```

响应状态码

| 状态码 | 说明           |
| ------ | -------------- |
| 0      | 验证礼品码成功 |
| 1      | 验证礼品码失败 |

## 6.第三方库

gin

```
用于api服务，go web 框架
代码： github.com/gin-gonic/gin

```

proto

```
用于消息数据协议
包含：proto.Unmarshal，proto.Marshal 数据序列化
代码："github.com/golang/protobuf/proto"

```

redis

```
用于礼品码数据存储
包含：hash，string 
代码："github.com/gomodule/redigo/redis"
```

mongodb

```
用于玩家数据存储
代码：go.mongodb.org/mongo-driver/mongo
```



## 7.如何编译执行

```
#切换主目录下
cd ./app/
#编译
go build
```

## 8.todo 

```
后续优化，连接验证
```



