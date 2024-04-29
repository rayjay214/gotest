package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

var (
	//clientId           = "402880906dfd7de4016e0180834f19cb"  //xx
	//clientId                = "4028808f81003d8801811814eab5406a" //qx
	clientId                = "402880208e080e4c018e0d4844c12865" //ms
	baseUrl                 = "http://cmpapi.zdm2m.com/"
	getMethod               = "smsApi.do?querycard"
	queryPackageMethod      = "smsApi.do?batchQueryCardsRenewPackageInfo"
	packageRenewMethod      = "smsApi.do?packageRenew"
	additionalPackageMethod = "smsApi.do?additionalPackage"
)

type ZdPackageMsg struct {
	Total       float64 `json:"total"`
	Used        float64 `json:"used"`
	Allowance   float64 `json:"allowance"`
	PType       string  `json:"ptype"`
	EnableTime  string  `json:"enabletime"`
	FailureTime string  `json:"failuretime"`
}

type ZdCardMsg struct {
	CardNo            string  `json:"cardno"`
	Iccid             string  `json:"iccid"`
	OrderID           string  `json:"order_id"`
	OperatorType      string  `json:"operatortype"`
	Status            string  `json:"status"`
	SaleDate          string  `json:"saledate"`
	ValidDate         string  `json:"validdate"`
	ActivationDate    string  `json:"activationdate"`
	PerActivationDate string  `json:"peractivationdate"`
	CardAccount       float64 `json:"cardaccount"`
	ActiveStatus      string  `json:"activestatus"`
	Imsi              string  `json:"imsi"`
	SilentDateEndTime string  `json:"silentdateendtime"`
	OperatorTypeCode  string  `json:"operatortypecode"`
	PoolNo            string  `json:"poolno"`
}

type ZdGetResult struct {
	PackageMsg    []ZdPackageMsg `json:"packagemsg"`
	CardMsg       ZdCardMsg      `json:"cardmsg"`
	BalanceAmount float64        `json:"balanceamount"`
	ResultMsg     string         `json:"resultmsg"`
	ResultCode    string         `json:"resultcode"`
}

type ZdAdditionalPackageInfo struct {
	Price       float64 `json:"price"`
	PackageCode string  `json:"packageCode"`
	PackageName string  `json:"packageName"`
}

type ZdMinPackageInfo struct {
	Price       float64 `json:"price"`
	PackageCode string  `json:"packageCode"`
	PackageName string  `json:"packageName"`
}

type ZdRenewPackageInfo struct {
	CardNo                string                    `json:"cardno"`
	Iccid                 string                    `json:"iccid"`
	AdditionalPackageInfo []ZdAdditionalPackageInfo `json:"additionalPackageInfo"`
	MinPackageInfo        ZdMinPackageInfo          `json:"minPackageInfo"`
}

type ZdBatchRenewPackageInfo struct {
	Data       []ZdRenewPackageInfo `json:"data"`
	ResultMsg  string               `json:"resultmsg"`
	ResultCode string               `json:"resultcode"`
}

type ZdRenewResult struct {
	ResultMsg  string `json:"resultmsg"`
	ResultCode string `json:"resultcode"`
	Balance    string `json:"balance"`
	Cost       string `json:"cost"`
}

func genSign(signStr string) string {
	hash := md5.New()
	hash.Write([]byte(signStr))
	hashValue := hash.Sum(nil)
	hashString := hex.EncodeToString(hashValue)
	return strings.ToUpper(hashString)
}

func getUrl(url string) ([]byte, error) {
	fmt.Println(url)
	res, err := http.Get(url)
	if err != nil {
		return nil, errors.New("http get failed")
	}

	defer res.Body.Close()

	content, errs := ioutil.ReadAll(res.Body)
	if errs != nil {
		return nil, errors.New("parse failed")
	}

	return content, nil
}

// 单卡查询
func ZdGet(iccid string) (ZdGetResult, error) {
	var result ZdGetResult

	url := baseUrl + getMethod
	url += fmt.Sprintf("&clientid=%v", clientId)
	url += fmt.Sprintf("&cardno=%v", iccid)
	url += fmt.Sprintf("&sign=%v", genSign(fmt.Sprintf("clientid=%v&cardno=%v", clientId, iccid)))
	response, err := getUrl(url)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(response, &result)
	if err != nil {
		return result, err
	}
	fmt.Println(result)
	return result, nil
}

// 查询可续费套餐
func ZdQueryPackage(iccid string) (ZdBatchRenewPackageInfo, error) {
	var result ZdBatchRenewPackageInfo

	url := baseUrl + queryPackageMethod
	url += fmt.Sprintf("&clientid=%v", clientId)
	url += fmt.Sprintf("&cardno=%v", iccid)
	url += fmt.Sprintf("&sign=%v", genSign(fmt.Sprintf("clientid=%v&cardno=%v", clientId, iccid)))
	response, err := getUrl(url)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(response, &result)
	if err != nil {
		return result, err
	}
	fmt.Println(result)
	return result, nil
}

// 续费短信，基础套餐
func ZdPackageRenew(iccid string, genre, amount int) (ZdRenewResult, error) {
	var result ZdRenewResult

	url := baseUrl + packageRenewMethod
	url += fmt.Sprintf("&clientid=%v", clientId)
	url += fmt.Sprintf("&cardno=%v", iccid)
	url += fmt.Sprintf("&amount=%v", amount)
	url += fmt.Sprintf("&genre=%v", genre)
	url += fmt.Sprintf("&sign=%v", genSign(fmt.Sprintf("clientid=%v&cardno=%v&amount=%v&genre=%v", clientId, iccid, amount, genre)))
	response, err := getUrl(url)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(response, &result)
	if err != nil {
		return result, err
	}
	fmt.Println(result)
	return result, nil
}

// 续费加油包
func ZdAdditionalPackage(iccid, packagecode string) (ZdRenewResult, error) {
	var result ZdRenewResult

	url := baseUrl + additionalPackageMethod
	url += fmt.Sprintf("&clientid=%v", clientId)
	url += fmt.Sprintf("&cardno=%v", iccid)
	url += fmt.Sprintf("&packagecode=%v", packagecode)
	url += fmt.Sprintf("&sign=%v", genSign(fmt.Sprintf("clientid=%v&cardno=%v&packagecode=%v", clientId, iccid, packagecode)))
	response, err := getUrl(url)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(response, &result)
	if err != nil {
		return result, err
	}
	fmt.Println(result)
	return result, nil
}

func main() {
	ZdGet("898604B8262270090149")
	ZdQueryPackage("898604B8262270090149")
}
