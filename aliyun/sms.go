package main

import (
	"github.com/aliyun-sdk/sms-go"
	"gopkg.in/ini.v1"
	"log"
)

var err error
var client *sms.Client
var ak, sk, sn, tc string

func init() {
	cfg, err := ini.Load("config.ini")
	section := cfg.Section("general")

	ak = section.Key("key").Value()
	sk = section.Key("secret").Value()
	sn = section.Key("sn").Value()
	tc = section.Key("tc").Value()
	client, err = sms.New(ak, sk, sms.SignName(sn), sms.Template(tc))
	if err != nil {
		log.Fatalln(err)
	}
}

func Send() {
	err = client.Send(
		sms.Mobile("19925233886"),
		sms.Parameter(map[string]string{
			"imei":  "1234567891",
			"alarm": "低电",
		}),
	)
	if err != nil {
		log.Println(err)
	}
}

func SendBatch() {
	items := []sms.BatchItem{
		{
			// "您的短信签名, 若已全局配置则可留空"
			Sign:   "",
			Mobile: "17757171482",
			Params: map[string]string{"code": "01234"},
		},
		{
			// "您的短信签名, 若已全局配置则可留空"
			Sign:   "",
			Mobile: "17757171483",
			Params: map[string]string{"code": "56789"},
		},
	}
	err = client.SendBatch(items)
	if err != nil {
		log.Println(err)
	}
}

func main() {
	Send()
}
