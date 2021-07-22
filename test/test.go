package main

import (
	"MongoPlayer/protoFile/protoFile"
	"bytes"
	"fmt"
	"github.com/golang/protobuf/proto"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

/**
api 测试
*/

// 发送GET请求
// url：         请求地址
// response：    请求返回的内容

func httpGet(url string) string {

	// 超时时间：5秒
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var buffer [512]byte
	result := bytes.NewBuffer(nil)
	for {
		n, err := resp.Body.Read(buffer[0:])
		result.Write(buffer[0:n])
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
	}

	return result.String()
}

// 发送Post请求
// url：         请求地址
// response：    请求返回的内容

func httpPost(urlStr string) string {

	//codeType:= "1"
	//codeType:= "2"
	codeType := "3"
	resp, err := http.PostForm(urlStr,
		url.Values{
			"codeType":   {codeType},
			"drawCount":  {"19"},
			"des":        {"des 这是金币的礼品码"},
			"validTime":  {"16345634354545"},
			"content":    {"{'1':1000,'2':10000}"},
			"createUser": {"qq"},
			"userId":     {"123456"},
		})

	if err != nil {
		// handle error
		log.Println("--resp err")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		log.Println("--ReadAll err")
	}

	fmt.Println(string(body))
	return string(body)
}

// 发送Post请求
// url：         请求地址
// response：    请求返回的内容

func httpPostCheck(urlStr, code, uid string) (resData protoFile.GeneralReward) {

	resp, err := http.PostForm(urlStr,
		url.Values{
			"code": {code},
			"uid":  {uid},
		})

	if err != nil {
		// handle error
		log.Println("--resp err")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		log.Println("--ReadAll err")
	} else {
		user := &protoFile.GeneralReward{}
		proto.UnmarshalMerge(body, user)
		resData = *user
	}
	fmt.Println("-aaaaaa", string(body))
	return
}

// 发送GET请求
// url：         请求地址
// response：    请求返回的内容

func httpGetProto(url string) (resData protoFile.GeneralReward) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	} else {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
		} else {
			user := &protoFile.GeneralReward{}
			proto.UnmarshalMerge(body, user)
			resData = *user
			log.Println("---httpGetProto=", *user)
		}

	}
	return
}

// 发送Post请求
// url：         请求地址
// response：    请求返回的内容

func httpPostLoginTest(urlStr, uniCode string) (resData protoFile.GeneralReward) {

	resp, err := http.PostForm(urlStr,
		url.Values{
			"uniCode": {uniCode},
		})

	if err != nil {
		// handle error
		log.Println("--resp err")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		log.Println("--ReadAll err")
	} else {
		user := &protoFile.GeneralReward{}
		proto.UnmarshalMerge(body, user)
		resData = *user
	}
	return
}

// 发送Post请求
// url：         请求地址
// response：    请求返回的内容

func httpPostLogin(urlStr, uniCode string) string {
	resp, err := http.PostForm(urlStr,
		url.Values{
			"uniCode": {uniCode},
		})

	if err != nil {
		// handle error
		log.Println("--resp err")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		log.Println("--ReadAll err")
	}

	fmt.Println(string(body))
	return string(body)
}

func main() {
	log.Println("--start test")

	result := ""
	urlStr := ""

	// 玩家登陆
	urlStr = "http://127.0.0.1:8000/playerLogin"

	resultData := httpPostLogin(urlStr, "uniCode001")
	log.Println("--结果", result, resultData)

	//// 验证礼品码
	urlStr = "http://127.0.0.1:8000/checkCode"
	// type 1-指定用户一次性消耗，2-不指定用户限制兑换次数，3-不限用户不限次数兑换
	//type=3 jG7a4lo8, type=2 90KKHauh,type=1 C72uloHO uid=123456
	//code:="jG7a4lo8"  //
	//code:="C72uloHO"  //
	//uid:="1234567"

	code := "90KKHauh" //
	uid := "8a601a2f-e101-437a-baa4-af37783c38f7"
	giftData := httpPostCheck(urlStr, code, uid)
	log.Println("--结果", giftData)
}
