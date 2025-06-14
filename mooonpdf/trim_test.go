// Package mooonpdf
// Wrote by yijian on 2025/06/14
package mooonpdf

import "testing"

// go test -v -run="TestTrimFile"
func TestTrimFile(t *testing.T) {
	inFile := "test.pdf"

	err := TrimFile(inFile, "out_dir/odd.pdf", []string{"odd"})
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	err = TrimFile(inFile, "out_dir/1-2.pdf", []string{"1-2"})
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	err = TrimFile(inFile, "out_dir/1_3.pdf", []string{"1", "3"})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	t.Log("success")
}
