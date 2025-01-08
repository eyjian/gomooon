// Package mooonpdf
// Wrote by yijian on 2024/12/12
package mooonpdf

import (
	"github.com/gen2brain/go-fitz"
	"strings"
)

// IsPdfFile 判断文件是否为 pdf 文件
func IsPdfFile(filepath string) bool {
	if filepath == "" {
		return false
	}
	if len(filepath) <= 4 {
		return false
	}
	if strings.ToLower(filepath[len(filepath)-4:]) != ".pdf" {
		return false
	}

	doc, err := fitz.New(filepath)
	if err != nil {
		return false
	} else {
		doc.Close()
		return true
	}
}

// GetPdfPageCount 获取 pdf 文件页数
func GetPdfPageCount(filepath string) (int, error) {
	doc, err := fitz.New(filepath)
	if err != nil {
		return 0, err
	} else {
		pageCount := doc.NumPage()
		doc.Close()
		return pageCount, nil
	}
}