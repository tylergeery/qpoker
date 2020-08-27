package qutils

// ToI64 converts numeric interface to int64
func ToI64(i interface{}) int64 {
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

// ToInt converts numeric interface to int
func ToInt(i interface{}) int {
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

// MaxInt64 gets max of int64 values
func MaxInt64(nums ...int64) int64 {
	if len(nums) == 0 {
		return 0
	}

	max := nums[0]
	for _, num := range nums {
		if num > max {
			max = num
		}
	}

	return max
}

// MaxInt gets max of int values
func MaxInt(nums ...int) int {
	if len(nums) == 0 {
		return 0
	}

	max := nums[0]
	for _, num := range nums {
		if num > max {
			max = num
		}
	}

	return max
}

// MinInt64 gets max of int64 values
func MinInt64(nums ...int64) int64 {
	if len(nums) == 0 {
		return 0
	}

	min := nums[0]
	for _, num := range nums {
		if num < min {
			min = num
		}
	}

	return min
}

// MinInt gets min of int values
func MinInt(nums ...int) int {
	if len(nums) == 0 {
		return 0
	}

	min := nums[0]
	for _, num := range nums {
		if num < min {
			min = num
		}
	}

	return min
}
