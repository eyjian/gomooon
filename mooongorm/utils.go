// Package mooongorm
// Wrote by yijian on 2024/05/24
package mooongorm

import (
    "errors"
    "gorm.io/gorm"
    "reflect"
    "strings"
)
import (
    "github.com/eyjian/gomooon/mooonstr"
)

// FilterFields 过滤掉指定的字段
// 应用场景：INSERT 或 UPDATE 操作时排除某些字段
// typeOfModel 通过表的 Model 而来，如：typeOfModel := reflect.TypeOf((*TableStruct)(nil)).Elem()
// 注意表的 Model 定义一定要有 gorm 的 tag，如：`gorm:"column:f_name"`，要求 column 后跟的是字段名
func FilterFields(typeOfModel reflect.Type, filteredFields *mooonstr.StringSet) []string {
    numFields := typeOfModel.NumField()              // 取得总的字段数
    unfilteredFields := make([]string, 0, numFields) // 存放所有未被过滤掉的字段

    // 遍历所有字段，排除不需要更新的字段
    for i := 0; i < numFields; i++ {
        field := typeOfModel.Field(i)
        if !filteredFields.Contains(field.Name) { // field.Name 值为 struct 的字段名，如：Id、Name、CreateTime
            // 使用结构体标签（tag）获取字段的列名
            fieldName := field.Tag.Get("gorm") // 这里的 fieldName 值示例：column:f_id;primaryKey;autoIncrement
            parts := strings.Split(fieldName, ";")
            for _, part := range parts {
                if strings.HasPrefix(part, "column:") {
                    fieldName = strings.TrimPrefix(part, "column:")
                    break
                }
            }
            unfilteredFields = append(unfilteredFields, fieldName)
        }
    }

    return unfilteredFields
}

// CreateInBatchesWithoutTransaction inserts value in batches of batchSize without transaction
// 使用注意：
// 1、应当总是将 CreateInBatchesWithoutTransaction 放在一个事务中，如果返回错误调用者应当执行回滚操作；
// 2、调用的 db 的 db.CreateBatchSize 值应当为 0，原因时 db.Create 会调用 db.CreateInBatches，而 db.CreateInBatches 会开启子事务，导致出现嵌套事务。
func CreateInBatchesWithoutTransaction(db *gorm.DB, value interface{}, batchSize int) *gorm.DB {
    // 检查参数是否有效
    if batchSize <= 0 {
        db.AddError(errors.New("batch size must be a positive integer"))
        return db
    }

    reflectValue := reflect.Indirect(reflect.ValueOf(value))
    switch reflectValue.Kind() {
    case reflect.Slice, reflect.Array:
        var rowsAffected int64

        // 获取切片的长度
        reflectLen := reflectValue.Len()
        if reflectLen == 0 {
            db.AddError(errors.New("the slice is empty"))
            return db
        }

        // 按批次插入数据
        for i := 0; i < reflectLen; i += batchSize {
            end := i + batchSize
            if end > reflectLen {
                end = reflectLen
            }

            // 插入当前批次的数据
            err := db.Create(reflectValue.Slice(i, end).Interface()).Error
            if err != nil {
                db.AddError(err)
                return db // 一旦出错就返回，调用者应当执行回滚操作
            }
            rowsAffected += db.RowsAffected
            db.RowsAffected = rowsAffected
        }

        return db
    default:
        return db.Create(value)
    }
}
