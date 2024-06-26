// Package mooongorm
// Wrote by yijian on 2024/05/24
package mooongorm

import (
    "fmt"
    "github.com/eyjian/gomooon/mooonstr"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
    "reflect"
    "testing"
    "time"
)

type TableStruct struct {
    Id         uint32    `gorm:"column:f_id;primaryKey;autoIncrement"`
    Name       string    `gorm:"column:f_name"`
    CreateTime time.Time `gorm:"column:f_create_time"`
    UpdateTIme time.Time `gorm:"column:f_update_time"`
    Memo       string    `gorm:"column:f_memo"`
}

// go test -v -run="TestFilterFields$"
func TestFilterFields(t *testing.T) {
    typeOfModel := reflect.TypeOf((*TableStruct)(nil)).Elem()

    filteredFields := mooonstr.NewStringSet()
    filteredFields.Add("UpdateTIme")

    unfilteredFields := FilterFields(typeOfModel, &filteredFields)
    t.Logf("%+v\n", unfilteredFields)

    if unfilteredFields[0] != "f_id" {
        t.Errorf("f_id error: %s\n", unfilteredFields[0])
        return
    }
    if unfilteredFields[1] != "f_name" {
        t.Errorf("f_id error: %s\n", unfilteredFields[1])
        return
    }
    if unfilteredFields[2] != "f_create_time" {
        t.Errorf("f_id error: %s\n", unfilteredFields[2])
        return
    }
    if unfilteredFields[3] != "f_memo" {
        t.Errorf("f_id error: %s\n", unfilteredFields[3])
        return
    }
}

type TestModel struct {
    ID   int    `gorm:"primaryKey"`
    Name string `gorm:"unique"`
    Age  int
}

// go test -v -run="TestCreateInBatchesWithoutTransaction$"
func TestCreateInBatchesWithoutTransaction(t *testing.T) {
    // 初始化数据库连接
    db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Silent),
    })
    if err != nil {
        t.Fatalf("failed to connect to the database: %v", err)
    }

    // 自动迁移数据表
    err = db.AutoMigrate(&TestModel{})
    if err != nil {
        t.Fatalf("failed to migrate the schema: %v", err)
    }

    // 准备数据
    data := make([]TestModel, 10)
    for i := 0; i < 10; i++ {
        data[i] = TestModel{
            Name: fmt.Sprintf("test-%d", i),
        }
    }

    db = db.Debug()
    tx := db.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()

    // 调用 CreateInBatchesWithoutTransaction 函数
    tx = tx.Omit("Id", "Age")
    err = CreateInBatchesWithoutTransaction(tx, data, 2).Error
    if err != nil {
        t.Fatalf("failed to create records in batches: %v", err)
    }

    // 验证数据是否正确插入
    var count int64
    err = tx.Model(&TestModel{}).Count(&count).Error
    if err != nil {
        t.Fatalf("failed to count records: %v", err)
    } else {
        t.Logf("count:%d\n", count)
    }

    if count != int64(len(data)) {
        t.Fatalf("expected %d records, but got %d", len(data), count)
    }

    // 验证插入的数据是否正确
    var results []TestModel
    err = tx.Find(&results).Error
    if err != nil {
        t.Fatalf("failed to find records: %v", err)
    } else {
        t.Logf("%+v\n", results)
    }

    for i, result := range results {
        if result.Name != fmt.Sprintf("test-%d", i) {
            t.Fatalf("expected name to be %q, but got %q", fmt.Sprintf("test-%d", i), result.Name)
        } else {
            t.Logf("[%d] %+v\n", i, result)
        }
    }

    t.Log("test ok")
}

// go test -v -run="TestCreateInBatchesWithoutTransactionErr$"
func TestCreateInBatchesWithoutTransactionErr(t *testing.T) {
    // 初始化数据库连接
    db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Silent),
    })
    if err != nil {
        t.Fatalf("failed to connect to the database: %v", err)
    }

    // 自动迁移数据表
    err = db.AutoMigrate(&TestModel{})
    if err != nil {
        t.Fatalf("failed to migrate the schema: %v", err)
    }

    // 准备数据
    data := make([]TestModel, 10)
    for i := 0; i < 10; i++ {
        data[i] = TestModel{
            Name: fmt.Sprintf("test-%d", i),
        }
    }

    // 创建一个包含错误的数据
    dataWithError := []TestModel{
        {Name: "duplicate"},
        {Name: "duplicate"},
    }

    tx := db.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()

    // 调用 CreateInBatchesWithoutTransaction 函数
    err = CreateInBatchesWithoutTransaction(tx, dataWithError, 2).Error
    if err != nil {
        t.Logf("expected error: %v", err)
        tx.Rollback()
    } else {
        t.Fatal("expected an error, but got none")
    }

    // 验证数据是否已回滚
    var count int64
    err = db.Model(&TestModel{}).Count(&count).Error
    if err != nil {
        t.Fatalf("failed to count records: %v", err)
    }

    if count != 0 {
        t.Fatalf("expected 0 records after rollback, but got %d", count)
    } else {
        t.Logf("count is %d\n", count)
    }

    t.Log("test ok")
}
