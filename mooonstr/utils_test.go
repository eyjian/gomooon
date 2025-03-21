// Package mooonstr
// Wrote by yijian on 2024/09/03
package mooonstr

import "testing"

// go test -v -run="TestCamelCase"
func TestCamelCase(t *testing.T) {
	str := "hello_world"
	res := CamelCase(str)
	t.Logf("result: %s => %s\n", str, res)

	str = "hello"
	res = CamelCase(str)
	t.Logf("result: %s => %s\n", str, res)

	str = "hello1d"
	res = CamelCase(str)
	t.Logf("result: %s => %s\n", str, res)
}

// go test -v -run="TestJoinInt32$"
func TestJoinInt32(t *testing.T) {
	elems := []int32{-1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	res := JoinInt32(elems, ",")
	t.Logf("result: %s\n", res)
}

// go test -v -run="TestJoinUint32$"
func TestJoinUint32(t *testing.T) {
	elems := []uint32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	res := JoinUint32(elems, ",")
	t.Logf("result: %s\n", res)
}

// go test -v -run="TestJoinInt64$"
func TestJoinInt64(t *testing.T) {
	elems := []int64{-10, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	res := JoinInt64(elems, ",")
	t.Logf("result: %s\n", res)
}

// go test -v -run="TestJoinUint64$"
func TestJoinUint64(t *testing.T) {
	elems := []uint64{10, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	res := JoinUint64(elems, ",")
	t.Logf("result: %s\n", res)
}

// go test -v -run="TestLuhnCheck$"
func TestLuhnCheck(t *testing.T) {
	number := ""
	if LuhnCheck(number) {
		t.Errorf("`%s` ok\n", number)
	} else {
		t.Logf("`%s` error\n", number)
	}

	number = " "
	if LuhnCheck(number) {
		t.Errorf("`%s` ok\n", number)
	} else {
		t.Logf("`%s` error\n", number)
	}

	number = "1234567890"
	if LuhnCheck(number) {
		t.Errorf("`%s` ok\n", number)
	} else {
		t.Logf("`%s` error\n", number)
	}

	number = "abc"
	if LuhnCheck(number) {
		t.Errorf("`%s` ok\n", number)
	} else {
		t.Logf("`%s` error\n", number)
	}

	number = "zzz@"
	if LuhnCheck(number) {
		t.Errorf("`%s` ok\n", number)
	} else {
		t.Logf("`%s` error\n", number)
	}

	// 中国银行卡
	number = "6217858000107987452"
	if LuhnCheck(number) {
		t.Logf("`%s` ok\n", number)
	} else {
		t.Errorf("`%s` error\n", number)
	}

	// 广发银行卡
	number = "6214620221003803700"
	if LuhnCheck(number) {
		t.Logf("`%s` ok\n", number)
	} else {
		t.Errorf("`%s` error\n", number)
	}

	// 中国工商银行卡
	number = "6212260712007585433"
	if LuhnCheck(number) {
		t.Logf("`%s` ok\n", number)
	} else {
		t.Errorf("`%s` error\n", number)
	}

	// 中国农业银行卡
	number = "6228480425523649274"
	if LuhnCheck(number) {
		t.Logf("`%s` ok\n", number)
	} else {
		t.Errorf("`%s` error\n", number)
	}

	// 中国建设银行卡
	number = "6217002730013587765"
	if LuhnCheck(number) {
		t.Logf("`%s` ok\n", number)
	} else {
		t.Errorf("`%s` error\n", number)
	}

	// 中信银行
	number = "6217731800737179"
	if LuhnCheck(number) {
		t.Logf("`%s` ok\n", number)
	} else {
		t.Errorf("`%s` error\n", number)
	}

	// 民生银行
	number = "6226220158339332"
	if LuhnCheck(number) {
		t.Logf("`%s` ok\n", number)
	} else {
		t.Errorf("`%s` error\n", number)
	}

	// 湖北银行卡
	number = "190200120100010316"
	if LuhnCheck(number) {
		t.Logf("`%s` ok\n", number)
	} else {
		t.Errorf("`%s` error\n", number)
	}

	// 湖北省农村信用社卡
	number = "6224121177445138"
	if LuhnCheck(number) {
		t.Logf("`%s` ok\n", number)
	} else {
		t.Errorf("`%s` error\n", number)
	}

	// 恒丰银行卡
	number = "6230780100034858135"
	if LuhnCheck(number) {
		t.Logf("`%s` ok\n", number)
	} else {
		t.Errorf("`%s` error\n", number)
	}

	// 重庆银行卡
	number = "6230962100001869974"
	if LuhnCheck(number) {
		t.Logf("`%s` ok\n", number)
	} else {
		t.Errorf("`%s` error\n", number)
	}

	// 宁波银行卡
	number = "64090122000033385"
	if LuhnCheck(number) {
		t.Logf("`%s` ok\n", number)
	} else {
		t.Errorf("`%s` error\n", number)
	}

	// 徽商银行卡
	number = "225019421181000002"
	if LuhnCheck(number) {
		t.Logf("`%s` ok\n", number)
	} else {
		t.Errorf("`%s` error\n", number)
	}

	// 九江银行卡
	number = "337119700000009424"
	if LuhnCheck(number) {
		t.Logf("`%s` ok\n", number)
	} else {
		t.Errorf("`%s` error\n", number)
	}

	// 哈尔滨银行卡
	number = "6217524511104903114"
	if LuhnCheck(number) {
		t.Logf("`%s` ok\n", number)
	} else {
		t.Errorf("`%s` error\n", number)
	}

	// 中国邮政储蓄银行卡
	number = "6217993300063741836"
	if LuhnCheck(number) {
		t.Logf("`%s` ok\n", number)
	} else {
		t.Errorf("`%s` error\n", number)
	}
}

// go test -v -run="TestFormatCents"
func TestFormatCents(t *testing.T) {
	// 编写测试用例
	tests := []struct {
		cents    int
		expected string
	}{
		{0, "0"},
		{1, "0.01"},
		{10, "0.1"},
		{100, "1"},
		{1000, "10"},
	}

	// 遍历测试用例并执行测试
	for _, tt := range tests {
		actual := FormatCents(uint32(tt.cents))
		if actual != tt.expected {
			t.Errorf("FormatCents(%d) = %s, expected %s", tt.cents, actual, tt.expected)
		}
	}
}

// go test -v -run="TestIsAllChinese"
// TestIsAllChinese 测试是否都是中文
func TestIsAllChinese(t *testing.T) {
	// 编写测试用例
	tests := []struct {
		str      string
		expected bool
	}{
		{"", false},
		{"abc", false},
		{"abc123", false},
		{"123abc", false},
		{"123abc456", false},
		{"abc123456", false},
		{"123456", false},
		{"一二三", true},
		{"一，二三", true},
		{"一二三四五六", true},
		{"一二三四五六七八九", true},
		{"一二三四五六七八九十", true},
		{"一二三四五六七八九十123", false},
		{"123一二三四五六七八九十", false},
		{"123一二三四五六七八九十123", false},
		{"123一二三四五六七八九十", false},
		{"123一二三四五六七八九十123", false},
		{"123一二三四五六七八九十", false},
		{"123一二三四五六七八九十123", false},
		{"123一二三四五六七八九十", false},
	}

	// 遍历测试用例并执行测试
	for _, tt := range tests {
		actual := IsAllChinese(tt.str)
		if actual != tt.expected {
			t.Errorf("IsAllChinese(%q) = %v, expected %v", tt.str, actual, tt.expected)
		}
	}
}

// go test -v -run="TestContainsChinese"
func TestContainsChinese(t *testing.T) {
	// 编写测试用例
	tests := []struct {
		str      string
		expected bool
	}{
		{"", false},
		{"abc", false},
		{"abc123", false},
		{"123abc", false},
		{"123abc456", false},
		{"abc123456", false},
		{"123456", false},
		{"123456，", true},
		{"一二三", true},
		{"一二三四五六", true},
		{"一二三四五六七八九", true},
		{"一二三四五六七八九十", true},
		{"一二三四五六七八九十123", true},
		{"123一二三四五六七八九十", true},
		{"123一二三四五六七八九十123", true},
		{"123一二三四五六七八九十", true},
		{"123一二三四五六七八九十123", true},
		{"123一二三四五六七八九十", true},
		{"123一二三四五六七八九十123", true},
		{"123一二三四五六七八九十", true},
	}

	// 遍历测试用
	for _, tt := range tests {
		actual := ContainsChinese(tt.str)
		if actual != tt.expected {
			t.Errorf("ContainsChinese(%q) = %v, expected %v", tt.str, actual, tt.expected)
		}
	}
}