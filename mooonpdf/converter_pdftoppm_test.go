// Package mooonpdf
// Wrote by yijian on 2026/05/29
package mooonpdf

import (
	"os"
	"os/exec"
	"testing"

	"github.com/eyjian/gomooon/mooonerror"
)

// go test -v -run="TestPdftoppmConvertAllPages"
func TestPdftoppmConvertAllPages(t *testing.T) {
	converter := NewPdftoppmConverter()

	pdfPath := "../testdata/test.pdf"
	outDir := "../testdata/output"

	// 确保输出目录存在
	os.MkdirAll(outDir, os.ModePerm)
	//defer os.RemoveAll(outDir)

	files, cerr := converter.Convert(pdfPath, outDir, nil, nil)
	if cerr != nil {
		if cerr.ErrCode == mooonerror.ErrCodeToolNotFound {
			t.Skipf("skip test: %s", cerr.ErrMsg)
		}
		t.Fatalf("Convert failed: %s", cerr.ErrMsg)
	}

	if len(files) == 0 {
		t.Fatal("expected at least one image file, but got none")
	}

	for _, f := range files {
		if _, err := os.Stat(f); os.IsNotExist(err) {
			t.Errorf("generated image file does not exist: %s", f)
		}
	}

	t.Logf("generated %d image file(s)", len(files))
}

// go test -v -run="TestPdftoppmConvertSpecificPages"
func TestPdftoppmConvertSpecificPages(t *testing.T) {
	converter := NewPdftoppmConverter()

	pdfPath := "../testdata/test.pdf"
	outDir := "../testdata/output_pages"

	os.MkdirAll(outDir, os.ModePerm)
	defer os.RemoveAll(outDir)

	files, cerr := converter.Convert(pdfPath, outDir, []int{1, 3}, &Pdf2ImageOptions{
		DPI:    200,
		Format: ImageFormatJPG,
	})
	if cerr != nil {
		if cerr.ErrCode == mooonerror.ErrCodeToolNotFound {
			t.Skipf("skip test: %s", cerr.ErrMsg)
		}
		t.Fatalf("Convert failed: %s", cerr.ErrMsg)
	}

	if len(files) != 2 {
		t.Fatalf("expected 2 image files, but got %d", len(files))
	}

	for _, f := range files {
		if _, err := os.Stat(f); os.IsNotExist(err) {
			t.Errorf("generated image file does not exist: %s", f)
		}
	}

	t.Logf("generated %d image file(s)", len(files))
}

func TestPdftoppmConvertFileNotFound(t *testing.T) {
	// 如果 pdftoppm 不可用则跳过此测试
	if _, err := exec.LookPath(pdftoppmCmdName); err != nil {
		t.Skipf("skip test: pdftoppm not available")
	}

	converter := NewPdftoppmConverter()

	_, cerr := converter.Convert("../testdata/notexist.pdf", "../testdata", nil, nil)
	if cerr == nil {
		t.Fatal("expected error, but got nil")
	}
	if cerr.ErrCode != mooonerror.ErrCodeFileNotFound {
		t.Fatalf("expected error code %d, but got %d", mooonerror.ErrCodeFileNotFound, cerr.ErrCode)
	}
}

func TestPdftoppmConvertOutDirNotExist(t *testing.T) {
	// 如果 pdftoppm 不可用则跳过此测试
	if _, err := exec.LookPath(pdftoppmCmdName); err != nil {
		t.Skipf("skip test: pdftoppm not available")
	}

	converter := NewPdftoppmConverter()

	_, cerr := converter.Convert("../testdata/test.pdf", "../testdata/notexist", nil, nil)
	if cerr == nil {
		t.Fatal("expected error, but got nil")
	}
	if cerr.ErrCode != mooonerror.ErrCodeFileNotFound {
		t.Fatalf("expected error code %d, but got %d", mooonerror.ErrCodeFileNotFound, cerr.ErrCode)
	}
}

