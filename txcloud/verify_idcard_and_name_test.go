// Package txcloud
// Wrote by yijian on 2024/11/16
package txcloud

import (
	stderrors "errors"
	"os"
	"testing"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
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
	consistent, inconsistent, fail := txCloud.BatchVerifyIdcardAndName(2, len(data), data)
	t.Logf("consistent:%d, inconsistent:%d, fail:%d\n", consistent, inconsistent, fail)
	for k, v := range data {
		t.Logf("%s: %+v\n", k, v)
		if e, ok := v.Err.(*errors.TencentCloudSDKError); ok {
			t.Logf("e.Code is `%s`, e.Message is `%s`\n", e.Code, e.Message)
		}
	}
}

// go test -v -run="TestExtractRequestId"
func TestExtractRequestId(t *testing.T) {
	sdkErr := errors.NewTencentCloudSDKError("RequestLimitExceeded", "频率超限", "rid-sdk")
	if rid := extractRequestId(sdkErr); rid != "rid-sdk" {
		t.Errorf("sdk err: rid=%q, want %q", rid, "rid-sdk")
	}
	if rid := extractRequestId(stderrors.New("network timeout")); rid != "" {
		t.Errorf("non-sdk err: rid=%q, want empty", rid)
	}
}

// go test -v -run="TestJudgeFaceVerifyResult"
func TestJudgeFaceVerifyResult(t *testing.T) {
	newResp := func(result string) *FaceResponse {
		resp := &FaceResponse{}
		resp.Response.Result = result
		resp.Response.Description = "desc-" + result
		resp.Response.RequestId = "rid-" + result
		return resp
	}

	// 一致
	ok, _, _, err := judgeFaceVerifyResult(newResp("0"), idcardNameInconsistentResults)
	if !ok || err != nil {
		t.Fatalf("result=0: ok=%v, err=%v", ok, err)
	}

	// 各接口的确定性否定码判不一致
	cases := []struct {
		name  string
		codes []string
		set   map[string]struct{}
	}{
		{"idcard_name", []string{"-1", "-2", "-3", "-5"}, idcardNameInconsistentResults},
		{"idcard_phone", []string{"-4", "-6", "-7", "-8", "-9", "-13"}, idcardPhoneInconsistentResults},
		{"idcard_bankcard", []string{"-1", "-2", "-3", "-4", "-5", "-16", "-19", "-20"}, idcardBankcardInconsistentResults},
		{"name_bankcard", []string{"-1", "-2", "-3", "-4", "-15", "-18"}, nameBankcardInconsistentResults},
	}
	for _, c := range cases {
		for _, code := range c.codes {
			ok, desc, rid, err := judgeFaceVerifyResult(newResp(code), c.set)
			if ok || err != nil {
				t.Errorf("%s result=%s: ok=%v, err=%v", c.name, code, ok, err)
			}
			if desc != "desc-"+code || rid != "rid-"+code {
				t.Errorf("%s result=%s: desc=%q, rid=%q", c.name, code, desc, rid)
			}
		}
	}

	// 服务异常及未知码返回 *errors.TencentCloudSDKError
	errCases := []struct {
		name string
		code string
		set  map[string]struct{}
	}{
		{"idcard_name 证件库服务异常", "-4", idcardNameInconsistentResults},
		{"idcard_name 系统升级中", "-6", idcardNameInconsistentResults},
		{"idcard_name 超当日限制", "-7", idcardNameInconsistentResults},
		{"idcard_name 未知码", "-99", idcardNameInconsistentResults},
		{"idcard_phone 服务繁忙", "-11", idcardPhoneInconsistentResults},
		{"idcard_bankcard 服务繁忙", "-17", idcardBankcardInconsistentResults},
		{"name_bankcard 服务繁忙", "-16", nameBankcardInconsistentResults},
		{"name_bankcard 次数超限", "-17", nameBankcardInconsistentResults},
		{"idcard_bankcard 次数超限", "-18", idcardBankcardInconsistentResults},
	}
	for _, c := range errCases {
		ok, desc, _, err := judgeFaceVerifyResult(newResp(c.code), c.set)
		if ok || err == nil {
			t.Errorf("%s(%s): ok=%v, err=%v", c.name, c.code, ok, err)
			continue
		}
		if desc != "desc-"+c.code {
			t.Errorf("%s(%s): desc=%q, want %q", c.name, c.code, desc, "desc-"+c.code)
		}
		verr, isVerifyErr := err.(*errors.TencentCloudSDKError)
		if !isVerifyErr {
			t.Errorf("%s(%s): err type %T, want *errors.TencentCloudSDKError", c.name, c.code, err)
			continue
		}
		if verr.Code != c.code || verr.Message != "desc-"+c.code || verr.RequestId != "rid-"+c.code {
			t.Errorf("%s(%s): verr=%+v", c.name, c.code, verr)
		}
		// GetErrCodeAndErrMsg 应能统一提取
		if code, msg := GetErrCodeAndErrMsg(err); code != c.code || msg != "desc-"+c.code {
			t.Errorf("%s(%s): GetErrCodeAndErrMsg=(%q, %q)", c.name, c.code, code, msg)
		}
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
	ok, desc, requestId, err := txCloud.VerifyIdcardAndName(idcard, name)
	if err != nil {
		if e, ok := err.(*errors.TencentCloudSDKError); ok {
			t.Logf("e.Code is `%s`, e.Message is `%s`\n", e.Code, e.Message)
		} else {
			t.Logf("%s\n", err.Error())
		}
	} else if !ok {
		t.Errorf("(RequestId=%s) %s", requestId, desc)
	} else {
		t.Log(ok, desc)
	}
}
