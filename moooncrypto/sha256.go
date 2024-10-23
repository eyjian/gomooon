// Package moooncrypto
// Wrote by yijian on 2024/08/26
package moooncrypto

import (
    "crypto/sha256"
    "encoding/hex"
)

func Sha256Sign(str, key string) string {
    b := sha256.Sum256([]byte(str + key))
    return hex.EncodeToString(b[:])
}
