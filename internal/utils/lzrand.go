package utils

import (
	"math/rand"
	"sync"
	"time"
)

const (
	HexMeta = "0123456789abcdef"

	AsciiLitter = "qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM"
	UpperLitter = "QWERTYUIOPASDFGHJKLZXCVBNM"
	LowerLitter = "qwertyuiopasdfghjklzxcvbnm"
	Digits      = "1234567890"
)

var (
	rnd      = rand.New(rand.NewSource(time.Now().UnixNano()))
	rndMutex sync.Mutex
)

// 随机生成 MD5 HASH 值
func RandomMD5Hash() string {
	return RandStrWithMeta(32, HexMeta)
}

// 随机生成指定长度的字符串
func RandomHexString(size int) (ret string) {
	return RandStrWithMeta(size, HexMeta)
}

// 随机int, 但max不會隨機到，需要+1
func RandInt(min, max int) int {
	rndMutex.Lock()
	defer rndMutex.Unlock()
	if min >= max || max == 0 {
		return max
	}

	return rnd.Intn(max-min) + min
}

// 自定义meta的随机
func RandStrWithMeta(n int, metaData string) string {
	bytes := []byte(metaData)
	result := make([]byte, 0)
	//r :=
	for i := 0; i < n; i++ {
		result = append(result, bytes[RandInt(0, len(bytes))])
	}
	return string(result)
}

// 随机数字字符串
func RandDigital(n int) string {
	if n == 0 {
		n = RandInt(3, 20)
	}
	return RandStrWithMeta(n, Digits)
}

// 随机大小写英文字母字符串
func RandStr(n int) string {
	if n == 0 {
		n = RandInt(3, 20)
	}
	return RandStrWithMeta(n, AsciiLitter)
}

func RandStrAndDigital(n int) string {
	if n == 0 {
		n = RandInt(3, 20)
	}
	return RandStrWithMeta(n, AsciiLitter+Digits)
}

func RandBytes(n int) []byte {
	if n == 0 {
		n = RandInt(3, 20)
	}
	dst := make([]byte, 0)
	for i := 0; i < n; i++ {
		dst = append(dst, byte(RandInt(0, 255)))
	}
	return dst
}
