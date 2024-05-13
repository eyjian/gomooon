// Package moooncrypto
// Wrote by yijian on 2024/01/02
package moooncrypto

import (
    "crypto"
    "crypto/rsa"
    "crypto/sha256"
    "crypto/x509"
    "encoding/base64"
    "encoding/pem"
    "errors"
    "fmt"
)

// RSA（Rivest-Shamir-Adleman）是一种非对称加密算法，由 Ron Rivest、Adi Shamir 和 Leonard Adleman 于 1977 年提出

func RsaSha256SignWithPrivateKey(privateKey *rsa.PrivateKey, data []byte) (string, error) {
    // 签名数据
    hash := sha256.Sum256(data)
    signature, err := rsa.SignPKCS1v15(nil, privateKey, crypto.SHA256, hash[:])
    if err != nil {
        return "", fmt.Errorf("RSA-SHA256 sign error: %v", err)
    }

    // 将签名结果转换为 Base64 编码
    encodedSignature := make([]byte, base64.StdEncoding.EncodedLen(len(signature)))
    base64.StdEncoding.Encode(encodedSignature, signature)

    return string(encodedSignature), nil
}

func RsaSha256SignWithPrivateKeyStr(privateKeyStr []byte, data []byte) (string, error) {
    // 解析私钥
    block, _ := pem.Decode(privateKeyStr)
    if block == nil {
        return "", errors.New("decode private key error")
    }

    privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
    if err != nil {
        return "", fmt.Errorf("parse private key error: %s", err.Error())
    }

    return RsaSha256SignWithPrivateKey(privateKey, data)
}
