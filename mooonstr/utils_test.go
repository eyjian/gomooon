// Package mooonstr
// Wrote by yijian on 2024/09/03
package mooonstr

import "testing"

// go test -v -run="TestJoinInt32$"
func TestJoinInt32(t *testing.T) {
	elems := []int32{-1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	res := JoinInt32(elems, ",")
	t.Logf("result: %s\n", res)
}

// go test -v -run="TestJoinUint32$"
func TestJoinUint32(t *testing.T) {
	elems := []uint32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	res := JoinUint32(elems, ",")
	t.Logf("result: %s\n", res)
}

// go test -v -run="TestJoinInt64$"
func TestJoinInt64(t *testing.T) {
	elems := []int64{-10, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	res := JoinInt64(elems, ",")
	t.Logf("result: %s\n", res)
}

// go test -v -run="TestJoinUint64$"
func TestJoinUint64(t *testing.T) {
	elems := []uint64{10, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	res := JoinUint64(elems, ",")
	t.Logf("result: %s\n", res)
}