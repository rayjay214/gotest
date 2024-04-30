package main

import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"gopkg.in/ini.v1"
	"time"
)

func main() {
	// 创建阿里云客户端
	cfg, err := ini.Load("config.ini")
	section := cfg.Section("general")

	ak := section.Key("key").Value()
	sk := section.Key("secret").Value()

	aliClient, err := sdk.NewClientWithAccessKey("cn-hangzhou", ak, sk)
	if err != nil {
		panic(err)
	}

	strImei := "12345123451"
	strAlarm := "低电报警"
	type Param struct {
		Imei  string `json:"imei"`
		Alarm string `json:"alarm"`
	}
	param := Param{
		Imei:  fmt.Sprintf("%v于%v", strImei, time.Now().Format("2006-01-02 15:04:05")),
		Alarm: strAlarm,
	}
	jsonStr, _ := json.Marshal(param)
	fmt.Println(string(jsonStr))

	// 创建api请求并设置参数
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Domain = "dyvmsapi.aliyuncs.com"
	request.Version = "2017-05-25"
	request.ApiName = "SingleCallByTts"
	request.QueryParams["CalledNumber"] = "18956280556" // 被叫号码
	request.QueryParams["TtsCode"] = "TTS_260900017"    // 模板CODE
	request.QueryParams["PlayTimes"] = "1"              // 播放次数
	request.QueryParams["TtsParam"] = string(jsonStr)

	// 发送请求
	response, err := aliClient.ProcessCommonRequest(request)
	if err != nil {
		panic(err)
	}

	// 打印响应结果
	fmt.Println(response.GetHttpContentString())
}
