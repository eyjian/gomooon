// Package mooongorm
// Wrote by yijian on 2024/05/24
package mooongorm

import (
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