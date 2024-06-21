package model

import (
	"errors"
	"fmt"
	"reflect"
	"time"
)

type PayOrder struct {
	OrderId       string `json:"orderId" gorm:"primaryKey;type:varchar(128);comment:内部订单号"`
	TransactionId string `json:"transactionId" gorm:"type:varchar(128);comment:支付平台交易id"`
	Imei          int64  `json:"imei" gorm:"type:bigint;comment:设备号"`
	GoodsDesc     string `json:"goodsDesc" gorm:"type:varchar(32);comment:商品描述"`
	PayPlat       string `json:"payPlat" gorm:"type:varchar(4);comment:支付平台 字典：sim_pay_plat"`
	PayPrice      int64  `json:"payPrice" gorm:"type:bigint;comment:支付金额"`
	PackageId     int64  `json:"packageId" gorm:"type:bigint;comment:套餐id"`
	Remark        string `json:"remark" gorm:"type:varchar(512);comment:备注信息"`
	Status        string `json:"status" gorm:"type:varchar(4);comment:支付状态"`
	//CreateBy      int64     `json:"createBy" gorm:"index;comment:创建者"`
	//UpdateBy      int64     `json:"updateBy" gorm:"index;comment:更新者"`
	//CreatedAt     time.Time `json:"createdAt" gorm:"comment:创建时间"`
	//UpdatedAt     time.Time `json:"updatedAt" gorm:"comment:最后更新时间"`
	//DeletedAt     time.Time `json:"deletedAt" gorm:"index;comment:删除时间"`
}

type Device struct {
	Imei             int64     `json:"imei" gorm:"primaryKey;type:bigint unsigned;comment:设备号" db:"imei"`
	Pwd              string    `json:"pwd" gorm:"type:varchar(64);comment:设备号登录密码"`
	CarName          string    `json:"carName" gorm:"type:varchar(32);comment:设备名称"`
	StartTime        time.Time `json:"startTime" gorm:"type:datetime;comment:设备激活时间"`
	Expiration       int64     `json:"expiration" gorm:"type:bigint;comment:设备使用期限"`
	DeviceType       string    `json:"deviceType" gorm:"type:varchar(128);comment:设备型号" db:"device_type"`
	Iccid            string    `json:"iccid" gorm:"type:varchar(64);comment:设备上报ICCID"`
	BindPhone        string    `json:"bindPhone" gorm:"type:varchar(64);comment:绑定手机号"`
	Switch           int64     `json:"switch" gorm:"type:bigint unsigned;comment:增值开关"`
	FamilyId         int64     `json:"familyId" gorm:"type:bigint unsigned;comment:所属的familyid"`
	GroupId          int64     `json:"groupId" gorm:"type:bigint unsigned;comment:所属的分组"`
	TerminalFamilyId int64     `json:"terminalFamilyId" gorm:"type:bigint unsigned;comment:所属的终端familyid"`
	TerminalGroupId  int64     `json:"terminalGroupId" gorm:"type:bigint unsigned;comment:所属的终端分组id"`
	Mode             string    `json:"mode" gorm:"type:varchar(4);comment:设备模式"`
	Status           string    `json:"status" gorm:"type:varchar(4);comment:设备状态"`
	Protocol         string    `json:"protocol" gorm:"type:varchar(4);comment:设备协议"`
	ShakeValue       int32     `json:"shakeValue" gorm:"type:int;comment:震动敏感度"`
	//CreateBy   int64     `json:"createBy" gorm:"index;comment:创建者"`
	//UpdateBy   int64     `json:"updateBy" gorm:"index;comment:更新者"`
	//CreatedAt  time.Time `json:"createdAt" gorm:"comment:创建时间"`
	//UpdatedAt  time.Time `json:"updatedAt" gorm:"comment:最后更新时间"`
	//DeletedAt  time.Time `json:"deletedAt" gorm:"index;comment:删除时间"`
}

type Fence struct {
	Id          int64  `json:"id" comment:"id"`
	Name        string `json:"name" gorm:"type:varchar(500);comment:围栏名称"`
	FenceType   string `json:"fenceType" gorm:"type:varchar(4);comment:围栏类型, 字典：fence_type" db:"fence_type"`
	Imei        int64  `json:"imei" gorm:"type:bigint unsigned;comment:车载设备IMEI"`
	Value       string `json:"value" gorm:"type:varchar(3000);comment:json串"`
	FenceSwitch string `json:"fenceSwitch" gorm:"type:varchar(4);comment:围栏报警开关, 字典: fence_switch" db:"fence_switch"`
}

func MapToStruct(data []interface{}, structToMap interface{}) error {
	v := reflect.ValueOf(structToMap)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return errors.New("s must be a pointer to struct")
	}

	t := reflect.TypeOf(structToMap).Elem()
	for i := 0; i < t.NumField(); i++ {
		fieldName := t.Field(i).Name
		fieldValue := data[i]

		if fieldValue == nil {
			continue
		}

		var valueToSet interface{}

		fisInList := func(str string, list []string) bool {
			for _, s := range list {
				if str == s {
					return true
				}
			}
			return false
		}
		tranTimeList := []string{"CreatedAt", "UpdatedAt", "DeletedAt", "StartTime", "EndTime"}

		//时间实际上应该是东八区，但是打印的是UTC
		//todo fix
		if fisInList(fieldName, tranTimeList) {
			timeStr := fieldValue.(string)
			valueToSet, _ = time.Parse("2006-01-02 15:04:05.000", timeStr)
		} else {
			valueToSet = fieldValue
		}

		field := v.Elem().Field(i)
		if !field.CanSet() {
			return errors.New("field cannot be set")
		}

		value := reflect.ValueOf(valueToSet)
		if !value.Type().AssignableTo(field.Type()) {
			msg := fmt.Sprintf("field %s cannot be assigned type %s, idx %d", fieldName, value.Type(), i)
			return errors.New(msg)
		}

		field.Set(value)
	}

	return nil
}

func DiffStructFields(a, b interface{}) []string {
	var diffs []string

	typeOf := reflect.TypeOf(a)

	for i := 0; i < typeOf.NumField(); i++ {
		field := typeOf.Field(i)
		fieldName := field.Name

		// 获取字段的值
		valueA := reflect.ValueOf(a).Field(i)
		valueB := reflect.ValueOf(b).Field(i)

		// 比较字段值是否相等
		if !reflect.DeepEqual(valueA.Interface(), valueB.Interface()) {
			diff := fmt.Sprintf("%s: %v -> %v", fieldName, valueA.Interface(), valueB.Interface())
			diffs = append(diffs, diff)
		}
	}

	return diffs
}
