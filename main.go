package main

import (
    "errors"
    "fmt"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "gorm.io/gorm/schema"
    "gorm.io/sharding"
    "gotest/util"
    "hash/fnv"
    "reflect"
    "sort"
    "sync"
    "time"
)

type Ipc struct {
    Uid         string         `json:"uid" gorm:"type:varchar(32);comment:uid"`
    DeviceType  string         `json:"deviceType" gorm:"type:varchar(4);comment:设备型号，字典:ipc_type"`
    FamilyId    int64          `json:"familyId" gorm:"type:bigint unsigned;comment:ipc挂靠的familyid"`
    OnlineState string         `json:"onlineState" gorm:"type:varchar(4);comment:在线状态，字典:ipc_online_state"`
    MediaState  string         `json:"mediaState" gorm:"type:varchar(4);comment:流媒体状态, 字典:ipc_media_state"`
    NatType     string         `json:"natType" gorm:"type:varchar(4);comment:所处网络的nat类型,字典：ipc_nat_type"`
    CreatedAt   time.Time      `json:"createdAt" gorm:"comment:创建时间"`
    UpdatedAt   time.Time      `json:"updatedAt" gorm:"comment:最后更新时间"`
    DeletedAt   gorm.DeletedAt `json:"deletedAt" gorm:"index;comment:删除时间"`
    CreateBy    int64          `json:"createBy" gorm:"index;comment:创建者"`
    UpdateBy    int64          `json:"updateBy" gorm:"index;comment:更新者"`
}

type Compare struct {
    Uid         string         `json:"uid" gorm:"type:varchar(32);comment:uid"`
    DeviceType  string         `json:"deviceType" gorm:"type:varchar(4);comment:设备型号，字典:ipc_type"`
    FamilyId    int64          `json:"familyId" gorm:"type:bigint unsigned;comment:ipc挂靠的familyid"`
    OnlineState string         `json:"onlineState" gorm:"type:varchar(4);comment:在线状态，字典:ipc_online_state"`
    MediaState  string         `json:"mediaState" gorm:"type:varchar(4);comment:流媒体状态, 字典:ipc_media_state"`
    NatType     string         `json:"natType" gorm:"type:varchar(4);comment:所处网络的nat类型,字典：ipc_nat_type"`
    CreatedAt   time.Time      `json:"createdAt" gorm:"comment:创建时间"`
    UpdatedAt   time.Time      `json:"updatedAt" gorm:"comment:最后更新时间"`
    DeletedAt   gorm.DeletedAt `json:"deletedAt" gorm:"index;comment:删除时间"`
    CreateBy    int64          `json:"createBy" gorm:"index;comment:创建者"`
    UpdateBy    int64          `json:"updateBy" gorm:"index;comment:更新者"`
}

func hash(s string) uint64 {
    h := fnv.New64a()
    h.Write([]byte(s))
    return h.Sum64()
}

func main() {
    dsn := "admin:shht@tcp(114.215.190.173:8000)/ipc?charset=utf8mb4&parseTime=True&loc=Local&timeout=1000ms"
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
        NamingStrategy: schema.NamingStrategy{
            SingularTable: true, // use singular table name, table for `User` would be `user` with this option enabled
        },
    })
    sqlDB, err := db.DB()
    sqlDB.SetMaxOpenConns(100)
    sqlDB.SetMaxIdleConns(10)

    middleware := sharding.Register(sharding.Config{
        ShardingKey:           "uid",
        NumberOfShards:        8,
        PrimaryKeyGenerator:   sharding.PKCustom,
        PrimaryKeyGeneratorFn: func(tableIdx int64) int64 { return 0 },
        ShardingAlgorithm: func(column any) (suffix string, err error) {
            if t, ok := column.(string); ok {
                uid := hash(t)
                return fmt.Sprintf("_%02d", uid%8), nil
            }
            return "", errors.New("invalid uid")
        },
    }, "ipc")
    db.Use(middleware)

    if err != nil {
        fmt.Println(err)
    }

    //createTable(db)
    //createData(db)
    //queryOne(db)
    //ipcs := queryAll(db)
    //transferData(db, ipcs)
    //queryPage(db)
    //queryPageMod(db)
    //testReflect()
    createOne(db)
}

