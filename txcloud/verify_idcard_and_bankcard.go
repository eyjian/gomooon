// Package txcloud
// Wrote by yijian on 2024/11/16
package txcloud

import (
	"encoding/json"
	"fmt"
	"strings"
)
import (
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	faceid "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/faceid/v20180301"
)

// VerifyIdcardAndBankcard 验证身份证号码和银行卡是否匹配
// 返回值：一致性返回 true，否则返回 false，出错返回 error；第二个返回值为验证结果描述
// 参数 idcard 和 bankcard 可含有空格
func (t *Face) VerifyIdcardAndBankcard(idcard, name, bankcard string) (bool, string, error) {
	client, _ := faceid.NewClient(t.credential, "", t.clientProfile)
	request := faceid.NewBankCardVerificationRequest()
	request.IdCard = common.StringPtr(strings.ReplaceAll(idcard, " ", ""))
	request.Name = common.StringPtr(name)
	request.BankCard = common.StringPtr(strings.ReplaceAll(bankcard, " ", ""))
	request.CertType = common.Int64Ptr(0)

	// 发起请求
	response, err := client.BankCardVerification(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return false, "", fmt.Errorf("a txcloud API error has returned: %s", err.Error())
	}
	if err != nil {
		return false, "", err
	}

	// 解析响应
	resp := FaceResponse{}
	jsonStr := response.ToJsonString()
	err = json.Unmarshal([]byte(jsonStr), &resp)
	if err != nil {
		return false, "", err
	}

	// 判断是否一致
	if resp.Response.Result != "0" {
		return false, resp.Response.Description, nil // 不一致
	}
	return true, resp.Response.Description, nil // 一致
}
