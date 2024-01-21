// Package utils
// Wrote by yijian on 2024/01/18
package utils

import (
	"crypto/x509"
	"encoding/pem"
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

// 前面的密钥需要命令行传入
// go test -v -run="TestUpperHmacSHA256Sign" -args KEY DATA
func TestUpperHmacSHA256Sign(t *testing.T) {
	key := os.Args[len(os.Args)-2]
	data := os.Args[len(os.Args)-1]
	signature, err := UpperHmacSHA256Sign(data, key)
	if err != nil {
		t.Errorf("%s\n", err.Error())
	} else {
		t.Logf("signature: %s\n", signature)
	}
}

// 前面的密钥需要命令行传入
// go test -v -run="TestLowerHmacSHA256Sign" -args KEY DATA
func TestLowerHmacSHA256Sign(t *testing.T) {
	key := os.Args[len(os.Args)-2]
	data := os.Args[len(os.Args)-1]
	signature, err := LowerHmacSHA256Sign(data, key)
	if err != nil {
		t.Errorf("%s\n", err.Error())
	} else {
		t.Logf("signature: %s\n", signature)
	}
}

// go test -v -run="TestRsa256SignWithPrivateKey"
func TestRsa256SignWithPrivateKey(t *testing.T) {
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

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		t.Errorf("Error parsing private key: %s\n", err.Error())
		return
	}

	// 待签名数据
	data := []byte("data to be signed")

	// 调用 Rsa256SignWithPrivateKey 函数进行签名
	signature, err := Rsa256SignWithPrivateKey(privateKey, data)
	if err != nil {
		t.Errorf("Error signing data: %s\n", err.Error())
		return
	}

	t.Logf("Signature: %s\n", string(signature))
}

// go test -v -run="TestRsa256SignWithPrivateKeyStr"
func TestRsa256SignWithPrivateKeyStr(t *testing.T) {
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
	data := []byte("" +
		"")

	// 调用 Rsa256SignWithPrivateKeyStr 函数进行签名
	signature, err := Rsa256SignWithPrivateKeyStr(privateKeyBytes, data)
	if err != nil {
		t.Errorf("Error signing data: %s\n", err.Error())
		return
	}

	t.Logf("Signature: %s\n", string(signature))
}
