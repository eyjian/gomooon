// Package moooncrypto
// Wrote by yijian on 2024/01/02
package moooncrypto

import (
    "testing"
)

// 前面的密钥需要命令行传入
// go test -v -run="TestHmacSha256Sign" #-args KEY DATA
func TestHmacSha256Sign(t *testing.T) {
    key := "192006250b4c09247ec02edce69f6a2d"
    data := "appid=wxd930ea5d5a258f4f&body=test&device_info=1000&mch_id=10000100&nonce_str=ibuaiVcKdpRxkhJA&key=" + key
    signature, err := HmacSha256Sign(data, key, true)
    if err != nil {
        t.Errorf("%s\n", err.Error())
    } else {
        t.Logf("signature: %s\n", signature)
        excepted := "6A9AE1657590FD6257D693A078E1C3E4BB6BA4DC30B23E0EE2496E54170DACD6"
        if signature != excepted {
            t.Errorf("sign error: %s, excepted: %s\n", signature, excepted)
        } else {
            t.Logf("sign ok\n")
        }
    }
}

// go test -v -run="TestHmacSha256"
func TestHmacSha256(t *testing.T) {
    key := "192006250b4c09247ec02edce69f6a2d"
    str := "123456"
    signature := HmacSha256(str, key)
    t.Logf("signature: %s\n", signature)
}
