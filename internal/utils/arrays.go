package utils

// AddUnique is a helper function to add unique items to a slice
func AddUnique(slice []string, items ...string) []string {
	existing := make(map[string]bool)
	for _, s := range slice {
		existing[s] = true
	}

	for _, item := range items {
		if !existing[item] {
			slice = append(slice, item)
			existing[item] = true
		}
	}
	return slice
}

// RemoveItems is a helper function to remove items from a slice
func RemoveItems(slice []string, items ...string) []string {
	removeMap := make(map[string]bool)
	for _, item := range items {
		removeMap[item] = true
	}

	result := make([]string, 0, len(slice))
	for _, s := range slice {
		if !removeMap[s] {
			result = append(result, s)
		}
	}
	return result
}
