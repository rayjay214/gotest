package main

import (
    "context"
    "fmt"
    "os"

    "github.com/go-mysql-org/go-mysql/mysql"
    "github.com/go-mysql-org/go-mysql/replication"
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
            // 处理行数据变更事件
            if string(e.Table.Schema) == "xx" {
                fmt.Printf("Rows event: %s, time: %v\n", ev.Header.EventType, ev.Header.Timestamp)
                fmt.Printf("Schema: %s\n", e.Table.Schema)
                fmt.Printf("Table: %s\n", e.Table.Table)
                for _, row := range e.Rows {
                    fmt.Println("Row data:", row)
                }
            }
        case *replication.QueryEvent:
            // 处理 SQL 查询事件
            //fmt.Printf("database %s\n", e.Schema)
            //fmt.Printf("Query event: %s\n")

        default:
            // 其他类型的事件
            //fmt.Printf("Unknown event: %T\n", e)
        }
    }
}
