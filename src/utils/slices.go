package utils

// IntSliceHasValue returns whether an int is in Slice
func IntSliceHasValue(haystack []int, needle int) bool {
	for _, val := range haystack {
		if needle == val {
			return true
		}
	}

	return false
}

// StringSliceHasValue returns whether a string is in slice
func StringSliceHasValue(haystack []string, needle string) bool {
	for _, val := range haystack {
		if needle == val {
			return true
		}
	}

	return false
}
