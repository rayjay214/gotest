package main

import (
	"context"
	"fmt"
	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-mysql-org/go-mysql/replication"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"gotest/binlog/model"
	"os"
	"strconv"
)

var (
	rdb  *redis.Client
	pipe *redis.Pipeline
	db   *sqlx.DB
)

func initStorage() {
	var err error
	rdb = redis.NewClient(&redis.Options{
		Addr: "47.107.69.24:6480",
		DB:   0,
	})
	pipe = rdb.Pipeline().(*redis.Pipeline)

	address := fmt.Sprintf("%v:%v@tcp(%v)/xx?parseTime=true&timeout=5s&readTimeout=10s", "admin", "shht", "47.107.69.24:8000")
	db, err = sqlx.Open("mysql", address)
	if err != nil {
		log.Error("open mysql failed,", err)
		return
	}
}

func loadDevice() {
	var devices []model.Device

	err := db.Select(&devices, "select imei,device_type,protocol from device")
	if err != nil {
		log.Error("exec failed, ", err)
		return
	}

	for _, v := range devices {
		log.Println(v.Imei, v.DeviceType)
		key := fmt.Sprintf("imei_%v", v.Imei)
		pipe.HSet(context.Background(), key, map[string]interface{}{
			"device_type": v.DeviceType,
			"protocol":    v.Protocol,
			"iccid":       v.Iccid,
		})
	}
	pipe.Exec(context.Background())
}

func loadFence() {
	var fenceList []model.Fence

	err := db.Select(&fenceList, "select id,name,fence_type,imei,value,fence_switch from fence")
	if err != nil {
		log.Error("exec failed, ", err)
		return
	}

	for _, fence := range fenceList {
		setKey := "fenceset"
		infoKey := fmt.Sprintf("fenceinfo_%v", fence.Imei)
		pipe.SAdd(context.Background(), setKey, fence.Imei)
		pipe.HSet(context.Background(), infoKey, map[string]interface{}{
			strconv.FormatInt(fence.Id, 10): fence.Value,
		})
	}
	pipe.Exec(context.Background())
}

func syncFromBinlog() {
	// 创建一个 BinlogSyncer 实例
	cfg := replication.BinlogSyncerConfig{
		ServerID: 100,
		Flavor:   "mysql",
		Host:     "47.107.69.24",
		Port:     8000,
		User:     "admin",
		Password: "shht",
	}
	syncer := replication.NewBinlogSyncer(cfg)

	// 开始同步 binlog
	streamer, err := syncer.StartSync(mysql.Position{
		Name: "mysql-bin.000007",
		Pos:  132983966,
	})
	if err != nil {
		fmt.Println("Error starting syncer:", err)
		os.Exit(1)
	}

	// 监听 binlog 事件
	for {
		ev, err := streamer.GetEvent(context.Background())
		if err != nil {
			fmt.Println("Error getting event:", err)
			continue
		}

		// 处理不同类型的 binlog 事件
		switch e := ev.Event.(type) {
		case *replication.RowsEvent:
			// 处理行数据变更事件
			if string(e.Table.Schema) == "xx" {
				if string(e.Table.Table) == "device" {
					switch ev.Header.EventType {
					case replication.WRITE_ROWS_EVENTv2:
						syncInsertDevice(e)
					case replication.UPDATE_ROWS_EVENTv2:
						syncUpdateDevice(e)
					case replication.DELETE_ROWS_EVENTv2:
						syncDeleteDevice(e)
					}
				}
				if string(e.Table.Table) == "fence" {
					switch ev.Header.EventType {
					case replication.WRITE_ROWS_EVENTv2:
						syncInsertFence(e)
					case replication.UPDATE_ROWS_EVENTv2:
						syncUpdateFence(e)
					case replication.DELETE_ROWS_EVENTv2:
						syncDeleteFence(e)
					}
				}
			}
		case *replication.QueryEvent:
		default:
		}
	}
}

