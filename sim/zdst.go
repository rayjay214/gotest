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
	clientId           = "4028808f81003d8801811814eab5406a" //qx
	baseUrl            = "http://cmpapi.zdm2m.com/"
	getMethod          = "smsApi.do?querycard"
	queryPackageMethod = "smsApi.do?batchQueryCardsRenewPackageInfo"
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

func main() {
	ZdGet("1441594933091")
	ZdQueryPackage("1441594933091")
}
