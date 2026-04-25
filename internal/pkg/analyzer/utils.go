package analyzer

import "hummingbird/internal/models"

func keysToSlice(m map[string]bool) []string {
	var s []string
	for k := range m {
		s = append(s, k)
	}
	return s
}

func buildSlice(counts map[string]int, tracker map[string]map[string]bool) []models.Summary {
	var s []models.Summary
	for name, count := range counts {
		unique := len(tracker[name])
		s = append(s, models.Summary{Name: name, TotalRefs: count, UniqueFuncs: unique, FrictionScore: count * unique})
	}
	return s
}
