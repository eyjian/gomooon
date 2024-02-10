// Package utils
// Wrote by yijian on 2024/01/18
package utils

import (
    "os"
    "testing"
)

// go test -v -run="TestGetNonceStr"
func TestGetNonceStr(t *testing.T) {
    // 测试长度 1
    strLen := 1
    str := GetNonceStr(strLen)
    if len(str) == strLen {
        t.Logf("EXCEPTED: %s\n", str)
    } else {
        t.Errorf("UNEXCEPTED: %s\n", str)
    }

    // 测试长度 28
    strLen = 28
    str = GetNonceStr(28)
    if len(str) == strLen {
        t.Logf("EXCEPTED: %s\n", str)
    } else {
        t.Errorf("UNEXCEPTED: %s\n", str)
    }
}

// go test -v -run="TestDesensitizeStr"
func TestDesensitizeStr(t *testing.T) {
    str := "11204416541220243X"
    mask := DesensitizeStr(str, 2, 2)
    t.Logf("`%s` => %s\n", str, mask)

    mask = DesensitizeStr(str, 3, 2)
    t.Logf("`%s` => %s\n", str, mask)

    mask = DesensitizeStr(str, 2, 1)
    t.Logf("`%s` => %s\n", str, mask)

    mask = DesensitizeStr(str, 2, 0)
    t.Logf("`%s` => %s\n", str, mask)

    mask = DesensitizeStr(str, 0, 2)
    t.Logf("`%s` => %s\n", str, mask)

    str = "112044165412202"
    mask = DesensitizeStr(str, 2, 2)
    t.Logf("`%s` => %s\n", str, mask)
}

// go test -v -run="TestDesensitizeName"
func TestDesensitizeName(t *testing.T) {
    name := "张三"
    mask := DesensitizeName(name, 0, 1, 3)
    t.Logf("%s => %s\n", name, mask)

    name = "张三"
    mask = DesensitizeName(name, 1, 1, 3)
    t.Logf("%s => %s\n", name, mask)

    name = "王麻子"
    mask = DesensitizeName(name, 1, 1, 3)
    t.Logf("%s => %s\n", name, mask)

    name = "王麻子"
    mask = DesensitizeName(name, 0, 1, 3)
    t.Logf("%s => %s\n", name, mask)

    name = "欧阳大侠"
    mask = DesensitizeName(name, 1, 1, 3)
    t.Logf("%s => %s\n", name, mask)

    name = "欧阳大侠"
    mask = DesensitizeName(name, 0, 1, 3)
    t.Logf("%s => %s\n", name, mask)

    name = "张三.欧阳大侠"
    mask = DesensitizeName(name, 1, 1, 3)
    t.Logf("%s => %s\n", name, mask)

    name = "张三.欧阳大侠"
    mask = DesensitizeName(name, 0, 1, 3)
    t.Logf("%s => %s\n", name, mask)

    name = "张三"
    mask = DesensitizeName(name, 1, 0, 3)
    t.Logf("%s => %s\n", name, mask)

    name = "王麻子"
    mask = DesensitizeName(name, 1, 0, 3)
    t.Logf("%s => %s\n", name, mask)

    name = "欧阳大侠"
    mask = DesensitizeName(name, 1, 0, 3)
    t.Logf("%s => %s\n", name, mask)

    name = "张三.欧阳大侠"
    mask = DesensitizeName(name, 1, 0, 3)
    t.Logf("%s => %s\n", name, mask)
}

// 身份证号隐私数据，执行时指定
// go test -v -run="TestIsResidentIdentityCardNumber" -args FLAG ID // FLAG 只能取值 0 或者 1，1 表示 ID 为无效身份证号，0 表示为有效的身份证号
func TestIsResidentIdentityCardNumber(t *testing.T) {
    flag := os.Args[len(os.Args)-2]
    id := os.Args[len(os.Args)-1]
    if IsResidentIdentityCardNumber(id) {
        if flag == "0" {
            t.Errorf("%s is ID number\n", id)
        } else {
            t.Logf("%s is ID number\n", id)
        }
    } else {
        if flag == "0" {
            t.Logf("%s is not ID number\n", id)
        } else {
            t.Errorf("%s is not ID number\n", id)
        }
    }
}
