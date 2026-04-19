package report

import (
	"fmt"
	"hummingbird/internal/models"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
)

func GenerateSummaries(matches []models.Match) ([]models.Summary, []models.Summary) {
	tableMap := make(map[string]map[string]bool)
	funcMap := make(map[string]map[string]bool)
	tableCounts := make(map[string]int)
	funcCounts := make(map[string]int)

	for _, m := range matches {
		if m.Category == "table" {
			tableCounts[m.Name]++
			if tableMap[m.Name] == nil {
				tableMap[m.Name] = make(map[string]bool)
			}
			tableMap[m.Name][m.FunctionName] = true
		} else {
			funcCounts[m.Name]++
			if funcMap[m.Name] == nil {
				funcMap[m.Name] = make(map[string]bool)
			}
			funcMap[m.Name][m.FunctionName] = true
		}
	}

	tSum := buildSlice(tableCounts, tableMap)
	fSum := buildSlice(funcCounts, funcMap)

	// --- SORTING LOGIC ---
	// Sort Tables by Friction Score (Highest first)
	sort.Slice(tSum, func(i, j int) bool {
		return tSum[i].FrictionScore > tSum[j].FrictionScore
	})

	// Sort Functions by Friction Score (Highest first)
	sort.Slice(fSum, func(i, j int) bool {
		return fSum[i].FrictionScore > fSum[j].FrictionScore
	})

	return tSum, fSum
}

func buildSlice(counts map[string]int, tracker map[string]map[string]bool) []models.Summary {
	var s []models.Summary
	for name, count := range counts {
		unique := len(tracker[name])
		s = append(s, models.Summary{Name: name, TotalRefs: count, UniqueFuncs: unique, FrictionScore: count * unique})
	}
	return s
}

func PrintCLIReport(tSum, fSum []models.Summary) {
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 3, ' ', 0)

	// --- TABLE REPORT ---
	fmt.Println("\n📊 STRATEGIC TABLE SUMMARY")
	fmt.Fprintln(w, "TABLE NAME\tREFS\tFUNCS\tFRICTION\tRISK")
	fmt.Fprintln(w, "----------\t----\t-----\t--------\t------")
	for _, s := range tSum {
		status := getRiskStatus(s.FrictionScore)
		fmt.Fprintf(w, "%s\t%d\t%d\t%d\t%s\n", s.Name, s.TotalRefs, s.UniqueFuncs, s.FrictionScore, status)
	}
	w.Flush()

	// --- LOGIC REPORT ---
	fmt.Println("\n🧠 LOGIC CALL SUMMARY")
	fmt.Fprintln(w, "FUNCTION NAME\tREFS\tCALLERS\tFRICTION\tRISK")
	fmt.Fprintln(w, "-------------\t----\t-------\t--------\t------")
	for _, s := range fSum {
		// Only show functions that are reused (refs > 1) to keep the report focused
		if s.TotalRefs > 1 {
			status := getRiskStatus(s.FrictionScore)
			fmt.Fprintf(w, "%s\t%d\t%d\t%d\t%s\n", s.Name, s.TotalRefs, s.UniqueFuncs, s.FrictionScore, status)
		}
	}
	w.Flush()
}

// Helper to centralize risk logic
func getRiskStatus(score int) string {
	if score > 50 {
		return "🔥 CRITICAL"
	} else if score > 20 {
		return "⚠️  HIGH RISK"
	} else if score > 5 {
		return "⚖️  MEDIUM"
	}
	return "✅ LOW RISK"
}

// ExportToMermaid generates a styled Mermaid JS diagram with layered subgraphs
func ExportToMermaid(matches []models.Match) error {
	fmt.Println("🐦 Hummingbirds: Generating separated Mermaid diagrams...")

	var logicSB strings.Builder
	var dataSB strings.Builder

	// Headers
	logicSB.WriteString("graph LR\n    classDef function fill:#bbf,stroke:#333\n")
	dataSB.WriteString("graph LR\n    classDef function fill:#bbf,stroke:#333\n    classDef table fill:#f9f,stroke:#333\n")

	drawnLogic := make(map[string]bool)
	drawnData := make(map[string]bool)

	for _, m := range matches {
		if m.Category == "function" {
			// Logic Graph: Function -> Function
			edge := fmt.Sprintf("    %s([%s]) --> %s([%s]):::function", m.FunctionName, m.FunctionName, m.Name, m.Name)
			if !drawnLogic[edge] {
				logicSB.WriteString(edge + "\n")
				drawnLogic[edge] = true
			}
		} else if m.Category == "table" {
			// Data Graph: Function -> Table
			edge := fmt.Sprintf("    %s([%s]) --> %s[(%s)]:::table", m.FunctionName, m.FunctionName, m.Name, m.Name)
			if !drawnData[edge] {
				dataSB.WriteString(edge + "\n")
				drawnData[edge] = true
			}
		}
	}

	// Save Logic Graph
	if err := os.WriteFile("architecture_logic.mmd", []byte(logicSB.String()), 0644); err != nil {
		return err
	}

	// Save Data Graph
	return os.WriteFile("architecture_data.mmd", []byte(dataSB.String()), 0644)
}

func CalculateBlastRadius(tableName string, matches []models.Match) models.ImpactReport {
	direct := make(map[string]bool)
	indirect := make(map[string]bool)

	// 1. Find Direct Impact (Who touches the table?)
	for _, m := range matches {
		if m.Category == "table" && m.Name == tableName {
			direct[m.FunctionName] = true
		}
	}

	// 2. Find Indirect Impact (Who calls the callers?)
	// We iterate to find the next layer of the chain
	for {
		addedNew := false
		for _, m := range matches {
			if m.Category == "function" {
				// If this function calls a function already in our impact list
				if direct[m.Name] || indirect[m.Name] {
					// And it's not already tracked as an indirect impact
					if !indirect[m.FunctionName] && !direct[m.FunctionName] {
						indirect[m.FunctionName] = true
						addedNew = true
					}
				}
			}
		}
		if !addedNew {
			break
		} // Stop when no more callers are found
	}

	return models.ImpactReport{
		TargetTable:    tableName,
		DirectImpact:   keysToSlice(direct),
		IndirectImpact: keysToSlice(indirect),
		TotalRiskScore: len(direct) + (len(indirect) * 2), // Indirect hits are often riskier/harder to find
	}
}

func keysToSlice(m map[string]bool) []string {
	var s []string
	for k := range m {
		s = append(s, k)
	}
	return s
}
