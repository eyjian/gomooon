// Package mooonutils
// Wrote by yijian on 2024/01/18
package mooonutils

import (
	"os"
	"testing"
	"time"
)

// go test -v -run="TestGetNonceStr"
func TestGetNonceStr(t *testing.T) {
	// 测试长度 1
	strLen := 1
	str := GetNonceStr(strLen)
	if len(str) == strLen {
		t.Logf("EXCEPTED: %s\n", str)
	} else {
		t.Errorf("UNEXCEPTED: %s\n", str)
	}

	// 测试长度 28
	strLen = 28
	str = GetNonceStr(28)
	if len(str) == strLen {
		t.Logf("EXCEPTED: %s\n", str)
	} else {
		t.Errorf("UNEXCEPTED: %s\n", str)
	}
}

// go test -v -run="TestGetLowercaseNonceStr"
func TestGetLowercaseNonceStr(t *testing.T) {
	strLen := 32
	str := GetLowercaseNonceStr(strLen)
	if len(str) == strLen {
		t.Logf("EXCEPTED: %s\n", str)
	} else {
		t.Errorf("UNEXCEPTED: %s\n", str)
	}
}

// go test -v -run="TestGetUppercaseNonceStr"
func TestGetUppercaseNonceStr(t *testing.T) {
	strLen := 32
	str := GetUppercaseNonceStr(strLen)
	if len(str) == strLen {
		t.Logf("EXCEPTED: %s\n", str)
	} else {
		t.Errorf("UNEXCEPTED: %s\n", str)
	}
}

// go test -v -run="TestGetHexNonceStr"
func TestGetHexNonceStr(t *testing.T) {
	strLen := 32
	str := GetHexNonceStr(strLen)
	if len(str) == strLen {
		t.Logf("EXCEPTED: %s\n", str)
	} else {
		t.Errorf("UNEXCEPTED: %s\n", str)
	}
}

// go test -v -run="TestDesensitizeStr"
func TestDesensitizeStr(t *testing.T) {
	str := ""
	mask := DesensitizeStr(str, 2, 2)
	t.Logf("`%s` => %s\n", str, mask)

	str = "11204416541220243X"
	mask = DesensitizeStr(str, 2, 2)
	t.Logf("`%s` => %s\n", str, mask)

	mask = DesensitizeStr(str, 3, 2)
	t.Logf("`%s` => %s\n", str, mask)

	mask = DesensitizeStr(str, 2, 1)
	t.Logf("`%s` => %s\n", str, mask)

	mask = DesensitizeStr(str, 2, 0)
	t.Logf("`%s` => %s\n", str, mask)

	mask = DesensitizeStr(str, 0, 2)
	t.Logf("`%s` => %s\n", str, mask)

	str = "112044165412202"
	mask = DesensitizeStr(str, 2, 2)
	t.Logf("`%s` => %s\n", str, mask)

	str = "112044165412203"
	mask = DesensitizeStr(str, 100, 2)
	t.Logf("`%s` => %s\n", str, mask)

	str = "112044165412204"
	mask = DesensitizeStr(str, 2, 200)
	t.Logf("`%s` => %s\n", str, mask)

	str = "1"
	mask = DesensitizeStr(str, 2, 2)
	t.Logf("`%s` => %s\n", str, mask)

	str = "123"
	mask = DesensitizeStr(str, 2, 2)
	t.Logf("`%s` => %s\n", str, mask)

	str = ""
	mask = DesensitizeStr(str, 2, 2)
	t.Logf("[EMPTY] `%s` => %s\n", str, mask)
}

// go test -v -run="TestDesensitizeChineseName"
func TestDesensitizeChineseName(t *testing.T) {
	name := ""
	mask := DesensitizeChineseName(name, 0, 1)
	t.Logf("%s => %s\n", name, mask)

	name = "张三"
	mask = DesensitizeChineseName(name, 0, 1)
	t.Logf("%s => %s\n", name, mask)

	name = "张三"
	mask = DesensitizeChineseName(name, 1, 1)
	t.Logf("%s => %s\n", name, mask)

	name = "王麻子"
	mask = DesensitizeChineseName(name, 1, 1)
	t.Logf("%s => %s\n", name, mask)

	name = "王麻子"
	mask = DesensitizeChineseName(name, 0, 1)
	t.Logf("%s => %s\n", name, mask)

	name = "欧阳大侠"
	mask = DesensitizeChineseName(name, 1, 1)
	t.Logf("%s => %s\n", name, mask)

	name = "欧阳大侠"
	mask = DesensitizeChineseName(name, 0, 1)
	t.Logf("%s => %s\n", name, mask)

	name = "张三.欧阳大侠"
	mask = DesensitizeChineseName(name, 1, 1)
	t.Logf("%s => %s\n", name, mask)

	name = "张三.欧阳大侠"
	mask = DesensitizeChineseName(name, 0, 1)
	t.Logf("%s => %s\n", name, mask)

	name = "张三"
	mask = DesensitizeChineseName(name, 1, 0)
	t.Logf("%s => %s\n", name, mask)

	name = "王麻子"
	mask = DesensitizeChineseName(name, 1, 0)
	t.Logf("%s => %s\n", name, mask)

	name = "欧阳大侠"
	mask = DesensitizeChineseName(name, 1, 0)
	t.Logf("%s => %s\n", name, mask)

	name = "张三.欧阳大侠"
	mask = DesensitizeChineseName(name, 1, 0)
	t.Logf("%s => %s\n", name, mask)
}

