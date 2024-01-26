// Package utils
// Wrote by yijian on 2024/01/02
package utils

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

// Md5Sign MD5 签名
// data 需要签名的数据
func Md5Sign(data) string {
	hash := md5.Sum([]byte(data))
	return strings.ToUpper(hex.EncodeToString(hash[:]))
}