func syncInsertDevice(e *replication.RowsEvent) {
	for _, row := range e.Rows {
		var device model.Device
		err := model.MapToStruct(row, &device)
		if err != nil {
			fmt.Println(err)
			break
		}
		key := fmt.Sprintf("imei_%v", device.Imei)
		pipe.HSet(context.Background(), key, map[string]interface{}{
			"device_type": device.DeviceType,
			"protocol":    device.Protocol,
			"iccid":       device.Iccid,
		})
	}
	pipe.Exec(context.Background())
}

func syncUpdateDevice(e *replication.RowsEvent) {
	for i := 0; i < len(e.Rows); i += 2 {
		var oldDevice, newDevice model.Device
		err := model.MapToStruct(e.Rows[i], &oldDevice)
		if err != nil {
			fmt.Println(err)
			break
		}
		err = model.MapToStruct(e.Rows[i+1], &newDevice)
		if err != nil {
			fmt.Println(err)
			break
		}
		if newDevice != oldDevice {
			key := fmt.Sprintf("imei_%v", newDevice.Imei)
			pipe.HSet(context.Background(), key, map[string]interface{}{
				"device_type": newDevice.DeviceType,
				"protocol":    newDevice.Protocol,
				"iccid":       newDevice.Iccid,
			})
		}
	}
	pipe.Exec(context.Background())
}

func syncDeleteDevice(e *replication.RowsEvent) {
	for _, row := range e.Rows {
		var device model.Device
		err := model.MapToStruct(row, &device)
		if err != nil {
			fmt.Println(err)
			break
		}
		key := fmt.Sprintf("imei_%v", device.Imei)
		pipe.Del(context.Background(), key)
	}
	pipe.Exec(context.Background())
}

func syncInsertFence(e *replication.RowsEvent) {
	for _, row := range e.Rows {
		var fence model.Fence
		err := model.MapToStruct(row, &fence)
		if err != nil {
			fmt.Println(err)
			break
		}
		setKey := "fenceset"
		infoKey := fmt.Sprintf("fenceinfo_%v", fence.Imei)

		pipe.SAdd(context.Background(), setKey, fence.Imei)
		pipe.HSet(context.Background(), infoKey, map[string]interface{}{
			strconv.FormatInt(fence.Id, 10): fence.Value,
		})
	}
	pipe.Exec(context.Background())
}

func syncUpdateFence(e *replication.RowsEvent) {
	for i := 0; i < len(e.Rows); i += 2 {
		var oldFence, newFence model.Fence
		err := model.MapToStruct(e.Rows[i], &oldFence)
		if err != nil {
			fmt.Println(err)
			break
		}
		err = model.MapToStruct(e.Rows[i+1], &newFence)
		if err != nil {
			fmt.Println(err)
			break
		}
		if oldFence != newFence {
			infoKey := fmt.Sprintf("fenceinfo_%v", newFence.Imei)
			//不能修改围栏状态

			pipe.HSet(context.Background(), infoKey, map[string]interface{}{
				strconv.FormatInt(newFence.Id, 10): newFence.Value,
			})
		}
	}
	pipe.Exec(context.Background())
}

func syncDeleteFence(e *replication.RowsEvent) {
	for _, row := range e.Rows {
		var fence model.Fence
		err := model.MapToStruct(row, &fence)
		if err != nil {
			fmt.Println(err)
			break
		}
		setKey := "fenceset"
		infoKey := fmt.Sprintf("fenceinfo_%v", fence.Imei)
		pipe.HDel(context.Background(), infoKey, strconv.FormatInt(fence.Id, 10))
		hashLen, _ := rdb.HLen(context.Background(), "test").Result()
		if hashLen <= 0 {
			pipe.SRem(context.Background(), setKey, fence.Imei)
		}
	}
	pipe.Exec(context.Background())
}

func main() {
	initStorage()
	loadDevice()
	loadFence()
	syncFromBinlog()
}
