// Package txcloud
// Wrote by yijian on 2026/08/31
package txcloud

import (
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"os"
	"testing"
)

// go test -v -run="TestBatchVerifyEntNameAndCreditCode"
// 密钥从环境变量读取：TENCENTCLOUD_SECRET_ID 和 TENCENTCLOUD_SECRET_KEY
func TestBatchVerifyEntNameAndCreditCode(t *testing.T) {
	secretId := os.Getenv("TENCENTCLOUD_SECRET_ID")
	secretKey := os.Getenv("TENCENTCLOUD_SECRET_KEY")
	if secretId == "" || secretKey == "" {
		t.Skip("skip: TENCENTCLOUD_SECRET_ID or TENCENTCLOUD_SECRET_KEY not set")
	}

	// 下列均为虚构的占位数据（企业名称与信用代码非真实企业，仅用于演示接口调用方式），
	// 因此核验结果预期为不一致或查无结果。如需核验真实企业，请使用 TestVerifyEntNameAndCreditCode 并通过 -args 传入。
	data := map[string]*EntNameCreditCodeTuple{
		"91110000000000000X": {
			EntName:    "示例科技有限公司",
			CreditCode: "91110000000000000X",
		},
		"91440300000000001Y": {
			EntName:    "测试贸易有限公司",
			CreditCode: "91440300000000001Y",
		},
		"91310000000000002Z": {
			EntName:    "样例信息技术有限公司",
			CreditCode: "91310000000000002Z",
		},
	}

	txCloud := NewEntVerifier(secretId, secretKey)
	consistent, inconsistent, fail := txCloud.BatchVerifyEntNameAndCreditCode(2, len(data), data)
	t.Logf("consistent:%d, inconsistent:%d, fail:%d\n", consistent, inconsistent, fail)
	for k, v := range data {
		t.Logf("%s: %+v, OperatingStatusDesc: %s\n", k, v, OperatingStatusDesc(v.OperatingStatus))
		if e, ok := v.Err.(*errors.TencentCloudSDKError); ok {
			t.Logf("e.Code is `%s`, e.Message is `%s`\n", e.Code, e.Message)
		}
	}
}

// go test -v -run="TestVerifyEntNameAndCreditCode" -args ent_name credit_code
// 密钥从环境变量读取：TENCENTCLOUD_SECRET_ID 和 TENCENTCLOUD_SECRET_KEY
func TestVerifyEntNameAndCreditCode(t *testing.T) {
	secretId := os.Getenv("TENCENTCLOUD_SECRET_ID")
	secretKey := os.Getenv("TENCENTCLOUD_SECRET_KEY")
	if secretId == "" || secretKey == "" {
		t.Skip("skip: TENCENTCLOUD_SECRET_ID or TENCENTCLOUD_SECRET_KEY not set")
	}

	args := os.Args[5:]
	entName := args[0]
	creditCode := args[1]

	txCloud := NewEntVerifier(secretId, secretKey)
	ok, operatingStatus, requestId, err := txCloud.VerifyEntNameAndCreditCode(entName, creditCode)
	if err != nil {
		if e, ok := err.(*errors.TencentCloudSDKError); ok {
			t.Logf("e.Code is `%s`, e.Message is `%s`\n", e.Code, e.Message)
		} else {
			t.Logf("%s\n", err.Error())
		}
	} else if !ok {
		t.Errorf("(RequestId=%s) inconsistent, OperatingStatus: %s(%s)", requestId, OperatingStatusDesc(operatingStatus), operatingStatus)
	} else {
		t.Log(ok, OperatingStatusDesc(operatingStatus), requestId)
	}
}
