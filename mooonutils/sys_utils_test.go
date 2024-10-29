// Package mooonutils
// Wrote by yijian on 2024/10/25
package mooonutils

import (
	"testing"
)

// go test -v -run="TestGetProgramDir"
func TestGetProgramDir(t *testing.T) {
	dir, err := GetProgramDir()
	if err != nil {
		t.Errorf("GetProgramDir error: %s\n", err.Error())
	} else {
		t.Logf("GetProgramDir: %s\n", dir)
	}
}
