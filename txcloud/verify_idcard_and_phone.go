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
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	faceid "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/faceid/v20180301"
)

// IdcardPhoneTuple 验证结果
type IdcardPhoneTuple struct {
	Idcard string
	Name string
	Phone string
	Consistent bool // 是否一致
	ErrDesc string // 不一致原因描述
	Err error
}

// BatchVerifyIdcardAndPhone 批量验证身份证号码和手机号是否匹配
// 腾讯云频率限制：默认接口请求频率限制 20 次/秒
// 参数说明：
// concurrency 一次并发验证的个数（值应大于 0）
// data 的 key 为身份证，调用者应保证 key 同 value 的一致性
// 因为 data 内的对象不会变，只是对象的内容变化，因此实现不需要加锁
func (t *Face) BatchVerifyIdcardAndPhone(concurrency int, data map[string]*IdcardPhoneTuple) (int, int, int) {
	var wg sync.WaitGroup

	consistent := int32(0) // 一致的个数
	inconsistent := int32(0) // 不一致的个数
	fail := int32(0) // 出错的个数
	semaphore := make(chan struct{}, concurrency)
	for _, v := range data {
		wg.Add(1)
		semaphore <- struct{}{} // 获取信号量，限制并发数量

		go func(v *IdcardPhoneTuple) {
			defer wg.Done()
			defer func() { <-semaphore }() // 释放信号量

			ok, desc, err := t.VerifyIdcardAndPhone(v.Idcard, v.Name, v.Phone)
			if err != nil {
				v.Err = err
				v.Consistent = false
				v.ErrDesc = ""
			} else {
				v.Err = nil
				v.Consistent = ok
				v.ErrDesc = desc
			}
			if err != nil {
				atomic.AddInt32(&fail, 1)
			} else if ok {
				atomic.AddInt32(&consistent, 1)
			} else {
				atomic.AddInt32(&inconsistent, 1)
			}
		}(v)
	}

	wg.Wait()
	return int(consistent), int(inconsistent), int(fail)
}

// VerifyIdcardAndPhone 验证身份证号码和手机号是否匹配
// 腾讯云频率限制：默认接口请求频率限制 20 次/秒
// 返回值：一致性返回 true，否则返回 false，出错返回 error；第二个返回值为验证结果描述
// 参数 idcard 和 bankcard 可含有空格
func (t *Face) VerifyIdcardAndPhone(idcard, name, phone string) (bool, string, error) {
	client, _ := faceid.NewClient(t.credential, "", t.clientProfile)
	request := faceid.NewPhoneVerificationRequest()
	request.IdCard = common.StringPtr(strings.ReplaceAll(idcard, " ", ""))
	request.Name = common.StringPtr(name)
	request.Phone = common.StringPtr(strings.ReplaceAll(phone, " ", ""))

	// 发起请求
	response, err := client.PhoneVerification(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		//return false, "", fmt.Errorf("a txcloud API error has returned: %s", err.Error())
		return false, "", err
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
