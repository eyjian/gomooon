// Package txcloud
// Wrote by yijian on 2024/11/16
package txcloud

type Face struct {
	*TxCloud
}

// NewFace 实例化腾讯云人脸识别
func NewFace(secretId, secretKey string) *Face {
	return &Face{
		TxCloud: NewTxCloud(secretId, secretKey, "faceid.tencentcloudapi.com"),
	}
}
