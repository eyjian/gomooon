// Package mooonstr
// Wrote by yijian on 2024/09/03
package mooonstr

import (
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"regexp"
	"strconv"
	"strings"
)

// CamelCase 驼峰命名法
// 蛇形命名法可以使用 github.com/stoewer/go-strcase.SnakeCase
func CamelCase(s string) string {
	// 使用 strings.Split 按"_"分割字符串
	parts := strings.Split(s, "_")

	// 遍历分割后的字符串切片
	for i, part := range parts {
		// 将每个单词的首字母大写，其余字母小写
		// 也可使用 strings.Title，但 strings.Title 是个 Deprecated 函数
		parts[i] = cases.Title(language.English, cases.NoLower).String(strings.ToLower(part))
	}

	// 使用 strings.Join 将处理后的字符串切片连接起来
	return strings.Join(parts, "")
}

// JoinInt32 使用分隔符将多个 int32 串拼接成一个字符串
func JoinInt32(elems []int32, sep string) string {
	// 将 []int32 转换为 []string
	strElems := make([]string, len(elems))
	for i, num := range elems {
		strElems[i] = strconv.FormatInt(int64(num), 10)
	}

	// 使用 strings.Join 函数拼接字符串
	return strings.Join(strElems, sep)
}

// JoinUint32 使用分隔符将多个 uint32 串拼接成一个字符串
func JoinUint32(elems []uint32, sep string) string {
	// 将 []uint32 转换为 []string
	strElems := make([]string, len(elems))
	for i, num := range elems {
		strElems[i] = strconv.FormatInt(int64(num), 10)
	}

	// 使用 strings.Join 函数拼接字符串
	return strings.Join(strElems, sep)
}

// JoinInt64 使用分隔符将多个 int64 串拼接成一个字符串
func JoinInt64(elems []int64, sep string) string {
	// 将 []int64 转换为 []string
	strElems := make([]string, len(elems))
	for i, num := range elems {
		strElems[i] = strconv.FormatInt(num, 10)
	}

	// 使用 strings.Join 函数拼接字符串
	return strings.Join(strElems, sep)
}

// JoinUint64 使用分隔符将多个 uint64 串拼接成一个字符串
func JoinUint64(elems []uint64, sep string) string {
	// 将 []uint64 转换为 []string
	strElems := make([]string, len(elems))
	for i, num := range elems {
		strElems[i] = strconv.FormatUint(num, 10)
	}

	// 使用 strings.Join 函数拼接字符串
	return strings.Join(strElems, sep)
}

func isValidLuhnLength(str string) bool {
	match, _ := regexp.MatchString(`^\d{8,19}$`, str)
	return match
}

// LuhnCheck 模10算法实现，可用于效验银行卡等
func LuhnCheck(str string) bool {
	clean := strings.ReplaceAll(str, " ", "")
	if clean == "" {
		return false
	}
	if !isValidLuhnLength(clean) {
		return false
	}

	sum := 0
	alt := false
	for i := len(clean) - 1; i >= 0; i-- {
		//digit := int(clean[i] - '0')
		digit, err := strconv.Atoi(string(clean[i]))
		if err != nil {
			return false
		}
		if alt {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
		alt = !alt
	}
	return sum%10 == 0
}