// 身份证号隐私数据，执行时指定
// go test -v -run="TestIsResidentIdentityCardNumber1" -args FLAG ID // FLAG 只能取值 0 或者 1，1 表示 ID 为无效身份证号，0 表示为有效的身份证号
func TestIsResidentIdentityCardNumber1(t *testing.T) {
	flag := os.Args[len(os.Args)-2]
	id := os.Args[len(os.Args)-1]
	if IsResidentIdentityCardNumber(id) {
		if flag == "0" {
			t.Errorf("%s is ID number\n", id)
		} else {
			t.Logf("%s is ID number\n", id)
		}
	} else {
		if flag == "0" {
			t.Logf("%s is not ID number\n", id)
		} else {
			t.Errorf("%s is not ID number\n", id)
		}
	}
}

// go test -v -run="TestIsResidentIdentityCardNumber2"
func TestIsResidentIdentityCardNumber2(t *testing.T) {
	id := "110000202405290095"
	if IsResidentIdentityCardNumber(id) {
		t.Logf("%s ok\n", id)
	} else {
		t.Errorf("%s is not a valid id\n", id)
	}

	id = "110000202405290096"
	if IsResidentIdentityCardNumber(id) {
		t.Errorf("%s ok\n", id)
	} else {
		t.Logf("%s is not a valid id\n", id)
	}

	id = "610125199107084337"
	if IsResidentIdentityCardNumber(id) {
		t.Logf("%s ok\n", id)
	} else {
		t.Errorf("%s is not a valid id\n", id)
	}

	id = "320381198812252138"
	if IsResidentIdentityCardNumber(id) {
		t.Logf("%s ok\n", id)
	} else {
		t.Errorf("%s is not a valid id\n", id)
	}

	id = "230124196911070015"
	if IsResidentIdentityCardNumber(id) {
		t.Logf("%s ok\n", id)
	} else {
		t.Errorf("%s is not a valid id\n", id)
	}
}

// go test -v -run="TestGenerateResidentIdentityCardNumber$"
func TestGenerateResidentIdentityCardNumber(t *testing.T) {
	areaCode := "110000" // 北京市
	birthDate := time.Now().Format("20060102")
	for i := 1; i <= 10; i++ {
		id, err := GenerateResidentIdentityCardNumber(areaCode, birthDate, i)
		if err != nil {
			t.Errorf("%s\n", err.Error())
		} else {
			t.Logf("%s\n", id)
		}
	}
}

// go test -v -run="TestTruncateUtf8String"
func TestTruncateUtf8String(t *testing.T) {
	s1 := "公元2024年01月01日的天气阳光明媚"
	s2 := "mooonutils"
	s3 := "生活尽是美好无线"

	s10 := TruncateUtf8String(s1, 16)
	if s10 != "公元2024年01月01日的天气" {
		t.Errorf("%s truncated to 16: %s\n", s1, s10)
	} else {
		t.Logf("%s truncated to 16: %s\n", s1, s10)
	}

	s20 := TruncateUtf8String(s2, 5)
	if s20 != "mooon" {
		t.Errorf("%s truncated to 5: %s\n", s2, s20)
	} else {
		t.Logf("%s truncated to 5: %s\n", s2, s20)
	}

	s30 := TruncateUtf8String(s3, 2)
	if s30 != "生活" {
		t.Errorf("%s truncated to 2: %s\n", s3, s30)
	} else {
		t.Logf("%s truncated to 2: %s\n", s3, s30)
	}
}

// go test -v -run="TestIsExecutableFile"
func TestIsExecutableFile(t *testing.T) {
	filepath := "/usr/bin/ls"
	is, err := IsExecutableFile(filepath)
	if err != nil {
		t.Errorf("IsExecutableFile error: %s\n", err.Error())
	} else {
		if is {
			t.Logf("%s is executable\n", filepath)
		} else {
			t.Errorf("%s is unexecutable\n", filepath)
		}
	}

	filepath = "/etc/profile"
	is, err = IsExecutableFile(filepath)
	if err != nil {
		t.Errorf("IsExecutableFile error: %s\n", err.Error())
	} else {
		if is {
			t.Errorf("%s is executable\n", filepath)
		} else {
			t.Logf("%s is unexecutable\n", filepath)
		}
	}
}