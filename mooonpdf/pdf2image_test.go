// Package mooonpdf
// Wrote by yijian on 2024/10/24
package mooonpdf

import (
	"testing"
)

// go test -v -run="TestPdf2Image"
func TestPdf2Image(t *testing.T) {
	pdfFilepath := "test.pdf"

	imagePaths, err := Pdf2Png(pdfFilepath)
	if err != nil {
		t.Errorf("pdf2png error: %s", err.Error())
	} else {
		t.Logf("pdf2png success: %+v", imagePaths)
	}

	imagePaths, err = Pdf2Jpg(pdfFilepath)
	if err != nil {
		t.Errorf("Pdf2Jpg error: %s", err.Error())
	} else {
		t.Logf("Pdf2Jpg success: %+v", imagePaths)
	}
}
