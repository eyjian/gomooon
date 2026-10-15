// Package txcloud
// Wrote by yijian on 2024/11/16
package txcloud

type Face struct {
	*TxCloud
}

// NewFace 实例化腾讯云人脸识别
//
// endpoint 解析优先级：WithEndpoint option > 环境变量 TXCLOUD_FACEID_ENDPOINT > 默认值 faceid.tencentcloudapi.com
func NewFace(secretId, secretKey string, opts ...Option) *Face {
	return &Face{
		TxCloud: NewTxCloud(secretId, secretKey, resolveEndpoint(envFaceidEndpoint, defaultFaceidEndpoint, opts)),
	}
}
