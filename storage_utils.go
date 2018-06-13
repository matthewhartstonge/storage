package storage

// stringArrayEquals returns a bool based on the equality of two arrays
func stringArrayEquals(a []string, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i, v := range a {
		if b[i] != v {
			return false
		}
	}
	return true
}
