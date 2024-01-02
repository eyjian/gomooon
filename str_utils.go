// Package gomooon
// Wrote by yijian on 2024/01/02
package gomooon

import (
	"math/rand"
	"sync"
	"time"
)

const allCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const hexCharset = "abcdefABCDEF0123456789"

var (
	r      *rand.Rand
	mu     sync.Mutex
	once   sync.Once
	buffer sync.Pool
)

func init() {
	source := rand.NewSource(time.Now().UnixNano())
	r = rand.New(source)
	buffer.New = func() interface{} {
		return make([]byte, 0, 64)
	}
}

func getNonceStr(length int, charset string) string {
	once.Do(func() {
		mu.Lock()
		defer mu.Unlock()
		r.Seed(time.Now().UnixNano())
	})

	mu.Lock()
	defer mu.Unlock()

	// Get a buffer from the pool and reset its length to the desired value
	buf := buffer.Get().([]byte)[:length]

	for i := range buf {
		buf[i] = charset[r.Intn(len(charset))]
	}

	// Convert the buffer to a string, put it back into the pool, and return the result
	result := string(buf)
	buffer.Put(buf)
	return result
}

func GetNonceStr(length int) string {
	return getNonceStr(length, allCharset)
}

func GetHexNonceStr(length int) string {
	return getNonceStr(length, hexCharset)
}
