package qutils

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
