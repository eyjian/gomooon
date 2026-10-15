// Package txcloud
// Wrote by yijian on 2024/11/16
package txcloud

import (
	"os"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
)

// 各服务的默认 endpoint
const (
	defaultFaceidEndpoint = "faceid.tencentcloudapi.com"
	defaultSesEndpoint    = "ses.tencentcloudapi.com"
	defaultOcrEndpoint    = "ocr.tencentcloudapi.com"
)

// 各服务 endpoint 对应的环境变量名
const (
	envFaceidEndpoint = "TXCLOUD_FACEID_ENDPOINT"
	envSesEndpoint    = "TXCLOUD_SES_ENDPOINT"
	envOcrEndpoint    = "TXCLOUD_OCR_ENDPOINT"
)

// Option 腾讯云可选配置，用于各构造函数的可变参数
type Option func(*TxCloud)

// WithEndpoint 显式指定 endpoint，优先级高于环境变量和默认值
func WithEndpoint(endpoint string) Option {
	return func(t *TxCloud) {
		t.Endpoint = endpoint
	}
}

// resolveEndpoint 按优先级解析 endpoint：WithEndpoint option > 环境变量 > 默认值
func resolveEndpoint(envKey, defaultEndpoint string, opts []Option) string {
	txCloud := &TxCloud{}
	for _, opt := range opts {
		if opt != nil {
			opt(txCloud)
		}
	}
	if txCloud.Endpoint != "" {
		return txCloud.Endpoint
	}
	if endpoint := os.Getenv(envKey); endpoint != "" {
		return endpoint
	}
	return defaultEndpoint
}

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

// 各接口视为"不一致"的 Result 结果码集合：这些码均为确定性否定（用户输入问题或权威库明确无记录），重试结果不变，
// 判不一致可引导用户检查输入；不在集合内且非 "0"（一致）的结果码（如证件库服务异常、验证中心服务繁忙、次数超限，
// 以及未知的未来新增结果码）一律以 *errors.TencentCloudSDKError 返回兜底（Code 为 Result 结果码），
// 防止将腾讯云侧故障误判为不一致，且调用方可沿用 GetErrCodeAndErrMsg 统一处理

// idcardNameInconsistentResults 身份证二要素核验（IdCardVerification）：
// -1 不一致；-2 非法身份证号；-3 非法姓名；-5 证件库中无此身份证记录
var idcardNameInconsistentResults = map[string]struct{}{"-1": {}, "-2": {}, "-3": {}, "-5": {}}

// idcardPhoneInconsistentResults 手机号三要素核验（PhoneVerification）：
// -4 三要素信息不一致；-6 手机号码不合法；-7 身份证号码有误；-8 姓名校验不通过；-9 没有记录；
// -13 该号段不支持验证（确定性结果，重试不变，判不一致以引导用户换号而非无限重试）
var idcardPhoneInconsistentResults = map[string]struct{}{"-4": {}, "-6": {}, "-7": {}, "-8": {}, "-9": {}, "-13": {}}

// idcardBankcardInconsistentResults 银行卡三要素核验（BankCardVerification）：
// -1 认证未通过；-2 姓名校验不通过；-3 身份证号码有误；-4 银行卡号格式有误；
// -5~-16、-20 均为收费的确定性卡状态否定（持卡人信息有误、无效卡号、挂失、过期等）；
// -19 该证件号暂不支持核验（确定性结果，重试不变，判不一致以引导用户改用二代身份证）；
// -17 验证中心服务繁忙、-18 验证次数超限 返回错误
var idcardBankcardInconsistentResults = map[string]struct{}{
	"-1": {}, "-2": {}, "-3": {}, "-4": {}, "-5": {}, "-6": {}, "-7": {}, "-8": {},
	"-9": {}, "-10": {}, "-11": {}, "-12": {}, "-13": {}, "-14": {}, "-15": {}, "-16": {}, "-19": {}, "-20": {},
}

// nameBankcardInconsistentResults 银行卡二要素核验（BankCard2EVerification），注意码表同三要素错开，不可混用：
// -1 认证未通过；-2 姓名校验不通过；-3 银行卡号格式有误；
// -4~-15、-18 均为收费的确定性卡状态否定（持卡人信息有误、无效卡号、挂失、过期、卡状态异常等）；
// -16 验证中心服务繁忙、-17 验证次数超限 返回错误
var nameBankcardInconsistentResults = map[string]struct{}{
	"-1": {}, "-2": {}, "-3": {}, "-4": {}, "-5": {}, "-6": {}, "-7": {}, "-8": {},
	"-9": {}, "-10": {}, "-11": {}, "-12": {}, "-13": {}, "-14": {}, "-15": {}, "-18": {},
}

// extractRequestId 从错误中提取腾讯云 RequestId，非腾讯云错误时返回空串
func extractRequestId(err error) string {
	if e, ok := err.(*errors.TencentCloudSDKError); ok {
		return e.RequestId
	}
	return ""
}

// judgeFaceVerifyResult 判定 face 类核验接口结果：
// Result 为 "0" 返回一致；在 inconsistentResults 中返回不一致；
// 其余返回 *errors.TencentCloudSDKError（Code 为 Result 结果码，Message 为 Description），
// 出错时第二返回值仍返回腾讯云原始描述，便于调用方直接打日志
func judgeFaceVerifyResult(resp *FaceResponse, inconsistentResults map[string]struct{}) (bool, string, string, error) {
	if resp.Response.Result == "0" {
		return true, resp.Response.Description, resp.Response.RequestId, nil // 一致
	}
	if _, ok := inconsistentResults[resp.Response.Result]; ok {
		return false, resp.Response.Description, resp.Response.RequestId, nil // 不一致
	}
	return false, resp.Response.Description, resp.Response.RequestId,
		errors.NewTencentCloudSDKError(resp.Response.Result, resp.Response.Description, resp.Response.RequestId)
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