func TestExtractPageNumber(t *testing.T) {
	tests := []struct {
		filePath       string
		pdfBaseWithExt string
		ext            string
		expected       int
	}{
		{"/tmp/output/test.pdf-1.png", "test.pdf", ".png", 1},
		{"/tmp/output/test.pdf-10.png", "test.pdf", ".png", 10},
		{"/tmp/output/文档.pdf-3.jpg", "文档.pdf", ".jpg", 3},
	}

	for _, tt := range tests {
		result := extractPageNumber(tt.filePath, tt.pdfBaseWithExt, tt.ext)
		if result != tt.expected {
			t.Errorf("extractPageNumber(%s, %s, %s) = %d, expected %d",
				tt.filePath, tt.pdfBaseWithExt, tt.ext, result, tt.expected)
		}
	}
}

// go test -v -run="TestPdftoppmConvertWithCropBox"
// 验证 UseCropBox 选项的效果：
//
//	UseCropBox=true  (默认) → 使用 CropBox，图片尺寸对应裁剪区域，内容充满整张图
//	UseCropBox=false         → 使用 MediaBox，图片尺寸对应整张画布，内容可能只占一小部分
func TestPdftoppmConvertWithCropBox(t *testing.T) {
	converter := NewPdftoppmConverter()

	pdfPath := "../testdata/test.pdf"

	// 1) 默认行为（UseCropBox=nil，等同于 true）
	outDirCrop := "../testdata/output_crop"
	os.MkdirAll(outDirCrop, os.ModePerm)
	defer os.RemoveAll(outDirCrop)

	files, cerr := converter.Convert(pdfPath, outDirCrop, nil, nil)
	if cerr != nil {
		if cerr.ErrCode == mooonerror.ErrCodeToolNotFound {
			t.Skipf("skip test: %s", cerr.ErrMsg)
		}
		t.Fatalf("Convert with cropbox failed: %s", cerr.ErrMsg)
	}
	if len(files) == 0 {
		t.Fatal("expected at least one image file with cropbox, but got none")
	}

	// 2) 显式 UseCropBox=false（使用 MediaBox）
	outDirMedia := "../testdata/output_media"
	os.MkdirAll(outDirMedia, os.ModePerm)
	defer os.RemoveAll(outDirMedia)

	useCropBoxFalse := false
	files2, cerr := converter.Convert(pdfPath, outDirMedia, nil, &Pdf2ImageOptions{
		UseCropBox: &useCropBoxFalse,
	})
	if cerr != nil {
		if cerr.ErrCode == mooonerror.ErrCodeToolNotFound {
			t.Skipf("skip test: %s", cerr.ErrMsg)
		}
		t.Fatalf("Convert with mediabox failed: %s", cerr.ErrMsg)
	}
	if len(files2) == 0 {
		t.Fatal("expected at least one image file with mediabox, but got none")
	}

	// 3) 显式 UseCropBox=true
	outDirCropExplicit := "../testdata/output_crop_explicit"
	os.MkdirAll(outDirCropExplicit, os.ModePerm)
	defer os.RemoveAll(outDirCropExplicit)

	useCropBoxTrue := true
	files3, cerr := converter.Convert(pdfPath, outDirCropExplicit, nil, &Pdf2ImageOptions{
		UseCropBox: &useCropBoxTrue,
	})
	if cerr != nil {
		if cerr.ErrCode == mooonerror.ErrCodeToolNotFound {
			t.Skipf("skip test: %s", cerr.ErrMsg)
		}
		t.Fatalf("Convert with explicit cropbox failed: %s", cerr.ErrMsg)
	}
	if len(files3) == 0 {
		t.Fatal("expected at least one image file with explicit cropbox, but got none")
	}

	t.Logf("cropbox (default): %s", files[0])
	t.Logf("mediabox (false):  %s", files2[0])
	t.Logf("cropbox (true):    %s", files3[0])
}
