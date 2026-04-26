package analyzer

import (
	"fmt"
	"hummingbird/internal/models"
	"sort"
)

// PrintBlastRadius calculates and prints the blast radius for a target table, showing direct and indirect impacts.
func PrintBlastRadius(target string, matches []models.Match) {
	radius := CalculateBlastRadius(target, matches)
	fmt.Printf("\n☢️  BLAST RADIUS: %s\n", target)
	fmt.Printf("   Directly Impacted:   %d functions\n", len(radius.DirectImpact))
	fmt.Printf("   Indirectly Impacted: %d functions\n", len(radius.IndirectImpact))
	fmt.Printf("   Total Risk Score:    %d\n", radius.TotalRiskScore)
}

// CalculateBlastRadius recursively determines all functions that directly or indirectly reference a target table.
func CalculateBlastRadius(tableName string, matches []models.Match) models.ImpactReport {
	direct := make(map[string]bool)
	indirect := make(map[string]bool)

	for _, m := range matches {
		if m.Category == "table" && m.Name == tableName {
			direct[m.FunctionName] = true
		}
	}

	for {
		addedNew := false
		for _, m := range matches {
			if m.Category == "function" {
				if direct[m.Name] || indirect[m.Name] {
					if !indirect[m.FunctionName] && !direct[m.FunctionName] {
						indirect[m.FunctionName] = true
						addedNew = true
					}
				}
			}
		}
		if !addedNew {
			break
		}
	}

	return models.ImpactReport{
		TargetTable:    tableName,
		DirectImpact:   keysToSlice(direct),
		IndirectImpact: keysToSlice(indirect),
		TotalRiskScore: len(direct) + (len(indirect) * 2),
	}
}

// GenerateSummaries aggregates match data to produce prioritized summaries of table usage and function dependencies.
func GenerateSummaries(allFuncs []string, allTables []string, matches []models.Match) ([]models.Summary, []models.Summary) {
	tableMap := make(map[string]map[string]bool)
	funcMap := make(map[string]map[string]bool)
	tableCounts := make(map[string]int)
	funcCounts := make(map[string]int)

	// 1. Initialize all known entities to 0 to ensure they show up in the report
	for _, t := range allTables {
		tableCounts[t] = 0
		tableMap[t] = make(map[string]bool)
	}
	for _, f := range allFuncs {
		funcCounts[f] = 0
		funcMap[f] = make(map[string]bool)
	}

	// 2. Process matches
	for _, m := range matches {
		if m.Category == "table" {
			tableCounts[m.Name]++
			if tableMap[m.Name] == nil {
				tableMap[m.Name] = make(map[string]bool)
			}
			tableMap[m.Name][m.FunctionName] = true
		} else {
			// Exclude self-references (the function calling itself or the definition)
			if m.Name == m.FunctionName {
				continue
			}

			funcCounts[m.Name]++
			if funcMap[m.Name] == nil {
				funcMap[m.Name] = make(map[string]bool)
			}
			funcMap[m.Name][m.FunctionName] = true
		}
	}

	tSum := buildSlice(tableCounts, tableMap)
	fSum := buildSlice(funcCounts, funcMap)

	// 3. Sort by Friction Score (Highest first)
	sort.Slice(tSum, func(i, j int) bool {
		return tSum[i].FrictionScore > tSum[j].FrictionScore
	})

	sort.Slice(fSum, func(i, j int) bool {
		return fSum[i].FrictionScore > fSum[j].FrictionScore
	})

	return tSum, fSum
}
