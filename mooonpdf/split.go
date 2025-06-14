// Package mooonpdf
// Wrote by yijian on 2025/06/14
package mooonpdf

import (
	pdfapi "github.com/pdfcpu/pdfcpu/pkg/api"
	pdfmodel "github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

// SplitFile 将 pdf 文件的每页都拆分为一个 pdf 文件
// span 拆分跨度，1 表示每页拆分为一个 pdf 文件，2 表示每 2 页拆分为一个 pdf 文件，依次类推。但 0 表示书签进行拆分。
func SplitFile(inFile, outDir string, span int) error {
	config := pdfmodel.NewDefaultConfiguration()
	return pdfapi.SplitFile(inFile, outDir, span, config)
}
