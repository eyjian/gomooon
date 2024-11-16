// Package txcloud
// Wrote by yijian on 2024/11/16
package txcloud

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
)
import (
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	faceid "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/faceid/v20180301"
)

// IdcardBankcardTuple 验证结果
type IdcardBankcardTuple struct {
	Idcard string
	Name string
	Bankcard string
	Consistent bool // 是否一致
	ErrDesc string // 不一致原因描述
	Err error
}

// BatchVerifyIdcardAndBankcard 批量验证身份证号码和银行卡是否匹配
// 参数说明：
// concurrency 一次并发验证的个数
// data 的 key 为身份证
func (t *Face) BatchVerifyIdcardAndBankcard(concurrency int, data map[string]*IdcardBankcardTuple) (int, int, int) {
	var wg sync.WaitGroup
	var mutex sync.Mutex

	consistent := 0 // 一致的个数
	inconsistent := 0 // 不一致的个数
	fail := 0 // 出错的个数
	semaphore := make(chan struct{}, concurrency)
	for _, v := range data {
		wg.Add(1)
		semaphore <- struct{}{} // 获取信号量，限制并发数量

		go func(v *IdcardBankcardTuple) {
			defer wg.Done()
			defer func() { <-semaphore }() // 释放信号量

			ok, desc, err := t.VerifyIdcardAndBankcard(v.Idcard, v.Name, v.Bankcard)
			if err != nil {
				v.Err = err
				v.Consistent = false
				v.ErrDesc = ""
			} else {
				v.Err = nil
				v.Consistent = ok
				v.ErrDesc = desc
			}

			mutex.Lock()
			defer mutex.Unlock()
			data[v.Idcard] = v

			if err != nil {
				fail++
			} else if ok {
				consistent++
			} else {
				inconsistent++
			}
		}(v)
	}

	wg.Wait()
	return consistent, inconsistent, fail
}

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
