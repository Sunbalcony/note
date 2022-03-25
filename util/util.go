package util

import (
	"bytes"
	"math/rand"
	"time"
)

const char = "abcdefghijklmnopqrstuvwxyz0123456789"

func RandChar(size int) string {
	rand.Seed(time.Now().UnixNano())
	//rand.NewSource(time.Now().UnixNano())
	var s bytes.Buffer
	for i := 0; i < size; i++ {
		s.WriteByte(char[rand.Int63()%int64(len(char))])
	}
	return s.String()
}
