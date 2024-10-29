// Package mooonutils
// Wrote by yijian on 2024/04/02
package mooonutils

import (
	"testing"
	"time"
)

// go test -v -run="TestGetCurrentTimestamp"
func TestGetCurrentTimestamp(t *testing.T) {
	timezone := "Asia/Shanghai"
	timestamp := GetCurrentTimestamp(timezone)

	if timestamp <= 0 {
		t.Errorf("GetCurrentTimestamp(%s) = %d; want a positive integer", timezone, timestamp)
	}
}

// go test -v -run="TestGetCurrentTimestampString"
func TestGetCurrentTimestampString(t *testing.T) {
	timezone := "Asia/Shanghai"
	timestampStr := GetCurrentTimestampString(timezone)

	if len(timestampStr) != 10 {
		t.Errorf("len(GetCurrentTimestampString(%s)) = %d; want 10", timezone, len(timestampStr))
	}
}

// go test -v -run="TestNormalizeDateTimeString"
func TestNormalizeDateTimeString(t *testing.T) {
	testCases := []struct {
		dateStr     string
		expectedStr string
	}{
		{"2020-01-02", "2020-01-02 00:00:00"},
		{"2021年01月02日", "2021-01-02 00:00:00"},
		{"2022-01-02 12:03:04", "2022-01-02 12:03:04"},
		{"2023年01月02日 12时03分04秒", "2023-01-02 12:03:04"},
		{"2024年01月02日12时03分04秒", "2024-01-02 12:03:04"},
		{"2025/01/02 12:03:04", "2025-01-02 12:03:04"},
		{"2026/01/02", "2026-01-02 00:00:00"},
	}
	for _, tc := range testCases {
		t.Run(tc.dateStr, func(t *testing.T) {
			cleanedStr := NormalizeDateTimeString(tc.dateStr, true)
			if cleanedStr != tc.expectedStr {
				t.Errorf("CleanDateTimeString(%s) => %s; want %s", tc.dateStr, cleanedStr, tc.expectedStr)
			} else {
				t.Logf("CleanDateTimeString(%s) => %s; want %s", tc.dateStr, cleanedStr, tc.expectedStr)
			}
		})
	}
}

// go test -v -run="TestString2Time"
func TestString2Time(t *testing.T) {
	testCases := []struct {
		dateStr      string
		expectedTime time.Time
	}{
		{"2024-01-02", time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)},
		{"2024-1-2", time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)},
		{"2024/01/02", time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)},
		{"2024/1/2", time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)},
		{"2024.01.02", time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)},
		{"2024.1.2", time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)},
		{"20240102", time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)},
		{"01-02-24", time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)},
		{"01/02/24", time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)},
		{"02-Jan-2024", time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)},
		{"02 Jan 2024", time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)},
		{"2024-Jan-02", time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)},
		{"2024 Jan 02", time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)},
		{"02/Jan/2024", time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)},
		{"02.Jan.2024", time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)},
		{"2024/Jan/02", time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)},
		{"2024.Jan.02", time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)},
		{"Sunday, 02 Jan 2024", time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)},
		{"Sun, 02 Jan 2024", time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)},
		{"02 Jan 2024", time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)},
		{"Jan 02, 2024", time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)},
		{"Jan 02 2024", time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)},
		{"2024 Jan 02", time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)},
		{"2024年02月01日", time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)},
	}

	for _, tc := range testCases {
		t.Run(tc.dateStr, func(t *testing.T) {
			result, err := String2Time(tc.dateStr)
			if err != nil {
				t.Errorf("String2Time(%s) returned error: %v", tc.dateStr, err)
			}
			if !result.Equal(tc.expectedTime) {
				t.Errorf("String2Time(%s) = %v; want %v", tc.dateStr, result, tc.expectedTime)
			} else {
				t.Logf("%s ok\n", tc.dateStr)
			}
		})
	}
}

// go test -v -run="TestString2Time_InvalidFormat"
func TestString2Time_InvalidFormat(t *testing.T) {
	dateStr := "invalid format"
	_, err := String2Time(dateStr)

	if err == nil {
		t.Errorf("String2Time(%s) did not return an error; want an error", dateStr)
	}
}