package connection

import (
	"fmt"
	"strconv"
)

func interfaceInt64(val interface{}) int64 {
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
			fmt.Printf("Error handling chip request value: %s\n", err)
		}
		value = int64(v)
		break
	default:
		break
	}

	return value
}
