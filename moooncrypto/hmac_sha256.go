// Package moooncrypto
// Wrote by yijian on 2024/01/02
package moooncrypto

import (
    "crypto/hmac"
    "crypto/sha256"
    "encoding/hex"
    "strings"
)

// HMAC: 哈希运算消息认证码（Hash-based Message Authentication Code），带密钥的 Hash 函数。
// 由 H.Krawezyk、M.Bellare 和 R.Canetti 于 1996 年提出的一种基于 Hash 函数和密钥进行消息认证的方法，并于 1997 年作为 RFC2104 被公布。

// HmacSha256Sign SHA256 签名
// data 需要签名的数据
// key 签名密钥
// toUpper 为 true 返回大写的签名字符串，为 false 返回小写的签名字符串
func HmacSha256Sign(data, key string, toUpper bool) (string, error) {
    hash := hmac.New(sha256.New, []byte(key))
    _, err := hash.Write([]byte(data))
    if err != nil {
        return "", err
    }

    // hex.EncodeToString 返回的是小写的十六进制字符串
    if toUpper {
        return strings.ToUpper(hex.EncodeToString(hash.Sum(nil))), nil
    } else {
        return strings.ToLower(hex.EncodeToString(hash.Sum(nil))), nil
    }
}

// HmacSha256 HmacSha256 签名
func HmacSha256(data, key string) string {
    hashed := hmac.New(sha256.New, []byte(key))
    hashed.Write([]byte(data))
    return string(hashed.Sum(nil))
}
