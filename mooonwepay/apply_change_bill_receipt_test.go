// Package mooonwepay
// Wrote by yijian on 2024/08/23
package mooonwepay

import (
	"context"
	"github.com/eyjian/gomooon/moooncrypto"
	"github.com/eyjian/gomooon/mooonutils"
	"net/http"
	"os"
	"testing"
	"time"
)

// go test -v -run="TestApplyChangeBillReceipt$" -args private_key.pem mchid serial_no out_batch_no <out_detail_no>
func TestApplyChangeBillReceipt(t *testing.T) {
	numArgs := len(os.Args)
	t.Log("args num:", numArgs)
	if numArgs != 9 && numArgs != 10 {
		t.Error("args error")
		return
	}

	var args []string
	if numArgs == 9 {
		args = os.Args[len(os.Args)-4:]
	} else {
		args = os.Args[len(os.Args)-5:]
	}
	privateKeyFilepath := args[0]
	mchid := args[1]
	serialNo := args[2]
	outBatchNo := args[3]
	outDetailNo := ""
	if len(args) > 4 {
		outDetailNo = args[4]
	}

	privateKey, err := moooncrypto.Filepath2PrivateKey(privateKeyFilepath)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("privateKey ok\n")

	timestamp := time.Now().Unix()
	nonceStr := mooonutils.GetNonceStr(32)
	req := &ApplyChangeBillReceiptReq{
		Ctx:        context.Background(),
		HttpClient: &http.Client{},
		PrivateKey: privateKey,

		Host:      "https://api.mch.weixin.qq.com",
		NonceStr:  nonceStr,
		Timestamp: timestamp,

		Mchid:       mchid,
		SerialNo:    serialNo,
		OutBatchNo:  outBatchNo,
		OutDetailNo: outDetailNo,
		AcceptType:  "BATCH_TRANSFER",
	}
	t.Log(*req)
	resp, err := ApplyChangeBillReceipt(req)
	if err != nil {
		t.Error(err)
	}
	if resp != nil {
		t.Logf("%+v\n", *resp)
	}
}
