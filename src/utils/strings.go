package utils

import (
	"math/rand"
)

const slugBytes = "abcdefghijklmnopqrstuvwxyz123456789"

func GenerateSlug(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = slugBytes[rand.Int63()%int64(len(slugBytes))]
	}

	return string(b)
}
