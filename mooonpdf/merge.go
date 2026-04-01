// Package mooonpdf
// Wrote by yijian on 2025/06/14
package mooonpdf

import (
	pdfapi "github.com/pdfcpu/pdfcpu/pkg/api"
	pdfmodel "github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

// MergeFiles 将多个 pdf 文件合并为一个 pdf 文件
func MergeFiles(inFiles []string, outFile string) error {
	return pdfapi.MergeCreateFile(inFiles, outFile, false, pdfmodel.NewDefaultConfiguration())
}
