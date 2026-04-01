// Package mooonpdf
// Wrote by yijian on 2024/12/12
package mooonpdf

import (
	"testing"
)

// go test -v -run="TestIsPdfFile"
func TestIsPdfFile(t *testing.T) {
	if IsPdfFile("test.pdf") {
		t.Log("test.pdf is pdf file")
	} else {
		t.Error("test.pdf is not pdf file")
	}
	if IsPdfFile("test.txt") {
		t.Error("test.txt is pdf file")
	} else {
		t.Log("test.txt is not pdf file")
	}
}

// go test -v -run="TestGetPdfPageCount"
func TestGetPdfPageCount(t *testing.T) {
	count, err := GetPdfPageCount("test.pdf")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(count)
	}
}

// go test -v -run="TestValidatePdf"
func TestValidatePdf(t *testing.T) {
	err := ValidatePdf("test.pdf")
	if err != nil {
		t.Error(err)
	} else {
		t.Log("validate ok")
	}
}

// go test -v -run="TestOptimizePdf"
func TestOptimizePdf(t *testing.T) {
	err := OptimizePdf("test.pdf", "new_test.pdf")
	if err != nil {
		t.Error(err)
	} else {
		t.Log("optimize ok")
	}
}