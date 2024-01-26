// Package crypto
// Wrote by yijian on 2024/01/02
package crypto

import (
	"crypto/x509"
	"encoding/pem"
	"os"
	"testing"
)

// 前面的密钥需要命令行传入
// go test -v -run="TestUpperHmacSHA256Sign" #-args KEY DATA
func TestUpperHmacSHA256Sign(t *testing.T) {
	key := "192006250b4c09247ec02edce69f6a2d"
	data := "appid=wxd930ea5d5a258f4f&body=test&device_info=1000&mch_id=10000100&nonce_str=ibuaiVcKdpRxkhJA&key=" + key
	signature, err := UpperHmacSHA256Sign(data, key)
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

// go test -v -run="TestSHA256RSASignWithPrivateKey"
// 在线工具：https://8gwifi.org/rsasignverifyfunctions.jsp
func TestSHA256RSASignWithPrivateKey(t *testing.T) {
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
	signature, err := SHA256RSASignWithPrivateKey(privateKey, data)
	if err != nil {
		t.Errorf("Error signing data: %s\n", err.Error())
		return
	}

	t.Logf("data: %s\nsignature: %s\n", string(data), signature)
}

// go test -v -run="TestSHA256RSASignWithPrivateKeyStr"
func TestSHA256RSASignWithPrivateKeyStr(t *testing.T) {
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
	signature, err := SHA256RSASignWithPrivateKeyStr(privateKeyBytes, data)
	if err != nil {
		t.Errorf("Error signing data: %s\n", err.Error())
		return
	}

	t.Logf("data: %s\nsignature: %s\n", string(data), signature)
}
