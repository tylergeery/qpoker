package qutils

// HasTrueValues checks if the map has any true values
func HasTrueValues(m map[string]bool) bool {
	if m == nil {
		return false
	}

	for _, val := range m {
		if val {
			return true
		}
	}

	return false
}
