// Package txcloud
// Wrote by yijian on 2026/06/02
package txcloud

import (
	"context"
	"os"
	"strconv"
	"testing"
)

// 测试前请设置以下环境变量：
//   TENCENTCLOUD_SECRET_ID     - 腾讯云 SecretId
//   TENCENTCLOUD_SECRET_KEY    - 腾讯云 SecretKey
//   TENCENTCLOUD_SES_REGION    - 腾讯云 SES 地域（可选，默认 ap-hongkong）
//   TENCENTCLOUD_SES_FROM      - 发件人邮箱地址
//   TENCENTCLOUD_SES_TO        - 收件人邮箱地址
//   TENCENTCLOUD_SES_TEMPLATE_ID - 模板 ID（可选，用于模板发送测试）
//
// go test -v -run=TestSendEmail ./txcloud/

// TestSendEmailSimple 测试发送纯文本邮件
// go test -v -run="TestSendEmailSimple"
func TestSendEmailSimple(t *testing.T) {
	secretId := os.Getenv("TENCENTCLOUD_SECRET_ID")
	secretKey := os.Getenv("TENCENTCLOUD_SECRET_KEY")
	from := os.Getenv("TENCENTCLOUD_SES_FROM")
	to := os.Getenv("TENCENTCLOUD_SES_TO")
	region := os.Getenv("TENCENTCLOUD_SES_REGION")

	if secretId == "" || secretKey == "" {
		t.Skip("TENCENTCLOUD_SECRET_ID or TENCENTCLOUD_SECRET_KEY not set")
	}
	if from == "" || to == "" {
		t.Skip("TENCENTCLOUD_SES_FROM or TENCENTCLOUD_SES_TO not set")
	}

	sender := NewSesEmailSender(secretId, secretKey, region)
	result, err := sender.SendEmail(context.Background(), &EmailRequest{
		FromEmailAddress: from,
		Subject:          "gomooon 邮件测试 - 纯文本",
		Destination:      []string{to},
		Text:             "这是一封由 gomooon txcloud 包发送的测试邮件（纯文本）。",
		TriggerType:      1, // 触发类
	})
	if err != nil {
		t.Fatalf("SendEmail failed: %v", err)
	}
	t.Logf("SendEmail success: RequestId=%s, MessageId=%s", result.RequestId, result.MessageId)
}

// TestSendEmailHTML 测试发送 HTML 邮件
func TestSendEmailHTML(t *testing.T) {
	secretId := os.Getenv("TENCENTCLOUD_SECRET_ID")
	secretKey := os.Getenv("TENCENTCLOUD_SECRET_KEY")
	from := os.Getenv("TENCENTCLOUD_SES_FROM")
	to := os.Getenv("TENCENTCLOUD_SES_TO")
	region := os.Getenv("TENCENTCLOUD_SES_REGION")

	if secretId == "" || secretKey == "" {
		t.Skip("TENCENTCLOUD_SECRET_ID or TENCENTCLOUD_SECRET_KEY not set")
	}
	if from == "" || to == "" {
		t.Skip("TENCENTCLOUD_SES_FROM or TENCENTCLOUD_SES_TO not set")
	}

	html := `<html><body><h1 style="color:blue;">邮件测试</h1><p>这是一封由 <b>gomooon txcloud</b> 包发送的 HTML 测试邮件。</p></body></html>`

	sender := NewSesEmailSender(secretId, secretKey, region)
	result, err := sender.SendEmail(context.Background(), &EmailRequest{
		FromEmailAddress: from,
		Subject:          "gomooon 邮件测试 - HTML",
		Destination:      []string{to},
		Html:             html,
		TriggerType:      1,
	})
	if err != nil {
		t.Fatalf("SendEmail failed: %v", err)
	}
	t.Logf("SendEmail success: RequestId=%s, MessageId=%s", result.RequestId, result.MessageId)
}

// TestSendEmailTemplate 测试使用模板发送邮件
// go test -v -run="TestSendEmailTemplate"
func TestSendEmailTemplate(t *testing.T) {
	secretId := os.Getenv("TENCENTCLOUD_SECRET_ID")
	secretKey := os.Getenv("TENCENTCLOUD_SECRET_KEY")
	from := os.Getenv("TENCENTCLOUD_SES_FROM")
	to := os.Getenv("TENCENTCLOUD_SES_TO")
	region := os.Getenv("TENCENTCLOUD_SES_REGION")
	templateIDStr := os.Getenv("TENCENTCLOUD_SES_TEMPLATE_ID")

	if secretId == "" || secretKey == "" {
		t.Skip("TENCENTCLOUD_SECRET_ID or TENCENTCLOUD_SECRET_KEY not set")
	}
	if from == "" || to == "" {
		t.Skip("TENCENTCLOUD_SES_FROM or TENCENTCLOUD_SES_TO not set")
	}
	if templateIDStr == "" {
		t.Skip("TENCENTCLOUD_SES_TEMPLATE_ID not set")
	}

	templateID, err := strconv.ParseUint(templateIDStr, 10, 64)
	if err != nil {
		t.Fatalf("invalid TENCENTCLOUD_SES_TEMPLATE_ID: %v", err)
	}

	sender := NewSesEmailSender(secretId, secretKey, region)
	result, err := sender.SendEmail(context.Background(), &EmailRequest{
		FromEmailAddress: from,
		Subject:          "gomooon 邮件测试 - 模板",
		Destination:      []string{to},
		TemplateID:       templateID,
		TemplateData:     `{"code":"123456"}`,
		TriggerType:      1,
	})
	if err != nil {
		t.Fatalf("SendEmail failed: %v", err)
	}
	t.Logf("SendEmail success: RequestId=%s, MessageId=%s", result.RequestId, result.MessageId)
}
