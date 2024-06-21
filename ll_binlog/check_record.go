package main

import (
    "context"
    "fmt"
    "gotest/binlog/model"
    "os"
    "time"

    "github.com/go-mysql-org/go-mysql/mysql"
    "github.com/go-mysql-org/go-mysql/replication"
)

var (
    gKeys = int64(868886005747228)
)

func main() {
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
        Pos:  0,
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
            operTime := time.Unix(int64(ev.Header.Timestamp), 0)
            if string(e.Table.Schema) == "gps" {
                switch ev.Header.EventType {
                case replication.WRITE_ROWS_EVENTv2:
                    LogInsertEvent(e, operTime)
                case replication.UPDATE_ROWS_EVENTv2:
                    LogUpdateEvent(e, operTime)
                case replication.DELETE_ROWS_EVENTv2:
                    LogDeleteEvent(e, operTime)
                }
            }
        case *replication.QueryEvent:

        default:

        }
    }
}

func LogInsertEvent(e *replication.RowsEvent, time time.Time) {
    table := string(e.Table.Table)
    switch table {
    case "device":
        for _, row := range e.Rows {
            var device model.Device
            err := model.MapToStruct(row, &device)
            if err != nil {
                fmt.Println(err)
                break
            }
            if device.Imei == gKeys {
                fmt.Printf("%s insert device: %+v\n", time, device)
            }
        }
    case "pay_order":
        for _, row := range e.Rows {
            var order model.PayOrder
            err := model.MapToStruct(row, &order)
            if err != nil {
                fmt.Println(err)
                break
            }
            if order.Imei == gKeys {
                fmt.Printf("%s insert order: %+v\n", time, order)
            }
        }
    }
}

func LogUpdateEvent(e *replication.RowsEvent, time time.Time) {
    switch string(e.Table.Table) {
    case "device":
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
            if oldDevice.Imei == gKeys {
                changes := model.DiffStructFields(oldDevice, newDevice)
                //fmt.Printf("%s device change from %+v to %+v\n", time, oldDevice, newDevice)
                fmt.Printf("%s device change field %v\n", time, changes)
            }
        }
    case "pay_order":
        for i := 0; i < len(e.Rows); i += 2 {
            var oldOrder, newOrder model.PayOrder
            err := model.MapToStruct(e.Rows[i], &oldOrder)
            if err != nil {
                fmt.Println(err)
                break
            }
            err = model.MapToStruct(e.Rows[i+1], &newOrder)
            if err != nil {
                fmt.Println(err)
                break
            }
            if oldOrder.Imei == gKeys {
                changes := model.DiffStructFields(oldOrder, newOrder)
                //fmt.Printf("%s order change from %+v to %+v\n", time, oldOrder, newOrder)
                fmt.Printf("%s pay_order change field %v\n", time, changes)
            }
        }
    }
}

func LogDeleteEvent(event *replication.RowsEvent, time time.Time) {

}
