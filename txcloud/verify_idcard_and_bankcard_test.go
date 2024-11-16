// Package txcloud
// Wrote by yijian on 2024/11/16
package txcloud

import (
	"os"
	"testing"
)

// go test -v -run="TestVerifyIdcardAndBankcard" -args secret_id secret_key idcard name bankcard
func TestVerifyIdcardAndBankcard(t *testing.T) {
	args := os.Args[5:]
	secretId := args[0]
	secretKey := args[1]
	idcard := args[2]
	name := args[3]
	bankcard := args[4]

	txCloud := NewFace(secretId, secretKey)
	ok, desc, err := txCloud.VerifyIdcardAndBankcard(idcard, name, bankcard)
	if err != nil {
		t.Fatal(err)
	} else if !ok {
		t.Errorf("%s",desc)
	} else {
		t.Log(ok, desc)
	}
}
