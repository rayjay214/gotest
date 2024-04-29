package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
	"net/smtp"
	"strings"
	"time"
)

func SendToMail(user, password, host, subject, body, mailtype string, to []string) error {
	hp := strings.Split(host, ":")
	auth := smtp.PlainAuth("", user, password, hp[0])
	var content_type string
	if mailtype == "html" {
		content_type = "Content-Type: text/" + mailtype + "; charset=UTF-8"
	} else {
		content_type = "Content-Type: text/plain" + "; charset=UTF-8"
	}

	to_address := strings.Join(to, ";")
	msg := []byte("To: " + to_address + "\r\nFrom: " + user + "\r\nSubject: " + subject + "\r\n" + content_type + "\r\n\r\n" + body)

	err := smtp.SendMail(host, auth, user, to, msg)
	return err
}

func getValue(str, findS string) (int, string) {
	sLeft := ""
	if str == "" {
		return 0, sLeft
	}

	nPos := 0
	if findS != "" {
		nPos = strings.Index(str, findS)
	}
	if nPos < 0 {
		return 0, sLeft
	}

	sLeft = str[nPos+len(findS):]
	sLeft = strings.TrimLeft(sLeft, " ")
	sLeft = strings.TrimLeft(sLeft, "\r\n")

	nPos = strings.Index(sLeft, "\r\n")
	if nPos >= 0 {
		sLeft = sLeft[:nPos]
	}

	sLeft = strings.TrimRight(sLeft, " ")
	sLeft = strings.TrimRight(sLeft, "\r\n")

	return len(sLeft), sLeft
}

func test() {
	//str := "\r\nAT+GPS\r\n\r\n10\r\nOK"
	str := "\r\nAT+STATUS\r\nsimReady:1,sleepStatus:1,CSQ:11,ACC:0,gpsStarNum:0,ADC:4389,Voltage:4389mV,Gsensor:0,IP:vip.gps666.net,PORT:8885,sensorMotionId:0\r\n"
	_, value := getValue(str, "AT+STATUS")
	fmt.Println(value)
}

func test1() {
	var err error
	MysqlDB, err := sql.Open("mysql", "admin:shht@tcp(47.107.69.24:8000)/xx?charset=utf8&parseTime=True&loc=Local")

	if err != nil {
		return
	}
	// set pool params
	MysqlDB.SetMaxOpenConns(2000)
	MysqlDB.SetMaxIdleConns(1000)
	MysqlDB.SetConnMaxLifetime(time.Minute * 60)

	var total, used int
	row := MysqlDB.QueryRow("SELECT total, used FROM additional_service WHERE imei=? AND start_time<? AND end_time>? "+
		"AND (total-used>0) AND service_type=? order by end_time limit 1",
		9246243122, time.Now(), time.Now(), 2)

	err = row.Scan(&total, &used)
	fmt.Println(total, used, err)
}

func test2() {
	uid := "1AA01BA798D0"
	rdb := redis.NewClient(&redis.Options{
		Addr: "114.215.190.173:6480",
		DB:   0,
	})
	key := fmt.Sprintf("ipcinfo_%v", uid)
	result, _ := rdb.HGetAll(context.Background(), key).Result()
	fmt.Println(result)
}

func test3() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "114.215.190.173:6480",
		DB:   0,
	})
	rdb.HSet(context.Background(), "test", map[string]interface{}{
		"1": "aaa",
	})
	rdb.HSet(context.Background(), "test", map[string]interface{}{
		"2": "bbb",
	})
	hashLen, _ := rdb.HLen(context.Background(), "test").Result()
	fmt.Println(hashLen)
	rdb.HDel(context.Background(), "test", "1")
	hashLen, _ = rdb.HLen(context.Background(), "test").Result()
	fmt.Println(hashLen)
	//rdb.HDel(context.Background(), "test", "2")
	hashLen, _ = rdb.HLen(context.Background(), "test").Result()
	fmt.Println(hashLen)
}

func main() {
	//test()
	//test1()
	//test2()
	test3()

	/*
		user := "ipc@email.gps666.net"
		password := "FCst20221128"
		host := "smtpdm.aliyun.com:80"
		to := []string{"526528945@qq.com"}

		subject := "test Golang to sendmail"
		mailtype := "text"

		body := "test"
		fmt.Println("send email")
		err := SendToMail(user, password, host, subject, body, mailtype, to)
		if err != nil {
			fmt.Println("Send mail error!")
			fmt.Println(err)
		} else {
			fmt.Println("Send mail success!")
		}
	*/
}
