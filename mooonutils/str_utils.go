// Package mooonutils
// Wrote by yijian on 2024/01/02
package mooonutils

import (
	"fmt"
	"math/rand"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

const allCharset1 = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const allCharset2 = "abcdefghijklmnopqrstuvwxyz0123456789" // 不含大写字母
const allCharset3 = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789" // 不含小写字母
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
	return getNonceStr(length, allCharset1)
}

func GetLowercaseNonceStr(length int) string {
	return getNonceStr(length, allCharset2)
}

func GetUppercaseNonceStr(length int) string {
	return getNonceStr(length, allCharset3)
}

func GetHexNonceStr(length int) string {
	return getNonceStr(length, hexCharset)
}

// DesensitizeStr 脱敏字符串
// m 保留的前 m 个字
// n 保留的后 n 个字
func DesensitizeStr(str string, m, n int) string {
	strLen := len(str)

	// 显示前 m 位
	start := 0
	if m > strLen {
		m = strLen
	}
	visibleStart := str[start:m]

	// 显示后 n 位
	end := strLen - n
	if end < 0 {
		end = 0
	}
	visibleEnd := str[end:]

	// 生成脱敏后的值
	masked := visibleStart + strings.Repeat("*", strLen-m-n) + visibleEnd

	return masked
}

// DesensitizeName 脱敏姓名
// name 姓名，少数民族的姓名中间可能有点号
// m 保留的前 m 个字
// n 保留的后 n 个字
// k 单个字的字节数
func DesensitizeName(name string, m, n, k int) string {
	if k <= 0 {
		k = 3
	}

	var result []rune
	runes := []rune(name)
	dotIndex := -1

	for i, r := range runes {
		if r == '.' {
			dotIndex = i
			break
		}
	}

	// If the name has only two characters
	if len(runes) == 2 {
		if m != 0 {
			n = 0
		}
	}

	if dotIndex == -1 {
		result = append(result, runes[:m]...)
		result = append(result, []rune(strings.Repeat("*", len(runes)-m-n))...)
		result = append(result, runes[len(runes)-n:]...)
	} else {
		result = append(result, runes[:m]...)
		result = append(result, []rune(strings.Repeat("*", dotIndex-m))...)
		result = append(result, '.')

		remain := len(runes) - dotIndex - 1 - n
		if remain > 0 {
			result = append(result, []rune(strings.Repeat("*", remain))...)
		}

		result = append(result, runes[len(runes)-n:]...)
	}

	return string(result)
}

// IsResidentIdentityCardNumber 判断是否为居民身份证号
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
		return calculateCheckDigit(id) == id[17:]
	}

	return true
}

// calculateCheckDigit 计算居民身份证号的校验码
func calculateCheckDigit(id string) string {
	weights := []int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
	checkSum := 0
	for i := 0; i < 17; i++ {
		digit, _ := strconv.Atoi(string(id[i]))
		checkSum += digit * weights[i]
	}
	checkCodes := "10X98765432"
	return string(checkCodes[checkSum%11])
}

// GenerateResidentIdentityCardNumber 生成有效的居民身份证号，最后一位是根据前17位计算得出的校验码
// areaCode 六位数字行政区划代码，如：440306
// birthDate 八位数字出生日期，如：20240529
// sequence 三位顺序码，奇数分配给男性，偶数分配给女性
func GenerateResidentIdentityCardNumber(areaCode string, birthDate string, sequence int) (string, error) {
	if len(areaCode) != 6 {
		return "", fmt.Errorf("area code must be 6 digits")
	}
	if len(birthDate) != 8 {
		return "", fmt.Errorf("birthdate must be 8 digits")
	}
	if sequence < 1 || sequence > 999 {
		return "", fmt.Errorf("sequence number must be an integer between 1 and 999")
	}
	sequenceStr := fmt.Sprintf("%03d", sequence)
	id := areaCode + birthDate + sequenceStr
	checkDigit := calculateCheckDigit(id)
	return id + checkDigit, nil
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

// TruncateUtf8String 截取 UTF8 字符串，使其字数（不是字节数，一个数字、字母和汉字都分别计 1）不超过 maxCharCount
func TruncateUtf8String(utf8Str string, maxCharCount int) string {
	if utf8.RuneCountInString(utf8Str) <= maxCharCount {
		return utf8Str
	}

	var truncated string
	var charCount int
	for len(utf8Str) > 0 {
		r, size := utf8.DecodeRuneInString(utf8Str)
		if charCount+1 > maxCharCount {
			break
		}
		truncated += string(r)
		charCount++
		utf8Str = utf8Str[size:]
	}
	return truncated
}

// CountUtf8Characters 计算字数，一个数字、字母和汉字都分别计 1
func CountUtf8Characters(utf8Str string) int {
	return utf8.RuneCountInString(utf8Str)
}

// ExtractUrlPath 提取 url 路径
func ExtractUrlPath(urlStr string) string {
	parsedUrl, _ := url.Parse(urlStr)
	parsedUrl.Scheme = ""
	parsedUrl.Host = ""
	return parsedUrl.String()
}

// ConvertDateFormat 将 YYYY-MM-DD 格式的日期转换为 YYYYMMDD 格式
func ConvertDateFormat(dateStr string) string {
	return strings.ReplaceAll(dateStr, "-", "")
}