func testReflect() {
    type IpcOrder struct {
        UidOrder         string `form:"uidOrder"  search:"type:order;column:uid;table:ipc"`
        DeviceTypeOrder  string `form:"deviceTypeOrder"  search:"type:order;column:device_type;table:ipc"`
        FamilyIdOrder    string `form:"familyIdOrder"  search:"type:order;column:family_id;table:ipc"`
        OnlineStateOrder string `form:"onlineStateOrder"  search:"type:order;column:online_state;table:ipc"`
        MediaStateOrder  string `form:"mediaStateOrder"  search:"type:order;column:media_state;table:ipc"`
        NatTypeOrder     string `form:"natTypeOrder"  search:"type:order;column:nat_type;table:ipc"`
        CreateByOrder    string `form:"createByOrder"  search:"type:order;column:create_by;table:ipc"`
        UpdateByOrder    string `form:"updateByOrder"  search:"type:order;column:update_by;table:ipc"`
        CreatedAtOrder   string `form:"createdAtOrder"  search:"type:order;column:created_at;table:ipc"`
        UpdatedAtOrder   string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:ipc"`
        DeletedAtOrder   string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:ipc"`
    }

    type IpcGetPageReq struct {
        TableSuffix string `search:"-"`

        Uid         string    `form:"uid"  search:"type:exact;column:uid;table:ipc" comment:""`
        DeviceType  string    `form:"deviceType"  search:"type:exact;column:device_type;table:ipc" comment:"设备型号，字典:ipc_type"`
        FamilyId    int64     `form:"familyId"  search:"type:exact;column:family_id;table:ipc" comment:"ipc挂靠的familyid"`
        OnlineState string    `form:"onlineState"  search:"type:exact;column:online_state;table:ipc" comment:"在线状态，字典:ipc_online_state"`
        MediaState  string    `form:"mediaState"  search:"type:exact;column:media_state;table:ipc" comment:"流媒体状态, 字典:ipc_media_state"`
        NatType     string    `form:"natType"  search:"type:exact;column:nat_type;table:ipc" comment:"所处网络的nat类型,字典：ipc_nat_type"`
        CreateBy    int64     `form:"createBy"  search:"type:exact;column:create_by;table:ipc" comment:"创建者"`
        UpdateBy    int64     `form:"updateBy"  search:"type:exact;column:update_by;table:ipc" comment:"更新者"`
        CreatedAt   time.Time `form:"createdAt"  search:"type:exact;column:created_at;table:ipc" comment:"创建时间"`
        UpdatedAt   time.Time `form:"updatedAt"  search:"type:exact;column:updated_at;table:ipc" comment:"最后更新时间"`
        DeletedAt   time.Time `form:"deletedAt"  search:"type:exact;column:deleted_at;table:ipc" comment:"删除时间"`
        IpcOrder
    }

    q := IpcGetPageReq{}
    q.IpcOrder.UidOrder = "asc"

    type OrderSeq struct {
        Key   string
        Value string
    }
    var seq []OrderSeq
    qType := reflect.TypeOf(q.IpcOrder)
    qValue := reflect.ValueOf(q.IpcOrder)
    for i := 0; i < qType.NumField(); i++ {
        field := qValue.Field(i)
        fieldType := qType.Field(i)
        if reflect.DeepEqual(field.Interface(), reflect.Zero(field.Type()).Interface()) {
            fmt.Println("empty")
        } else {
            s := OrderSeq{
                Key:   fieldType.Name[:len(fieldType.Name)-5],
                Value: field.String(),
            }
            seq = append(seq, s)
        }
    }

    fmt.Println(seq)
}

func transferData(db *gorm.DB, ipcs []Ipc) {
    var compares []Compare
    for _, v := range ipcs {
        c := Compare(v)
        compares = append(compares, c)
    }
    db.Create(compares)
}

