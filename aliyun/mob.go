package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func main() {
	sendUrl := "https://webapi.sms.mob.com/sms/sendmsg"

	payload := url.Values{}
	payload.Set("appkey", "36f7446b78c59")
	payload.Set("phone", "18956280556")
	payload.Set("zone", "86")
	payload.Set("templateCode", "3076917")

	req, err := http.NewRequest("POST", sendUrl, bytes.NewBufferString(payload.Encode()))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	fmt.Println(string(body))
}
