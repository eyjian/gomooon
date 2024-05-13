// Package moooncrypto
// Wrote by yijian on 2024/01/02
package moooncrypto

import (
    "crypto/x509"
    "encoding/pem"
    "os"
    "testing"
)

// go test -v -run="TestRsaSha256SignWithPrivateKey"
// 在线工具：https://8gwifi.org/rsasignverifyfunctions.jsp
func TestRsaSha256SignWithPrivateKey(t *testing.T) {
    // 读取私钥文件
    privateKeyBytes, err := os.ReadFile("./id_rsa_256") // 256 * 8 = 2048 位
    if err != nil {
        t.Errorf("error reading private key file: %s\n", err.Error())
        return
    }

    // 解析私钥
    block, _ := pem.Decode(privateKeyBytes)
    if block == nil {
        t.Error("failed to decode private key")
        return
    }

    privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
    if err != nil {
        t.Errorf("Error parsing private key: %s\n", err.Error())
        return
    }

    // 待签名数据
    data := []byte("data to be signed")

    // 调用 Rsa256SignWithPrivateKey 函数进行签名
    signature, err := RsaSha256SignWithPrivateKey(privateKey, data)
    if err != nil {
        t.Errorf("Error signing data: %s\n", err.Error())
        return
    }

    t.Logf("data: %s\nsignature: %s\n", string(data), signature)
}

// go test -v -run="TestRsaSha256SignWithPrivateKeyStr"
func TestRsaSha256SignWithPrivateKeyStr(t *testing.T) {
    // 读取私钥文件
    privateKeyBytes, err := os.ReadFile("./id_rsa_256")
    if err != nil {
        t.Errorf("error reading private key file: %s\n", err.Error())
        return
    }

    // 解析私钥
    block, _ := pem.Decode(privateKeyBytes)
    if block == nil {
        t.Error("failed to decode private key")
        return
    }

    // 待签名数据
    data := []byte("data to be signed")

    // 调用 Rsa256SignWithPrivateKeyStr 函数进行签名
    signature, err := RsaSha256SignWithPrivateKeyStr(privateKeyBytes, data)
    if err != nil {
        t.Errorf("Error signing data: %s\n", err.Error())
        return
    }

    t.Logf("data: %s\nsignature: %s\n", string(data), signature)
}
