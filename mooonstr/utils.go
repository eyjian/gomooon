// Package mooonstr
// Wrote by yijian on 2024/09/03
package mooonstr

import (
	"fmt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
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

// IsChineseName 判断字符串是否为合法的中文姓名
func IsChineseName(str string, maxChineseCharacters int) bool {
	chineseCharacters := utf8.RuneCountInString(str)
	if chineseCharacters < 2 || chineseCharacters > maxChineseCharacters {
		return false
	}

	// 汉字+间隔符校验（unicode.Han 覆盖 CJK 统一汉字）
	validChars := map[rune]bool{
		'\u00B7': true, // Unicode 中间点符号
	}
	for _, r := range str {
		if !unicode.Is(unicode.Han, r) && !validChars[r] {
			return false
		}
	}

	// 正则表达式增强校验（少数民族姓名中间可能含点）
	pattern := `^[\p{Han}·]+$`         // 反引号内直接写字符
	reg := regexp.MustCompile(pattern) // 包含中间点符号
	return reg.MatchString(str)
}

// FormatCents 格式化分值
func FormatCents(cents uint32) string {
	// 分解元、角、分（100分=1元，10分=1角，1分=1分）
	yuan := cents / 100      // 元部分
	remainder := cents % 100 // 角分组合值
	jiao := remainder / 10   // 角位 (十位)
	fen := remainder % 10    // 分位 (个位)

	// 动态选择格式化规则
	switch {
	case fen != 0:
		return fmt.Sprintf("%d.%d%d", yuan, jiao, fen) // 完整两位小数：12.34
	case jiao != 0:
		return fmt.Sprintf("%d.%d", yuan, jiao) // 只保留角位：12.3
	default:
		return fmt.Sprintf("%d", yuan) // 整数格式：12
	}
}

// IsAllChinese 判断字符串是否全由中文组成（包含常用CJK汉字）
// 空字符串返回false
func IsAllChinese(str string) bool {
	if len(str) == 0 {
		return false
	}
	for _, r := range str {
		if !unicode.Is(unicode.Han, r) && !unicode.IsPunct(r) {
			return false
		}
	}
	return true
}

// ContainsChinese 判断字符串中是否包含中文字符，包括中文标点符号
// 空字符串返回false
func ContainsChinese(str string) bool {
	for _, r := range str {
		if unicode.Is(unicode.Han, r) || unicode.IsPunct(r) {
			return true
		}
	}
	return false
}