// Package mooonutils
// Wrote by yijian on 2024/01/02
package mooonutils

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// GetCurrentTimestamp 获取指定时区的当前时间戳
func GetCurrentTimestamp(timezone string) int64 {
	// 加载时区信息（北京时区：Asia/Shanghai）
	location, _ := time.LoadLocation(timezone)

	// 获取当前时间，并转换为北京时区的时间
	now := time.Now().In(location)

	// 将北京时区的时间转换为时间戳（精确到秒）
	timestamp := now.Unix()

	// 返回时间戳
	return timestamp
}

func GetCurrentTimestampString(timezone string) string {
	currentTimestamp := GetCurrentTimestamp(timezone)
	return strconv.FormatInt(currentTimestamp, 10)
}

func IsValidTime(s string) bool {
	// 选择一个合适的时间格式，例如：2006-01-02 15:04:05
	layout := "2006-01-02 15:04:05"
	_, err := time.Parse(layout, s)
	return err == nil
}

// NormalizeDateTimeString 将中文日期或时间字符串规整为“YYYY-MM-DD hh:mm:ss”格式
// withHms 参数 str 的值是否包含了"时、分、秒"，当值为 true 时，如果 str 本身不含时分秒，则结果会自动添加上
// str 日期或者时间字符串，格式可为：
// 1）YYYY年MM月DD日 hh时mm分ss秒
// 2）YYYY年MM月DD日hh时mm分ss秒
// 3）YYYY年MM月DD日
// 4）YYYY-MM-DD hh:mm:ss
// 5）YYYY-MM-DD
// 6）YYYY/MM/DD hh:mm:ss
// 7）YYYY/MM/DD
func NormalizeDateTimeString(str string, withHms bool) string {
	var builder strings.Builder
	runes := []rune(str)     // 将字符串转换为 rune 切片
	builder.Grow(len(runes)) // 预分配足够的空间以提高性能

	for i, r := range runes {
		switch r {
		case '年', '月':
			builder.WriteString("-")
		case '日':
			// 如果"日"后面没有空格，我们添加一个
			if i < len(runes)-1 && runes[i+1] != ' ' {
				builder.WriteString(" ")
			}
		case '时', '分':
			if withHms {
				builder.WriteRune(':')
			}
		case '秒':
			if withHms {
				builder.WriteString("")
			}
		default:
			builder.WriteRune(r)
		}
	}
	if withHms && builder.Len() == len("YYYY-MM-DD") {
		builder.WriteString(" 00:00:00")
	}

	return strings.ReplaceAll(builder.String(), "/", "-")
}

// String2Time 将日期字符串转换为时间对象
// 如果 dateStr 为"YYYY年MM月DD日"格式，可将"年、月、日"替换为"-"后再调用此函数
func String2Time(dateStr string) (time.Time, error) {
	// 定义支持的日期格式
	formats := []string{
		"2006-01-02", // YYYY-MM-DD
		"2006-1-2",   // YYYY-M-D
		"2006/01/02", // YYYY/MM/DD
		"2006/1/2",   // YYYY/M/D
		"2006.01.02", // YYYY.MM.DD
		"2006.1.2",   // YYYY.M.D
		"20060102",   // YYYYMMDD
		//"01-02-06",            // MM-DD-YY
		//"01/02/06",            // MM/DD/YY
		"02-Jan-2006",         // DD-MMM-YYYY (e.g., 02-Jan-2006)
		"02 Jan 2006",         // DD MMM YYYY (e.g., 02 Jan 2006)
		"2006-Jan-02",         // YYYY-MMM-DD (e.g., 2006-Jan-02)
		"2006 Jan 02",         // YYYY MMM DD (e.g., 2006 Jan 02)
		"02/Jan/2006",         // DD/MMM/YYYY (e.g., 02/Jan/2006)
		"02.Jan.2006",         // DD.MMM.YYYY (e.g., 02.Jan.2006)
		"2006/Jan/02",         // YYYY/MMM/DD (e.g., 2006/Jan/02)
		"2006.Jan.02",         // YYYY.MMM.DD (e.g., 2006.Jan.02)
		"Sunday, 02 Jan 2006", // Sunday, DD MMM YYYY (e.g., Sunday, 02 Jan 2006)
		"Sun, 02 Jan 2006",    // Sun, DD MMM YYYY (e.g., Sun, 02 Jan 2006)
		"02 Jan 2006",         // DD MMM YYYY (e.g., 02 Jan 2006)
		"Jan 02, 2006",        // MMM DD, YYYY (e.g., Jan 02, 2006)
		"Jan 02 2006",         // MMM DD YYYY (e.g., Jan 02 2006)
		"2006 Jan 02",         // YYYY MMM DD (e.g., 2006 Jan 02)
	}

	// 尝试使用不同的格式解析日期
	for _, format := range formats {
		t, err := time.Parse(format, dateStr)
		if err == nil {
			return t, nil
		}
	}

	// 如果所有格式都无法解析，返回错误
	return time.Time{}, fmt.Errorf("unsupported date format")
}