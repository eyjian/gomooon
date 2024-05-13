// Package mooonutils
// Wrote by yijian on 2024/01/02
package mooonutils

import (
    "fmt"
    "strconv"
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

func String2Time(dateStr string) (time.Time, error) {
    // 定义支持的日期格式
    formats := []string{
        "2006-01-02",          // YYYY-MM-DD
        "2006-1-2",            // YYYY-M-D
        "2006/01/02",          // YYYY/MM/DD
        "2006/1/2",            // YYYY/M/D
        "2006.01.02",          // YYYY.MM.DD
        "2006.1.2",            // YYYY.M.D
        "20060102",            // YYYYMMDD
        "01-02-06",            // MM-DD-YY
        "01/02/06",            // MM/DD/YY
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
