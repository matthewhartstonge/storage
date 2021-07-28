package utils

// AppendToStringSet adds the provided items to the slice, if they don't already
// exist.
func AppendToStringSet(strings []string, items ...string) []string {
	for i := range items {
		found := false
		for j := range strings {
			if items[i] == strings[j] {
				found = true

				break
			}
		}

		if !found {
			strings = append(strings, items[i])
		}
	}

	return strings
}

// RemoveFromStringSet removes all matching items from the provided slice.
func RemoveFromStringSet(strings []string, items ...string) []string {
	for i := range items {
		for j := len(strings) - 1; j >= 0; j-- {
			if items[i] == strings[j] {
				copy(strings[j:], strings[j+1:])
				strings[len(strings)-1] = ""
				strings = strings[:len(strings)-1]
			}
		}
	}

	return strings
}
