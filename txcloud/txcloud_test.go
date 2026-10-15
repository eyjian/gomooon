// Package txcloud
// Wrote by yijian on 2026/10/15
package txcloud

import "testing"

// TestResolveEndpoint 验证 endpoint 解析优先级：WithEndpoint option > 环境变量 > 默认值
func TestResolveEndpoint(t *testing.T) {
	const (
		envKey          = "TXCLOUD_TEST_ENDPOINT"
		defaultEndpoint = "default.tencentcloudapi.com"
		envEndpoint     = "env.tencentcloudapi.com"
		optionEndpoint  = "option.tencentcloudapi.com"
	)

	// 默认值兜底：无 option、无环境变量
	if got := resolveEndpoint(envKey, defaultEndpoint, nil); got != defaultEndpoint {
		t.Errorf("expected default endpoint %q, got %q", defaultEndpoint, got)
	}

	// 环境变量优先于默认值
	t.Setenv(envKey, envEndpoint)
	if got := resolveEndpoint(envKey, defaultEndpoint, nil); got != envEndpoint {
		t.Errorf("expected env endpoint %q, got %q", envEndpoint, got)
	}

	// WithEndpoint option 优先于环境变量
	if got := resolveEndpoint(envKey, defaultEndpoint, []Option{WithEndpoint(optionEndpoint)}); got != optionEndpoint {
		t.Errorf("expected option endpoint %q, got %q", optionEndpoint, got)
	}

	// WithEndpoint 传空串时回退到环境变量
	if got := resolveEndpoint(envKey, defaultEndpoint, []Option{WithEndpoint("")}); got != envEndpoint {
		t.Errorf("expected env endpoint %q, got %q", envEndpoint, got)
	}

	// nil option 不 panic
	if got := resolveEndpoint(envKey, defaultEndpoint, []Option{nil}); got != envEndpoint {
		t.Errorf("expected env endpoint %q, got %q", envEndpoint, got)
	}
}

// TestNewFaceEndpoint 验证 NewFace 的 endpoint 解析
func TestNewFaceEndpoint(t *testing.T) {
	face := NewFace("id", "key")
	if face.Endpoint != defaultFaceidEndpoint {
		t.Errorf("expected default endpoint %q, got %q", defaultFaceidEndpoint, face.Endpoint)
	}
	if face.clientProfile.HttpProfile.Endpoint != defaultFaceidEndpoint {
		t.Errorf("expected clientProfile endpoint %q, got %q", defaultFaceidEndpoint, face.clientProfile.HttpProfile.Endpoint)
	}

	customEndpoint := "faceid.custom.com"
	face = NewFace("id", "key", WithEndpoint(customEndpoint))
	if face.Endpoint != customEndpoint {
		t.Errorf("expected endpoint %q, got %q", customEndpoint, face.Endpoint)
	}
	if face.clientProfile.HttpProfile.Endpoint != customEndpoint {
		t.Errorf("expected clientProfile endpoint %q, got %q", customEndpoint, face.clientProfile.HttpProfile.Endpoint)
	}
}

// TestNewSesEmailSenderEndpoint 验证 NewSesEmailSender 的 endpoint 解析
func TestNewSesEmailSenderEndpoint(t *testing.T) {
	sender := NewSesEmailSender("id", "key", "")
	if sender.Endpoint != defaultSesEndpoint {
		t.Errorf("expected default endpoint %q, got %q", defaultSesEndpoint, sender.Endpoint)
	}
	if sender.Region != defaultEmailRegion {
		t.Errorf("expected default region %q, got %q", defaultEmailRegion, sender.Region)
	}

	customEndpoint := "ses.custom.com"
	sender = NewSesEmailSender("id", "key", "ap-shanghai", WithEndpoint(customEndpoint))
	if sender.Endpoint != customEndpoint {
		t.Errorf("expected endpoint %q, got %q", customEndpoint, sender.Endpoint)
	}
	if sender.Region != "ap-shanghai" {
		t.Errorf("expected region %q, got %q", "ap-shanghai", sender.Region)
	}
}

// TestNewGeneralInvoiceRecognizerEndpoint 验证 NewGeneralInvoiceRecognizer 的 endpoint 解析
func TestNewGeneralInvoiceRecognizerEndpoint(t *testing.T) {
	recognizer := NewGeneralInvoiceRecognizer("id", "key", "")
	if recognizer.Endpoint != defaultOcrEndpoint {
		t.Errorf("expected default endpoint %q, got %q", defaultOcrEndpoint, recognizer.Endpoint)
	}

	customEndpoint := "ocr.custom.com"
	recognizer = NewGeneralInvoiceRecognizer("id", "key", "", WithEndpoint(customEndpoint))
	if recognizer.Endpoint != customEndpoint {
		t.Errorf("expected endpoint %q, got %q", customEndpoint, recognizer.Endpoint)
	}
}

// TestNewEntVerifierEndpoint 验证 NewEntVerifier 的 endpoint 解析
func TestNewEntVerifierEndpoint(t *testing.T) {
	verifier := NewEntVerifier("id", "key")
	if verifier.Endpoint != defaultOcrEndpoint {
		t.Errorf("expected default endpoint %q, got %q", defaultOcrEndpoint, verifier.Endpoint)
	}

	customEndpoint := "ocr.custom.com"
	verifier = NewEntVerifier("id", "key", WithEndpoint(customEndpoint))
	if verifier.Endpoint != customEndpoint {
		t.Errorf("expected endpoint %q, got %q", customEndpoint, verifier.Endpoint)
	}
}
