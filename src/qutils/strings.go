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

// IfacetoI64 converts numeric interface to int64
func IfacetoI64(i interface{}) int64 {
	switch i.(type) {
	case int:
		return int64(i.(int))
	case int64:
		return i.(int64)
	case float64:
		return int64(i.(float64))
	default:
		// panic(fmt.Sprintf("Cannot convert %s to int64", i))
	}

	return int64(0)
}

// IfacetoInt converts numeric interface to int
func IfacetoInt(i interface{}) int {
	switch i.(type) {
	case int:
		return i.(int)
	case int64:
		return int(i.(int64))
	case float64:
		return int(i.(float64))
	default:
		// panic(fmt.Sprintf("Cannot convert %s to int64", i))
	}

	return 0
}
