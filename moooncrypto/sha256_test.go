// Package moooncrypto
// Wrote by yijian on 2024/08/26
package moooncrypto

import (
    "testing"
)

// go test -v -run="TestSha256Sign$"
func TestSha256Sign(t *testing.T) {
    signature := Sha256Sign("123456", "123456")
    t.Logf(signature)
}
