// Package mooonpdf
// Wrote by yijian on 2025/06/14
package mooonpdf

import "testing"

// go test -v -run="TestMergeFiles"
func TestMergeFiles(t *testing.T) {
	inFiles := []string{"out_dir/test_1.pdf", "out_dir/test_3.pdf"}
	outFile := "out_dir/out.pdf"

	err := MergeFiles(inFiles, outFile)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(outFile)
}
