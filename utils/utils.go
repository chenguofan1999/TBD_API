package utils

import (
	"fmt"
	"math/rand"
	"path"
	"time"
)

var r *rand.Rand

func init() {
	r = rand.New(rand.NewSource(time.Now().Unix()))
}

func GenerateRandomFileName(fileName string) string {
	fileSuffix := path.Ext(fileName)
	randPrefix := randString(20)
	return fmt.Sprintf("%s%s", randPrefix, fileSuffix)
}

// RandString 生成随机字符串
func randString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		b := r.Intn(26) + 97
		bytes[i] = byte(b)
	}
	return string(bytes)
}
