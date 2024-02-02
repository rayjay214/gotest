package util

import (
    "fmt"

    //"github.com/bwmarrin/snowflake"
    "github.com/soloohu/snowflake"
)

var node *snowflake.Node
var baseId int64

func Init() {
    node, _ = snowflake.NewNode(int64(1))
    baseId = 1 * 100000000
}

func GenerateId() int64 {
    if baseId == 0 {
        Init()
    }
    id := node.Generate().Int64()
    ID := snowflake.ParseInt64(id)
    fmt.Printf("id:%d, Time:%d, Node:%d, Step:%d", id, ID.Time(), ID.Node(), ID.Step())
    return id
}

func GenerateIdString() string {
    if baseId == 0 {
        Init()
    }
    return node.Generate().String()
}

func GetBaseId() int64 {
    return baseId
}
