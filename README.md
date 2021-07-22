# MongoPlayer
创建玩家信息存储结构、使用礼品码,Mongo与Protobuf使用,基于三

1.目录结构

```
目录：
liqianqian@liqianqian MongoPlayer % pwd
/Users/liqianqian/go/src/MongoPlayer
项目结构分析：
liqianqian@liqianqian MongoPlayer % tree
.
├── README.md					//技术文档
├── controller				// http api
│   └── base.go				// api 
├── go.mod
├── go.sum
├── locust						//
│   ├── __pycache__
│   │   ├── load.cpython-37.pyc
│   │   └── locust.cpython-37.pyc
│   ├── load.py				//压测脚步
│   └── report_1626961311.935244.html		//压测报告
├── main.go						//入口函数
├── model							//
│   └── requestModel.go		// model 礼品码模块
│   ├── loginModel.go			// model 登陆模块
│   ├── mongodb.go				// mongndb 初始化
│   ├── redis.go					// redis 初始化
├── protoFile							// proto 文件
│   ├── generalReward.proto //proto 文件
│   └── protoFile
│       └── generalReward.pb.go  proto 文件生成的.go文件
└── test
│   └── test.go				// 单元测试
└── utils
    └── tool.go       // 工具方法
8 directories, 17 files
liqianqian@liqianqian MongoPlayer % 

```

2。运行

```
go run main.go  
```

3.api 文档

3.1

```
1）玩家登陆
http post 
api: ip:port/playerLogin
请求体：
		"uniCode":uniCode001,			//	客户端唯一表示
响应体
json
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
状态码
0 ：创建成功
1 ：创建失败
```

3.2

```
客户端调用 - 验证礼品码：用户在客户端内输入礼品码并提交，如果礼品码合法且未被领取过，调用下方奖励接口，给用户增加奖励， 加奖励成功后，返回奖励内容供客户端内展示。
http post 
api: ip:port/checkCode
请求体：
		"code":90KKHauh,			//	礼品码
		"uid":8a601a2f-e101-437a-baa4-af37783c38f7,			    //	用户ID
响应体
protobuf model
{
	"code": 0,
	"msg": "验证礼品码成功",
	"changes": {1:1000,2:10000},  // 变化量，礼品内容 1-金币，2-钻石
	"balance": {1:2,2:4},  // 变化前，礼品内容 1-金币，2-钻石
	"counter": {1:1002,2:10004},  // 变化后，礼品内容 1-金币，2-钻石
}
状态码
0 ：验证礼品码成功
1 ：验证礼品码失败
```



