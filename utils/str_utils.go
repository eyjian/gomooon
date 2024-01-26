// Package utils
// Wrote by yijian on 2024/01/02
package utils

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
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

const allCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const hexCharset = "abcdefABCDEF0123456789"

var (
	r      *rand.Rand
	mu     sync.Mutex
	once   sync.Once
	buffer sync.Pool
)

func init() {
	source := rand.NewSource(time.Now().UnixNano())
	r = rand.New(source)
	buffer.New = func() interface{} {
		return make([]byte, 0, 64)
	}
}

func getNonceStr(length int, charset string) string {
	once.Do(func() {
		mu.Lock()
		defer mu.Unlock()
		r.Seed(time.Now().UnixNano())
	})

	mu.Lock()
	defer mu.Unlock()

	// Get a buffer from the pool and reset its length to the desired value
	buf := buffer.Get().([]byte)[:length]

	for i := range buf {
		buf[i] = charset[r.Intn(len(charset))]
	}

	// Convert the buffer to a string, put it back into the pool, and return the result
	result := string(buf)
	buffer.Put(buf)
	return result
}

func GetNonceStr(length int) string {
	return getNonceStr(length, allCharset)
}

func GetHexNonceStr(length int) string {
	return getNonceStr(length, hexCharset)
}

// IsResidentIdentityCardNumber 判断是否为身份证号
func IsResidentIdentityCardNumber(id string) bool {
	// 身份证号码的正则表达式
	pattern15 := `^\d{15}$`
	pattern18 := `^\d{17}(\d|X|x)$`
	reg15 := regexp.MustCompile(pattern15)
	reg18 := regexp.MustCompile(pattern18)
	is15 := reg15.MatchString(id)
	is18 := reg18.MatchString(id)

	if !is15 && !is18 {
		return false
	}

	// 提取出生日期
	var birthYear, birthMonth, birthDay int
	if is15 {
		birthYear, _ = strconv.Atoi("19" + id[6:8])
		birthMonth, _ = strconv.Atoi(id[8:10])
		birthDay, _ = strconv.Atoi(id[10:12])
	} else {
		birthYear, _ = strconv.Atoi(id[6:10])
		birthMonth, _ = strconv.Atoi(id[10:12])
		birthDay, _ = strconv.Atoi(id[12:14])
	}
	birthDate := time.Date(birthYear, time.Month(birthMonth), birthDay, 0, 0, 0, 0, time.UTC)
	if birthDate.Year() != birthYear || birthDate.Month() != time.Month(birthMonth) || birthDate.Day() != birthDay {
		return false
	}

	// 如果是18位身份证，需要检查校验码
	if is18 {
		weights := []int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
		checkSum := 0
		for i := 0; i < 17; i++ {
			num, _ := strconv.Atoi(string(id[i]))
			checkSum += num * weights[i]
		}
		checkCodes := "10X98765432"
		return string(checkCodes[checkSum%11]) == id[17:]
	}

	return true
}

// IsValidBirthdate 判断是否为有效的出生日期
func IsValidBirthdate(date string) bool {
	// 使用正则表达式匹配日期格式
	pattern := `^(\d{4})-(\d{2})-(\d{2})$`
	matched, err := regexp.MatchString(pattern, date)
	if err != nil {
		return false
	}
	if !matched {
		return false
	}

	// 解析日期字符串
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		return false
	}

	// 检查日期是否在合理范围内
	now := time.Now()
	minAge := 0
	maxAge := 120
	minBirthYear := now.Year() - maxAge
	maxBirthYear := now.Year() - minAge
	if t.Year() < minBirthYear || t.Year() > maxBirthYear {
		return false
	}

	return true
}

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
