// Package txcloud
// Wrote by yijian on 2026/06/02
package txcloud

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	ses "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ses/v20201002"
)

// EmailSender 邮件发送接口
//
// 调用方应当依赖该接口而非具体实现，以便后续平滑切换到其它邮件方案。
type EmailSender interface {
	// SendEmail 发送一封邮件
	//
	// 当腾讯云返回 *errors.TencentCloudSDKError 错误时，原样向上透传，
	// 调用方可使用 GetErrCodeAndErrMsg 提取错误码和错误描述。
	SendEmail(ctx context.Context, req *EmailRequest) (*EmailResult, error)
}

// EmailRequest 邮件发送请求
type EmailRequest struct {
	// FromEmailAddress 发件人邮箱地址
	// 不使用别名时直接填写邮箱地址，例如：noreply@mail.qcloud.com
	// 如需填写发件人别名，格式为：别名 <邮箱地址>，例如：运维团队 <noreply@mail.qcloud.com>
	FromEmailAddress string

	// Subject 邮件主题
	Subject string

	// Destination 收信人邮箱地址列表，最多支持群发 50 人
	Destination []string

	// ReplyToAddresses 回复邮箱地址，可选
	// 如果不填，收件人的回复邮件将会发送失败
	ReplyToAddresses string

	// Cc 抄送人邮箱地址列表，最多 20 人，可选
	Cc []string

	// Bcc 密送人邮箱地址列表，最多 20 人，可选
	Bcc []string

	// Html 邮件 HTML 内容（二选一，优先于 Text）
	Html string

	// Text 邮件纯文本内容（二选一）
	Text string

	// TemplateID 模板 ID，使用模板发送时填写
	// 如果设置了 TemplateID，则忽略 Html 和 Text
	TemplateID uint64

	// TemplateData 模板变量参数，JSON 字符串
	// 例如：{"name":"xxx","code":"1234"}
	TemplateData string

	// Attachments 附件列表，可选
	Attachments []EmailAttachment

	// TriggerType 邮件触发类型
	// 0: 非触发类（默认），营销类邮件、非即时类邮件等
	// 1: 触发类，验证码等即时发送类邮件
	TriggerType uint64

	// Unsubscribe 退订链接选项，可选
	// 0: 不加入退订链接 1: 简体中文 2: 英文 3: 繁体中文
	Unsubscribe string
}

// EmailAttachment 邮件附件
type EmailAttachment struct {
	// FileName 附件名称
	FileName string
	// Content 附件内容（Base64 编码），如果为空则从 FilePath 读取
	Content string
	// FilePath 本地文件路径，当 Content 为空时从此路径读取文件
	FilePath string
}

// EmailResult 邮件发送结果
type EmailResult struct {
	RequestId string // 腾讯云请求 ID
	MessageId string // 邮件唯一消息标识符
}

// 默认 region
const defaultEmailRegion = "ap-hongkong"

// SesEmailSender 基于腾讯云 SES 的邮件发送实现
type SesEmailSender struct {
	*TxCloud // 继承 Region 等字段
}

// 编译期断言：确保实现满足 EmailSender 接口
var _ EmailSender = (*SesEmailSender)(nil)

// NewSesEmailSender 实例化腾讯云 SES 邮件发送器
//
// 参数：
//
//	secretId / secretKey 腾讯云访问密钥
//	region              腾讯云地域，传空字符串将使用默认值 "ap-hongkong"
func NewSesEmailSender(secretId, secretKey, region string) *SesEmailSender {
	if region == "" {
		region = defaultEmailRegion
	}
	tx := NewTxCloud(secretId, secretKey, "ses.tencentcloudapi.com")
	tx.Region = region
	return &SesEmailSender{
		TxCloud: tx,
	}
}

