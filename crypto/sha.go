// Package crypto
// Wrote by yijian on 2024/01/02
package crypto

import (
	"crypto"
	"crypto/hmac"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"strings"
)

// 公钥加密、私钥解密、私钥签名、公钥验签
// 身份认证（信息劫持）：防冒充、防伪造（非对称加密），确保正确的用户访问正确的网站
// 数据加密（信息加密）：防窃听（对称加密） ，使第三方无法查看明文内容
// 数据一致性、完整性效验（信息篡改）：防篡改、防抵赖（哈希算法），内容被篡改能及时发现
// 证书的作用：效验签名、传递公钥（客户端加密数据用）
// 数字签名：私钥加密（服务端），公钥解密（客户端）
// 非对称加密：密钥配送（交换）

// UpperHmacSHA256Sign SHA256 签名
// data 需要签名的数据
// key 签名密钥
func UpperHmacSHA256Sign(data, key string) (string, error) {
	return hmacSHA256Sign(data, key, true)
}

func LowerHmacSHA256Sign(data, key string) (string, error) {
	return hmacSHA256Sign(data, key, false)
}

// RSA（Rivest-Shamir-Adleman）是一种非对称加密算法，由 Ron Rivest、Adi Shamir 和 Leonard Adleman 于 1977 年提出
// HMAC: 哈希运算消息认证码（Hash-based Message Authentication Code），带密钥的 Hash 函数。
// 由 H.Krawezyk、M.Bellare 和 R.Canetti 于 1996 年提出的一种基于 Hash 函数和密钥进行消息认证的方法，并于 1997 年作为 RFC2104 被公布。
func hmacSHA256Sign(data, key string, toUpper bool) (string, error) {
	hash := hmac.New(sha256.New, []byte(key))
	_, err := hash.Write([]byte(data))
	if err != nil {
		return "", err
	} else {
		if toUpper {
			return strings.ToUpper(hex.EncodeToString(hash.Sum(nil))), nil
		} else {
			return strings.ToLower(hex.EncodeToString(hash.Sum(nil))), nil
		}
	}
}

func SHA256RSASignWithPrivateKey(privateKey *rsa.PrivateKey, data []byte) (string, error) {
	// 签名数据
	hash := sha256.Sum256(data)
	signature, err := rsa.SignPKCS1v15(nil, privateKey, crypto.SHA256, hash[:])
	if err != nil {
		return "", fmt.Errorf("failed to sign data: %v", err)
	}

	// 将签名结果转换为 Base64 编码
	encodedSignature := make([]byte, base64.StdEncoding.EncodedLen(len(signature)))
	base64.StdEncoding.Encode(encodedSignature, signature)

	return string(encodedSignature), nil
}

func SHA256RSASignWithPrivateKeyStr(privateKey []byte, data []byte) (string, error) {
	// 解析私钥
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return "", errors.New("failed to decode private key")
	}

	privateRsa, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("failed to parse private key: %v", err)
	}

	// 签名数据
	hash := sha256.Sum256(data)
	signature, err := rsa.SignPKCS1v15(nil, privateRsa, crypto.SHA256, hash[:])
	if err != nil {
		return "", fmt.Errorf("failed to sign data: %v", err)
	}

	// 将签名结果转换为 Base64 编码
	encodedSignature := make([]byte, base64.StdEncoding.EncodedLen(len(signature)))
	base64.StdEncoding.Encode(encodedSignature, signature)
	return string(encodedSignature), nil
}
