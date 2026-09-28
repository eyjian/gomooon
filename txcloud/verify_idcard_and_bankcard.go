// Package txcloud
// Wrote by yijian on 2024/11/16
package txcloud

import (
	"encoding/json"
	"strings"
	"sync"
	"sync/atomic"
)
import (
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	faceid "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/faceid/v20180301"
)

// IdcardBankcardTuple 验证结果
// 过频错误示例：Code=RequestLimitExceeded, Message=您当前每秒请求 `103` 次，超过了每秒频率上限 `100`，请稍后重试。
type IdcardBankcardTuple struct {
	Idcard     string // 身份证号（如果证件号没有填写，则只验证 Name 和  Bankcard 的一致性）
	Name       string // 姓名
	Bankcard   string // 银行帐号/卡号
	Consistent bool   // 是否一致
	ErrDesc    string // 不一致原因描述
	RequestId  string
	Err        error // if txcloudErr, ok := err.(*errors.TencentCloudSDKError); ok {
}

// BatchVerifyIdcardAndBankcard 批量验证身份证号码和银行卡是否匹配
// 腾讯云频率限制：默认接口请求频率限制 20 次/秒
// 参数说明：
// concurrency 一次并发验证的个数（值应大于 0）
// failCount 接口调用失败数达到时中止执行，为 0 表示遇到一个错误即中止执行
// data 的 key 为身份证，调用者应保证 key 同 value 的一致性，传入的 Err 应为 nil，ErrDesc 为空，Consistent 为 false
// 因为 data 内的对象不会变，只是对象的内容变化，因此实现不需要加锁
// 遇到腾讯云接口调用失败数超过 failCount 时会中止执行，并在第三个返回值记录错误个数，此时对 data 的处理时不完整的，
// 可通过 data 的 Err 为 nil 和 ErrDesc 为空来判断，并且对应的 Consistent 必为 false
func (t *Face) BatchVerifyIdcardAndBankcard(concurrency, failCount int, data map[string]*IdcardBankcardTuple) (int, int, int) {
	var wg sync.WaitGroup

	consistent := int32(0)   // 一致的个数
	inconsistent := int32(0) // 不一致的个数
	fail := int32(0)         // 接口调用失败的个数
	semaphore := make(chan struct{}, concurrency)
	for _, v := range data {
		wg.Add(1)
		semaphore <- struct{}{} // 获取信号量，限制并发数量

		go func(v *IdcardBankcardTuple) {
			defer wg.Done()
			defer func() { <-semaphore }() // 释放信号量

			ok, desc, requestId, err := t.VerifyIdcardAndBankcard(v.Idcard, v.Name, v.Bankcard)
			if err != nil {
				v.Err = err
				v.Consistent = false
				v.ErrDesc = ""
				v.RequestId = requestId // 出错时也记录 RequestId，便于定位问题
			} else {
				v.Err = nil
				v.Consistent = ok
				v.ErrDesc = desc
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

// VerifyIdcardAndBankcard 验证身份证号码和银行卡是否匹配
// 腾讯云频率限制：默认接口请求频率限制 20 次/秒
// 返回值：一致性返回 true，否则返回 false，出错返回 error；第二个返回值为验证结果描述；第三个返回值为 RequestId
// 参数 idcard 和 bankcard 可含有空格
// Result 结果码分类（其余结果码及未知码一律返回 *errors.TencentCloudSDKError，Code 为 Result 结果码，防止将腾讯云侧故障误判为不一致）：
// 三要素（BankCardVerification）："0" 一致；"-1" 认证未通过；"-2"~"-4" 输入要素非法；"-5"~"-16"、"-19"、"-20" 卡状态/证件否定，均视为不一致
// 二要素（BankCard2EVerification，idcard 为空时）："0" 一致；"-1" 认证未通过；"-2"、"-3" 输入要素非法；"-4"~"-15"、"-18" 卡状态否定，均视为不一致
// 注意两个接口的码表是错开的（如 "-16" 在三要素为"发卡行不支持此交易"，在二要素为"验证中心服务繁忙"）
func (t *Face) VerifyIdcardAndBankcard(idcard, name, bankcard string) (bool, string, string, error) {
	var err error
	var jsonStr string
	inconsistentResults := idcardBankcardInconsistentResults
	client, _ := faceid.NewClient(t.credential, "", t.clientProfile)

	if idcard == "" {
		inconsistentResults = nameBankcardInconsistentResults
		var response *faceid.BankCard2EVerificationResponse

		request := faceid.NewBankCard2EVerificationRequest()
		request.Name = common.StringPtr(name)
		request.BankCard = common.StringPtr(strings.ReplaceAll(bankcard, " ", ""))

		response, err = client.BankCard2EVerification(request)
		if err == nil {
			jsonStr = response.ToJsonString()
		}
	} else {
		var response *faceid.BankCardVerificationResponse

		request := faceid.NewBankCardVerificationRequest()
		request.Name = common.StringPtr(name)
		request.BankCard = common.StringPtr(strings.ReplaceAll(bankcard, " ", ""))
		request.IdCard = common.StringPtr(strings.ReplaceAll(idcard, " ", ""))
		//request.CertType = common.Int64Ptr(0) // 赋值会大幅度降低通过率

		// 发起请求
		response, err = client.BankCardVerification(request)
		if err == nil {
			jsonStr = response.ToJsonString()
		}
	}
	// 出错时从 err 中提取 RequestId 一并返回，便于定位问题
	if err != nil {
		return false, "", extractRequestId(err), err
	}

	// 解析响应
	resp := FaceResponse{}
	err = json.Unmarshal([]byte(jsonStr), &resp)
	if err != nil {
		return false, "", "", err
	}
	// 判断是否一致（仅 "0" 为一致，确定性否定码判不一致，服务异常等其余码返回错误）
	return judgeFaceVerifyResult(&resp, inconsistentResults)
}
