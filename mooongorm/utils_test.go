// Package mooongorm
// Wrote by yijian on 2024/05/24
package mooongorm

import (
    "github.com/eyjian/gomooon/mooonstr"
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