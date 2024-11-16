// Package txcloud
// Wrote by yijian on 2024/11/16
package txcloud

import (
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
)

// TxCloud 腾讯云
type TxCloud struct {
	SecretId string
	SecretKey string
	Endpoint string // 示例："faceid.tencentcloudapi.com"
	Region string // 核验身份证和银行卡不用指定，取值如：ap-shanghai

	credential *common.Credential
	clientProfile *profile.ClientProfile
}

// FaceResponse 腾讯云 face 类接口的响应
//{
//	"Response": {
//		"Description": "非法姓名（长度、格式等不正确）",
//		"RequestId": "2ee9df1f-23b8-457d-89ba-337a9c57e4b9",
//		"Result": "-3"
//	}
//}
type FaceResponse struct {
	Response struct {
		Description string `json:"Description"`
		RequestId   string `json:"RequestId"`
		Result      string `json:"Result"`
	} `json:"Response"`
}

// GetErrCodeAndErrMsg 获取腾讯云错误码和错误信息，如果不是腾讯云错则返回两个空字符串
func GetErrCodeAndErrMsg(err error) (string, string) {
	if e, ok := err.(*errors.TencentCloudSDKError); ok {
		return e.Code, e.Message
	} else {
		return "", ""
	}
}

// NewTxCloud 创建腾讯云对象
// 密钥在控制台 https://console.cloud.tencent.com/cam/capi 上获取
func NewTxCloud(secretId, secretKey, endpoint string) *TxCloud {
	txCloud:= &TxCloud{
		SecretId: secretId,
		SecretKey: secretKey,
		Endpoint: endpoint,
		credential: common.NewCredential(
			secretId,
			secretKey,
		),
		clientProfile: profile.NewClientProfile(),
	}
	txCloud.clientProfile.HttpProfile.Endpoint = endpoint

	return txCloud
}