// SendEmail 发送邮件（实现 EmailSender 接口）
func (s *SesEmailSender) SendEmail(ctx context.Context, req *EmailRequest) (*EmailResult, error) {
	if req == nil {
		return nil, fmt.Errorf("email request is nil")
	}
	if len(req.Destination) == 0 && len(req.Cc) == 0 && len(req.Bcc) == 0 {
		return nil, fmt.Errorf("email request: at least one of Destination/Cc/Bcc must be provided")
	}
	if req.FromEmailAddress == "" {
		return nil, fmt.Errorf("email request: FromEmailAddress is required")
	}
	if req.Subject == "" {
		return nil, fmt.Errorf("email request: Subject is required")
	}

	client, err := ses.NewClient(s.credential, s.TxCloud.Region, s.clientProfile)
	if err != nil {
		return nil, fmt.Errorf("failed to create ses client: %w", err)
	}

	request := ses.NewSendEmailRequest()
	if ctx != nil {
		request.SetContext(ctx)
	}

	request.FromEmailAddress = common.StringPtr(req.FromEmailAddress)
	request.Subject = common.StringPtr(req.Subject)
	request.Destination = common.StringPtrs(req.Destination)

	if req.ReplyToAddresses != "" {
		request.ReplyToAddresses = common.StringPtr(req.ReplyToAddresses)
	}
	if len(req.Cc) > 0 {
		request.Cc = common.StringPtrs(req.Cc)
	}
	if len(req.Bcc) > 0 {
		request.Bcc = common.StringPtrs(req.Bcc)
	}
	if req.TriggerType > 0 {
		request.TriggerType = &req.TriggerType
	}
	if req.Unsubscribe != "" {
		request.Unsubscribe = common.StringPtr(req.Unsubscribe)
	}

	// 邮件内容：模板 > Html/Text
	if req.TemplateID > 0 {
		request.Template = &ses.Template{
			TemplateID:   &req.TemplateID,
			TemplateData: common.StringPtr(req.TemplateData),
		}
	} else if req.Html != "" || req.Text != "" {
		htmlContent := req.Html
		textContent := req.Text
		// 如果 Html 以 "file://" 开头，则从文件读取
		if strings.HasPrefix(htmlContent, "file://") {
			data, err := os.ReadFile(strings.TrimPrefix(htmlContent, "file://"))
			if err != nil {
				return nil, fmt.Errorf("failed to read html file: %w", err)
			}
			htmlContent = string(data)
		}
		// 如果 Text 以 "file://" 开头，则从文件读取
		if strings.HasPrefix(textContent, "file://") {
			data, err := os.ReadFile(strings.TrimPrefix(textContent, "file://"))
			if err != nil {
				return nil, fmt.Errorf("failed to read text file: %w", err)
			}
			textContent = string(data)
		}
		request.Simple = &ses.Simple{
			Html: common.StringPtr(base64.StdEncoding.EncodeToString([]byte(htmlContent))),
			Text: common.StringPtr(base64.StdEncoding.EncodeToString([]byte(textContent))),
		}
	} else {
		return nil, fmt.Errorf("email request: one of TemplateID, Html or Text must be provided")
	}

	// 附件
	if len(req.Attachments) > 0 {
		attachments := make([]*ses.Attachment, 0, len(req.Attachments))
		for i := range req.Attachments {
			att := &req.Attachments[i]
			content := att.Content
			if content == "" && att.FilePath != "" {
				data, err := os.ReadFile(att.FilePath)
				if err != nil {
					return nil, fmt.Errorf("failed to read attachment file %q: %w", att.FilePath, err)
				}
				content = base64.StdEncoding.EncodeToString(data)
			}
			if content == "" {
				continue
			}
			attachments = append(attachments, &ses.Attachment{
				FileName: common.StringPtr(att.FileName),
				Content:  common.StringPtr(content),
			})
		}
		request.Attachments = attachments
	}

	response, err := client.SendEmail(request)
	if err != nil {
		return nil, err
	}

	result := &EmailResult{}
	if response.Response != nil {
		if response.Response.RequestId != nil {
			result.RequestId = *response.Response.RequestId
		}
		if response.Response.MessageId != nil {
			result.MessageId = *response.Response.MessageId
		}
	}
	return result, nil
}
