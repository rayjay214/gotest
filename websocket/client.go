package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

type ReqMessage struct {
	Method   string `json:"method"`
	Uid      string `json:"uid"`
	Begin    uint32 `json:"begin"`
	End      uint32 `json:"end"`
	FileName uint32 `json:"fileName"`
}

type RespMessage struct {
	Method   string `json:"method"`
	Files    []File `json:"files"`
	FileName uint32 `json:"fileName"`
	Uid      string `json:"uid"`
	Errcode  int    `json:"errcode"`
}

type File struct {
	Begin uint32 `json:"begin"`
	Cos   uint32 `json:"cos"`
}

func makeMessage() []byte {
	type Message struct {
		Method string `json:"method"`
		Uid    string `json:"uid"`
		Begin  uint32 `json:"begin"`
		End    uint32 `json:"end"`
	}

	message := Message{
		Method: "GetSdFileList",
		Uid:    "12345678",
		Begin:  1714790592,
		End:    1714790593,
	}

	json, _ := json.Marshal(message)
	return json
}

func make8201(uid string, beginTime, endTime uint) []byte {
	type Message struct {
		Method string `json:"method"`
		Uid    string `json:"uid"`
		Begin  uint32 `json:"begin"`
		End    uint32 `json:"end"`
	}
	message := Message{
		Method: "GetSdFileList",
		Uid:    uid,
		Begin:  uint32(beginTime),
		End:    uint32(endTime),
	}

	json, _ := json.Marshal(message)
	return json
}

func make8202(uid string, fileName uint) []byte {
	type Message struct {
		Method   string `json:"method"`
		Uid      string `json:"uid"`
		FileName uint32 `json:"fileName"`
	}
	message := Message{
		Method:   "GetSdFile",
		Uid:      uid,
		FileName: uint32(fileName),
	}

	json, _ := json.Marshal(message)
	return json
}

func make8206(uid string, month uint) []byte {
	type Message struct {
		Method string `json:"method"`
		Uid    string `json:"uid"`
		Month  uint32 `json:"month"`
	}
	message := Message{
		Method: "MonthFileSummary",
		Uid:    uid,
		Month:  uint32(month),
	}

	json, _ := json.Marshal(message)
	return json
}

func make8207(uid string, filename uint) []byte {
	type Message struct {
		Method   string `json:"method"`
		Uid      string `json:"uid"`
		FileName uint32 `json:"fileName"`
	}
	message := Message{
		Method:   "CancelTask",
		Uid:      uid,
		FileName: uint32(filename),
	}

	json, _ := json.Marshal(message)
	return json
}

var (
	addr      string
	uid       string
	appid     string
	beginTime uint
	endTime   uint
	fileName  uint
	msgId     uint
	month     uint
)

func main() {
	flag.StringVar(&addr, "addr", "localhost:7080", "http service address")
	flag.StringVar(&uid, "uid", "123456789111111", "uid")
	flag.StringVar(&appid, "appid", "abc", "appid")
	flag.UintVar(&beginTime, "beginTime", 0, "beginTime")
	flag.UintVar(&endTime, "endTime", 0, "endTime")
	flag.UintVar(&fileName, "fileName", 0, "fileName")
	flag.UintVar(&msgId, "msgId", 0, "msgId")
	flag.UintVar(&month, "month", 0, "month")

	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: addr, Path: "/ws"}
	q := u.Query()
	q.Set("appId", appid)
	u.RawQuery = q.Encode()

	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	var msg []byte
	if msgId == 0x8201 {
		msg = make8201(uid, beginTime, endTime)
	} else if msgId == 0x8202 {
		msg = make8202(uid, fileName)
	} else if msgId == 0x8206 {
		msg = make8206(uid, month)
	} else if msgId == 0x8207 {
		msg = make8207(uid, fileName)
	} else {
		return
	}

	err = c.WriteMessage(websocket.TextMessage, msg)
	if err != nil {
		log.Println("write:", err)
		return
	}

	var file *os.File

	go func() {
		defer close(done)
		for {
			mt, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			if mt == websocket.BinaryMessage {
				log.Printf("receive binary, len %v", len(message))
				file.Write(message)
			} else {
				var resp RespMessage
				err = json.Unmarshal(message, &resp)
				if err != nil {
					log.Println(err)
					break
				}
				if resp.Method == "BeginSendFile" {
					fileName := fmt.Sprintf("%v.avi", resp.FileName)
					file, err = os.Create(fileName)
					if err != nil {
						fmt.Println("Error opening file:", err)
					}
				}
				if resp.Method == "SendNotify" {
					fmt.Printf("receive finish %v:%v\n", resp.Uid, resp.Errcode)
				}
			}
		}
	}()

	ticker := time.NewTicker(time.Second * 300)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		//case t := <-ticker.C:
		case <-ticker.C:
			json := makeMessage()
			err := c.WriteMessage(websocket.TextMessage, json)
			if err != nil {
				log.Println("write:", err)
				return
			}
		case <-interrupt:
			log.Println("interrupt")
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
