// Package mooonutils
// Wrote by yijian on 2025/06/14
package mooonutils

import "testing"

// go test -v -run="TestDirExists"
func TestDirExists(t *testing.T) {
	dir := "./tmp"
	exists, err := DirExists(dir)
	if err != nil {
		t.Log(err)
	} else {
		t.Errorf("dir exists: %v", exists)
	}

	dir = "."
	exists, err = DirExists(dir)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("dir `%s` exists: %v", dir, exists)
}