func queryPage(db *gorm.DB) {
    var pageNo, pageSize int
    pageNo = 1
    pageSize = 15
    offset := pageNo * pageSize
    shardingSize := 8

    type queryRet struct {
        ipcs      []Ipc
        table_idx int
    }

    //第一次查询
    //var ipcs []Ipc
    var wg sync.WaitGroup
    retChan := make(chan queryRet, shardingSize)
    queryOffset := offset / shardingSize
    fmt.Println("queryOffset", queryOffset)
    f1 := func(idx int, rc chan queryRet, wg *sync.WaitGroup) {
        defer wg.Done()
        var ipc_split []Ipc
        sql := fmt.Sprintf("select * from ipc_%02d limit %d offset %d", idx, pageSize, queryOffset)
        db.Raw(sql).Scan(&ipc_split)
        ret := queryRet{
            ipcs:      ipc_split,
            table_idx: idx,
        }
        rc <- ret
    }

    for i := 0; i < shardingSize; i++ {
        wg.Add(1)
        go f1(i, retChan, &wg)
    }
    wg.Wait()

    //取最小值
    uidMin := "0"
    splitMax := make(map[int]string, shardingSize) //每个分片的最大值
    splitMin := make(map[int]string, shardingSize) //每个分片的最小值
    for i := 0; i < shardingSize; i++ {
        queryRet := <-retChan
        fmt.Println("len is ", len(queryRet.ipcs), queryRet.table_idx)
        if len(queryRet.ipcs) > 0 {
            splitMax[queryRet.table_idx] = queryRet.ipcs[len(queryRet.ipcs)-1].Uid
            splitMin[queryRet.table_idx] = queryRet.ipcs[0].Uid
            splitMin := queryRet.ipcs[0].Uid
            if splitMin < uidMin || uidMin == "0" {
                uidMin = splitMin
            }
        } else {
            splitMax[queryRet.table_idx] = "0"
            splitMin[queryRet.table_idx] = "0"
        }

    }
    fmt.Println(uidMin)

    //第二次查询
    f2 := func(idx int, rc chan queryRet, wg *sync.WaitGroup) {
        defer wg.Done()
        var ipc_split []Ipc
        uidMax := "Z"
        if splitMax[idx] != "0" {
            uidMax = splitMax[idx]
        }
        sql := fmt.Sprintf("select * from ipc_%02d where uid between '%s' and '%s'", idx, uidMin, uidMax)
        fmt.Println(sql)
        db.Raw(sql).Scan(&ipc_split)
        ret := queryRet{
            ipcs:      ipc_split,
            table_idx: idx,
        }
        rc <- ret
    }

    for i := 0; i < shardingSize; i++ {
        wg.Add(1)
        go f2(i, retChan, &wg)
    }
    wg.Wait()

    //计算uidMin的全局offset(每个分片的offset之和)
    uidMinOffset := 0
    //splitOffset := make(map[int]int, 8) //每个分片的offset
    var ipcs []Ipc
    for i := 0; i < shardingSize; i++ {
        queryRet := <-retChan
        table_idx := queryRet.table_idx
        //跟第一次的返回结果相比，有多少个比第一次的最小值还小的
        offset := queryOffset + 1
        for _, v := range queryRet.ipcs {
            uid := v.Uid

            if uid == splitMin[table_idx] {
                break
            }

            offset -= 1
        }
        if queryRet.ipcs[0].Uid != uidMin {
            offset -= 1
        }
        fmt.Println("offset:", offset, table_idx)
        uidMinOffset += offset
        ipcs = append(ipcs, queryRet.ipcs...)
    }
    fmt.Println(uidMinOffset)

    sort.Slice(ipcs, func(i, j int) bool {
        return ipcs[i].Uid < ipcs[j].Uid
    })

    beginIdx := offset - uidMinOffset + 1
    result := ipcs[beginIdx : beginIdx+pageSize]

    for _, v := range result {
        fmt.Println(v.Uid)
    }

}

func queryPageMod(db *gorm.DB) {
    var pageNo, pageSize int
    pageNo = 1
    pageSize = 15
    offset := pageNo * pageSize
    shardingSize := 8

    type queryRet struct {
        ipcs      []Ipc
        table_idx int
    }

    //第一次查询
    //var ipcs []Ipc
    var wg sync.WaitGroup
    retChan := make(chan queryRet, shardingSize)
    queryOffset := offset / shardingSize
    fmt.Println("queryOffset", queryOffset)
    f1 := func(idx int, rc chan queryRet, wg *sync.WaitGroup) {
        defer wg.Done()
        var ipc_split []Ipc
        tableName := fmt.Sprintf("ipc_%02d", idx)
        db.Table(tableName).Limit(pageSize).Offset(queryOffset).Order("order by created_at desc").Find(&ipc_split)

        ret := queryRet{
            ipcs:      ipc_split,
            table_idx: idx,
        }
        rc <- ret
    }

    for i := 0; i < shardingSize; i++ {
        wg.Add(1)
        go f1(i, retChan, &wg)
    }
    wg.Wait()

    //取最小值
    uidMin := "0"
    splitMax := make(map[int]string, shardingSize) //每个分片的最大值
    splitMin := make(map[int]string, shardingSize) //每个分片的最小值
    for i := 0; i < shardingSize; i++ {
        queryRet := <-retChan
        fmt.Println("len is ", len(queryRet.ipcs), queryRet.table_idx)
        if len(queryRet.ipcs) > 0 {
            splitMax[queryRet.table_idx] = queryRet.ipcs[len(queryRet.ipcs)-1].Uid
            splitMin[queryRet.table_idx] = queryRet.ipcs[0].Uid
            splitMin := queryRet.ipcs[0].Uid
            if splitMin < uidMin || uidMin == "0" {
                uidMin = splitMin
            }
        } else {
            splitMax[queryRet.table_idx] = "0"
            splitMin[queryRet.table_idx] = "0"
        }

    }
    fmt.Println(uidMin)

    //第二次查询
    f2 := func(idx int, rc chan queryRet, wg *sync.WaitGroup) {
        defer wg.Done()
        var ipc_split []Ipc
        tableName := fmt.Sprintf("ipc_%02d", idx)
        s := db.Table(tableName).Where("uid >= ?", uidMin)
        if splitMax[idx] != "0" {
            s = s.Where("uid <= ?", splitMax[idx])
        }
        s.Find(&ipc_split)

        ret := queryRet{
            ipcs:      ipc_split,
            table_idx: idx,
        }
        rc <- ret
    }

    for i := 0; i < shardingSize; i++ {
        wg.Add(1)
        go f2(i, retChan, &wg)
    }
    wg.Wait()

    //计算uidMin的全局offset(每个分片的offset之和)
    uidMinOffset := 0
    //splitOffset := make(map[int]int, 8) //每个分片的offset
    var ipcs []Ipc
    for i := 0; i < shardingSize; i++ {
        queryRet := <-retChan
        table_idx := queryRet.table_idx
        //跟第一次的返回结果相比，有多少个比第一次的最小值还小的
        offset := queryOffset + 1
        for _, v := range queryRet.ipcs {
            uid := v.Uid

            if uid == splitMin[table_idx] {
                break
            }

            offset -= 1
        }
        if queryRet.ipcs[0].Uid != uidMin {
            offset -= 1
        }
        fmt.Println("offset:", offset, table_idx)
        uidMinOffset += offset
        ipcs = append(ipcs, queryRet.ipcs...)
    }
    fmt.Println(uidMinOffset)

    sort.Slice(ipcs, func(i, j int) bool {
        return ipcs[i].Uid < ipcs[j].Uid
    })

    beginIdx := offset - uidMinOffset + 1
    result := ipcs[beginIdx : beginIdx+pageSize]

    for _, v := range result {
        fmt.Println(v.Uid)
    }

}

