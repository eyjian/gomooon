// Package utils
// Wrote by yijian on 2024/01/02
package utils

import (
	"fmt"
	"math/rand"
	"regexp"
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
	// 使用正则表达式匹配身份证号格式
	pattern := `^(\d{15}|\d{17}[\dxX])$`
	matched, err := regexp.MatchString(pattern, id)
	if err != nil {
		return false
	}
	if !matched {
		return false
	}

	// 如果是15位身份证号，则转换为18位
	if len(id) == 15 {
		id = ConvertTo18ResidentIdentityCardNumber(id)
	}

	// 计算校验码并与身份证号中的校验码进行比较
	checkCode := calculateCheckCode(id)
	if checkCode == "" || string(id[len(id)-1]) != checkCode {
		return false
	}

	return true
}

// ConvertTo18ResidentIdentityCardNumber 将15位身份证号转换为18位
func ConvertTo18ResidentIdentityCardNumber(id string) string {
	weights := []int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
	checkCodes := "10X98765432"

	sum := 0
	for i := 0; i < 15; i++ {
		digit, _ := parseInt(string(id[i]))
		sum += digit * weights[i]
	}

	checkIndex := sum % 11
	checkDigit := checkCodes[checkIndex]

	return id + string(checkDigit)
}

// calculateCheckCode 计算身份证号的校验码
func calculateCheckCode(id string) string {
	weights := []int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
	checkCodes := "10X98765432"

	sum := 0
	for i := 0; i < 17; i++ {
		digit, _ := parseInt(string(id[i]))
		sum += digit * weights[i]
	}

	checkIndex := sum % 11
	return string(checkCodes[checkIndex])
}

// parseInt 将字符串转换为整数
func parseInt(s string) (int, error) {
	return fmt.Sscanf(s, "%d")
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
