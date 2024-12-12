// Package mooonpdf
// Wrote by yijian on 2024/12/12
package mooonpdf

import "testing"

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