func createTable(db *gorm.DB) {
    for i := 0; i < 8; i += 1 {
        table := fmt.Sprintf("ipc_%02d", i)
        db.Exec(`DROP TABLE IF EXISTS ` + table)
        db.Exec(`CREATE TABLE ` + table + "(uid VARCHAR(64) NOT NULL,device_type VARCHAR(4) NOT NULL DEFAULT '0' COMMENT '设备型号，字典:ipc_type',family_id BIGINT UNSIGNED NOT NULL DEFAULT '0' COMMENT 'ipc挂靠的familyid',online_state VARCHAR(4) NOT NULL DEFAULT '0' COMMENT '在线状态，字典:ipc_online_state',media_state VARCHAR(4) NOT NULL DEFAULT '0' COMMENT '流媒体状态, 字典:ipc_media_state',nat_type VARCHAR(4) NOT NULL DEFAULT '0' COMMENT '所处网络的nat类型,字典：ipc_nat_type',create_by BIGINT DEFAULT NULL COMMENT '创建者',update_by BIGINT DEFAULT NULL COMMENT '更新者',created_at DATETIME(3) DEFAULT NULL COMMENT '创建时间',updated_at DATETIME(3) DEFAULT NULL COMMENT '最后更新时间',deleted_at DATETIME(3) DEFAULT NULL COMMENT '删除时间',PRIMARY KEY (`uid`));")
    }
}

func createData(db *gorm.DB) {
    for i := 0; i < 10; i++ {
        uid := util.GenerateIdString()
        db.Create(&Ipc{Uid: uid, DeviceType: "A9"})
    }
}

func createOne(db *gorm.DB) {
    err := db.Create(&Ipc{Uid: "841106208624896", DeviceType: "A9"}).Error
    fmt.Println(err)
}

func queryOne(db *gorm.DB) {
    //var ipcs []Ipc
    //db.Raw("SELECT * FROM ipc_04 WHERE uid = ?", "841086274427136").Scan(&ipcs)
    //var data Ipc
    var model Ipc
    //db.Model(&Ipc{}).Where("uid", "841087966630144").Unscoped().First(&model)
    db.Model(&Ipc{}).Unscoped().First(&model, "841087966630144")
    fmt.Printf("%#v\n", model)
}

func queryAll(db *gorm.DB) []Ipc {
    var ipcs []Ipc
    var wg sync.WaitGroup
    retChan := make(chan []Ipc, 8)
    f := func(idx int, rc chan []Ipc, wg *sync.WaitGroup) {
        defer wg.Done()
        var ipc_split []Ipc
        table_name := fmt.Sprintf("ipc_%02d", idx)
        db.Raw("SELECT * FROM " + table_name).Scan(&ipc_split)
        rc <- ipc_split
    }

    for i := 0; i < 8; i++ {
        wg.Add(1)
        go f(i, retChan, &wg)
    }

    wg.Wait()
    for i := 0; i < 8; i++ {
        ipc_split := <-retChan
        ipcs = append(ipcs, ipc_split...)
    }
    fmt.Println(len(ipcs))
    return ipcs
}
