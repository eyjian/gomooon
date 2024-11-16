// Package txcloud
// Wrote by yijian on 2024/11/16
package txcloud

import (
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
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

// NewTxCloud 创建腾讯云对象
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
