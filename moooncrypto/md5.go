// Package moooncrypto
// Wrote by yijian on 2024/01/02
package moooncrypto

import (
    "crypto/md5"
    "encoding/hex"
    "strings"
)

// Md5Sum MD5 计算
// data 需要计算的数据
func Md5Sum(data string, toUpper bool) string {
    hash := md5.Sum([]byte(data))
    if toUpper {
        return strings.ToUpper(hex.EncodeToString(hash[:]))
    } else {
        return strings.ToLower(hex.EncodeToString(hash[:]))
    }
}
