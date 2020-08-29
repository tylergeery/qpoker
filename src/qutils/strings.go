package qutils

import (
	"math/rand"
)

const slugBytes = "abcdefghijklmnopqrstuvwxyz123456789"

// GenerateSlug create a slug of length n
func GenerateSlug(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = slugBytes[rand.Int63()%int64(len(slugBytes))]
	}

	return string(b)
}

// GenerateVariedLengthSlug create a varied length slug between min and max
func GenerateVariedLengthSlug(min, max int) string {
	extra := rand.Intn(max - min)

	return GenerateSlug(min + extra)
}
