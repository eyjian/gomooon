// Package mooonpdf
// Wrote by yijian on 2024/10/24
package mooonpdf

import (
	"testing"
)

// go test -v -run="TestJpg2Png"
func TestJpg2Png(t *testing.T) {
	jpgFilepath := "test.jpg"
	pngFilepath, err := Jpg2Png(jpgFilepath)
	if err != nil {
		t.Errorf("jpg2png error: %s", err.Error())
	} else {
		t.Logf("jpg2png success: %s", pngFilepath)
	}
}
