// Package moooncrypto
// Wrote by yijian on 2024/02/02
package moooncrypto

import (
    "os"
    "testing"
)

// go test -v -run="TestMd5Sum" -args DATA
func TestMd5Sum(t *testing.T) {
    data := os.Args[len(os.Args)-1]
    md5str := Md5Sum(data, true)
    t.Logf("%s\n", md5str)
}
