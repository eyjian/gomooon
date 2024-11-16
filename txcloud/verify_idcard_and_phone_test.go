// Package txcloud
// Wrote by yijian on 2024/11/16
package txcloud

import (
	"os"
	"testing"
)

// go test -v -run="TestBatchVerifyIdcardAndPhone" -args secret_id secret_key
func TestBatchVerifyIdcardAndPhone(t *testing.T) {
	args := os.Args[5:]
	secretId := args[0]
	secretKey := args[1]

	data := map[string]*IdcardPhoneTuple{
		"123456789012345678": {
			Idcard: "123456789012345678",
			Name:   "张三",
			Phone: "13800138000",
		},
		"123456789012345679": {
			Idcard: "123456789012345679",
			Name:   "李四",
			Phone: "12345678901234567890",
		},
		"123456789012345680": {
			Idcard: "123456789012345680",
			Name:   "王五",
			Phone: "12345678901234567890",
		},
		"123456789012345681": {
			Idcard: "123456789012345681",
			Name:   "赵六",
			Phone: "12345678901234567890",
		},
	}

	txCloud := NewFace(secretId, secretKey)
	consistent, inconsistent, fail := txCloud.BatchVerifyIdcardAndPhone(2, data)
	t.Logf("consistent:%d, inconsistent:%d, fail:%d\n", consistent, inconsistent, fail)
	for k, v := range data {
		t.Logf("%s: %+v\n", k, v)
		errCode, errMsg := GetErrCodeAndErrMsg(v.Err)
		if errCode != "" || errMsg != "" {
			t.Logf("ErrCode is `%s`, ErrMessage is `%s`\n", errCode, errMsg)
		} else {
			t.Logf("%s\n", v.Err.Error())
		}
	}
}

// go test -v -run="TestVerifyIdcardAndPhone" -args secret_id secret_key idcard name phone
func TestVerifyIdcardAndPhone(t *testing.T) {
	args := os.Args[5:]
	secretId := args[0]
	secretKey := args[1]
	idcard := args[2]
	name := args[3]
	phone := args[4]

	txCloud := NewFace(secretId, secretKey)
	ok, desc, err := txCloud.VerifyIdcardAndPhone(idcard, name, phone)
	if err != nil {
		errCode, errMsg := GetErrCodeAndErrMsg(err)
		if errCode != "" || errMsg != "" {
			t.Logf("ErrCode is `%s`, ErrMessage is `%s`\n", errCode, errMsg)
		} else {
			t.Logf("%s\n", err.Error())
		}
	} else if !ok {
		t.Errorf("%s",desc)
	} else {
		t.Log(ok, desc)
	}
}
