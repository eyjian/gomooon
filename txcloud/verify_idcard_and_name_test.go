// Package txcloud
// Wrote by yijian on 2024/11/16
package txcloud

import (
	"os"
	"testing"
)

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
