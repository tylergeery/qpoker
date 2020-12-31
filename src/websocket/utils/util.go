package utils

import (
	"fmt"
	"strconv"
)

// InterfaceInt64 returns interface value as int64
func InterfaceInt64(val interface{}) int64 {
	value := int64(0)

	switch val.(type) {
	case float64:
		value = int64(val.(float64))
		break
	case int64:
		value = val.(int64)
		break
	case string:
		v, err := strconv.Atoi(val.(string))
		if err != nil {
			fmt.Printf("Error handling interface to int64: %s\n", err)
		}
		value = int64(v)
		break
	default:
		break
	}

	return value
}
