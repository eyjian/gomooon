// Package txcloud
// Wrote by yijian on 2026/08/31
package txcloud

import (
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
)
import (
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	ocr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ocr/v20181119"
)

// ocrEndpoint OCR 服务的接入点
const ocrEndpoint = "ocr.tencentcloudapi.com"

// EntVerifier 企业核验器（基于腾讯云 OCR），用于企业名称、统一社会信用代码等要素核验
type EntVerifier struct {
	*TxCloud // 复用 credential、clientProfile 等
}

// NewEntVerifier 创建企业核验器
// 密钥在控制台 https://console.cloud.tencent.com/cam/capi 上获取
func NewEntVerifier(secretId, secretKey string) *EntVerifier {
	return &EntVerifier{
		TxCloud: NewTxCloud(secretId, secretKey, ocrEndpoint),
	}
}

// operatingStatusDesc 企业状态码（OperatingStatus）同描述的映射
var operatingStatusDesc = map[string]string{
	"1": "开业（在营）/ 在营",
	"2": "迁出 / 非在营",
	"3": "注销",
	"4": "吊销",
	"5": "撤销",
	"6": "停业",
	"7": "撤销登记",
	"0": "其他",
	"/": "无法查询",
}

// OperatingStatusDesc 企业状态码（OperatingStatus）转描述，未知的状态码原样返回
// 枚举值：1：开业（在营）/ 在营；2：迁出 / 非在营；3：注销；4：吊销；5：撤销；6：停业；7：撤销登记；0：其他；/：无法查询
func OperatingStatusDesc(operatingStatus string) string {
	if desc, ok := operatingStatusDesc[operatingStatus]; ok {
		return desc
	} else {
		return operatingStatus
	}
}

// EntNameCreditCodeTuple 验证结果
// 营业执照核验（增强版）接口 VerifyBizLicenseEnterprise3，VerifyType 取值 ENT_2META（企业名称、统一社会信用代码二要素核验）
// 响应的 StatusCode：0 成功（计费），1 系统异常（不计费），2 查询无结果（不计费）
type EntNameCreditCodeTuple struct {
	EntName         string // 企业名称（禁止简称，标点符号需完全一致）
	CreditCode      string // 统一社会信用代码（18 位，字母全部大写，不含空格和横杠）
	Consistent      bool   // 是否一致
	OperatingStatus string // 企业状态码（描述见 OperatingStatusDesc），查询无结果时为空串
	RequestId       string
	Err             error // if txcloudErr, ok := err.(*errors.TencentCloudSDKError); ok {
}

// BatchVerifyEntNameAndCreditCode 批量验证企业名称和统一社会信用代码是否匹配
// 腾讯云频率限制：默认接口请求频率限制 20 次/秒
// 参数说明：
// concurrency 一次并发验证的个数（值应大于 0）
// failCount 接口调用失败数达到时中止执行，为 0 表示遇到一个错误即中止执行
// data 的 key 为统一社会信用代码，调用者应保证 key 同 value 的一致性，传入的 Err 应为 nil，Consistent 为 false
// 因为 data 内的对象不会变，只是对象的内容变化，因此实现不需要加锁
// 遇到腾讯云接口调用失败数超过 failCount 时会中止执行，并在第三个返回值记录错误个数，此时对 data 的处理时不完整的，
// 可通过 data 的 Err 为 nil 来判断，并且对应的 Consistent 必为 false
func (t *EntVerifier) BatchVerifyEntNameAndCreditCode(concurrency, failCount int, data map[string]*EntNameCreditCodeTuple) (int, int, int) {
	var wg sync.WaitGroup

	consistent := int32(0)   // 一致的个数
	inconsistent := int32(0) // 不一致的个数
	fail := int32(0)         //接口调用失败的个数
	semaphore := make(chan struct{}, concurrency)
	for _, v := range data {
		wg.Add(1)
		semaphore <- struct{}{} // 获取信号量，限制并发数量

		go func(v *EntNameCreditCodeTuple) {
			defer wg.Done()
			defer func() { <-semaphore }() // 释放信号量

			ok, operatingStatus, requestId, err := t.VerifyEntNameAndCreditCode(v.EntName, v.CreditCode)
			if err != nil {
				v.Err = err
				v.Consistent = false
				v.OperatingStatus = ""
			} else {
				v.Err = nil
				v.Consistent = ok
				v.OperatingStatus = operatingStatus
				v.RequestId = requestId
			}
			if err != nil {
				atomic.AddInt32(&fail, 1)
			} else if ok {
				atomic.AddInt32(&consistent, 1)
			} else {
				atomic.AddInt32(&inconsistent, 1)
			}
		}(v)

		// 如果出错，直接退出循环，特别是过频错误时及时释放额度
		if atomic.LoadInt32(&fail) > int32(failCount) {
			break
		}
	}

	wg.Wait()
	return int(consistent), int(inconsistent), int(fail)
}

