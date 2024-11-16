// Package txcloud
// Wrote by yijian on 2024/11/16
package txcloud

import (
	"encoding/json"
	"fmt"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	faceid "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/faceid/v20180301"
)

// Response 腾讯云接口响应
//{
//	"Response": {
//		"Description": "非法姓名（长度、格式等不正确）",
//		"RequestId": "2ee9df1f-23b8-457d-89ba-337a9c57e4b9",
//		"Result": "-3"
//	}
//}
type Response struct {
	Response struct {
		Description string `json:"Description"`
		RequestId   string `json:"RequestId"`
		Result      string `json:"Result"`
	} `json:"Response"`
}

// VerifyIdcardAndName 验证身份证号码和姓名是否匹配
// 返回值：一致性返回 true，否则返回 false，出错返回 error；第二个返回值为验证结果描述
func (t *TxCloud) VerifyIdcardAndName(idcard, name string) (bool, string, error) {
	client, _ := faceid.NewClient(t.credential, "", t.clientProfile)
	request := faceid.NewIdCardVerificationRequest()
	request.IdCard = common.StringPtr(idcard)
	request.Name = common.StringPtr(name)

	// 发起请求
	response, err := client.IdCardVerification(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return false, "", fmt.Errorf("a txcloud API error has returned: %s", err.Error())
	}
	if err != nil {
		return false, "", err
	}

	// 解析响应
	resp := Response{}
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
