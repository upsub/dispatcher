package util

// Contains checks if the given needle exists within the haystack string
func Contains(haystack []string, needle string) bool {
	for _, key := range haystack {
		if key == needle {
			return true
		}
	}

	return false
}