// VerifyEntNameAndCreditCode 验证企业名称和统一社会信用代码是否匹配（营业执照二要素核验）
// 腾讯云频率限制：默认接口请求频率限制 20 次/秒
// 返回值：一致返回 true，否则返回 false，出错返回 error；第二个返回值为企业状态码；第三个返回值为 RequestId
// 参数 creditCode 可含有空格和横杠，字母会被转为大写
// StatusCode 为 1（系统异常，不计费）时返回 error，可重试；为 2（查询无结果，不计费）时视为不一致
func (t *EntVerifier) VerifyEntNameAndCreditCode(entName, creditCode string) (bool, string, string, error) {
	client, err := ocr.NewClient(t.credential, t.Region, t.clientProfile)
	if err != nil {
		return false, "", "", err
	}
	request := ocr.NewVerifyBizLicenseEnterprise3Request()
	request.CreditCode = common.StringPtr(normalizeCreditCode(creditCode))
	request.EntName = common.StringPtr(entName)
	request.VerifyType = common.StringPtr("ENT_2META") // 企业名称、统一社会信用代码二要素核验

	// 发起请求（腾讯云 SDK 错误 *errors.TencentCloudSDKError 原样透传，调用方可据此判断错误码）
	response, err := client.VerifyBizLicenseEnterprise3(request)
	if err != nil {
		return false, "", "", err
	}
	if response == nil || response.Response == nil {
		return false, "", "", fmt.Errorf("empty response")
	}

	// RequestId 每次请求都会返回，这里做防御性处理，避免空指针
	requestId := ""
	if response.Response.RequestId != nil {
		requestId = *response.Response.RequestId
	}

	// StatusCode：0 成功（计费），1 系统异常（不计费），2 查询无结果（不计费）
	if response.Response.StatusCode == nil {
		return false, "", requestId, fmt.Errorf("response without StatusCode")
	}
	switch *response.Response.StatusCode {
	case 1: // 系统异常，不计费，可重试
		return false, "", requestId, fmt.Errorf("txcloud system error (StatusCode=1)")
	case 2: // 查询无结果，不计费，视为不一致
		return false, "", requestId, nil
	}

	// 仅 StatusCode 为 0 时才返回下列字段
	operatingStatus := ""
	if response.Response.OperatingStatus != nil {
		operatingStatus = *response.Response.OperatingStatus
	}

	// VerifyResult：1 完全匹配，0 不完全匹配（仅 StatusCode 为 0 时返回）
	if response.Response.VerifyResult == nil || *response.Response.VerifyResult != 1 {
		return false, operatingStatus, requestId, nil // 不一致
	}
	return true, operatingStatus, requestId, nil // 一致
}

// normalizeCreditCode 规范统一社会信用代码：去掉空格和横杠，字母转为大写
func normalizeCreditCode(creditCode string) string {
	creditCode = strings.ReplaceAll(creditCode, " ", "")
	creditCode = strings.ReplaceAll(creditCode, "-", "")
	return strings.ToUpper(creditCode)
}
