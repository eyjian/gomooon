// Package txcloud
// Wrote by yijian on 2024/11/16
package txcloud

import (
	"os"
	"testing"
)

// go test -v -run="TestBatchVerifyIdcardAndName" -args secret_id secret_key
func TestBatchVerifyIdcardAndName(t *testing.T) {
	args := os.Args[5:]
	secretId := args[0]
	secretKey := args[1]

	data := map[string]*IdcardNameTuple{
		"123456789012345678": {
			Idcard: "123456789012345678",
			Name:   "张三",
		},
		"123456789012345679": {
			Idcard: "123456789012345679",
			Name:   "李四",
		},
		"123456789012345680": {
			Idcard: "123456789012345680",
			Name:   "王五",
		},
		"123456789012345681": {
			Idcard: "123456789012345681",
			Name:   "赵六",
		},
	}

	txCloud := NewFace(secretId, secretKey)
	consistent, inconsistent, fail := txCloud.BatchVerifyIdcardAndName(2, data)
	t.Logf("consistent:%d, inconsistent:%d, fail:%d\n", consistent, inconsistent, fail)
	for k, v := range data {
		t.Logf("%s: %+v\n", k, v)
	}
}

// go test -v -run="TestVerifyIdcardAndName" -args secret_id secret_key idcard name
func TestVerifyIdcardAndName(t *testing.T) {
	args := os.Args[5:]
	secretId := args[0]
	secretKey := args[1]
	idcard := args[2]
	name := args[3]

	txCloud := NewFace(secretId, secretKey)
	ok, desc, err := txCloud.VerifyIdcardAndName(idcard, name)
	if err != nil {
		t.Fatal(err)
	} else if !ok {
		t.Errorf("%s",desc)
	} else {
		t.Log(ok, desc)
	}
}
