package analyzer

import "hummingbird/internal/models"

// keysToSlice converts a map of boolean flags into a slice of strings (the keys).
// This is used to extract the names of referenced functions or tables from the tracker maps.
func keysToSlice(m map[string]bool) []string {
	var s []string
	for k := range m {
		s = append(s, k)
	}
	return s
}

// buildSlice constructs a slice of Summary structs from frequency counts and call tracking maps.
// It calculates a FrictionScore based on the number of references and unique calling functions.
func buildSlice(counts map[string]int, tracker map[string]map[string]bool) []models.Summary {
	var s []models.Summary
	for name, count := range counts {
		unique := len(tracker[name])
		s = append(s, models.Summary{Name: name, TotalRefs: count, UniqueFuncs: unique, FrictionScore: count * unique})
	}
	return s
}
