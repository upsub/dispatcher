package util

// Merge two maps of strings together
func Merge(base map[string]string, new map[string]string) map[string]string {
	for key, value := range new {
		base[key] = value
	}

	return base
}
