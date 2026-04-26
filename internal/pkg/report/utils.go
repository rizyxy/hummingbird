package report

import (
	"fmt"
	"hummingbird/internal/models"
	"text/tabwriter"
)

// printTableSummaries renders a prioritized summary of database table usage, including reference counts and calculated risk metrics.
func printTableSummaries(w *tabwriter.Writer, tSum []models.Summary) {
	fmt.Println("\n📊 STRATEGIC TABLE SUMMARY")
	fmt.Fprintln(w, "TABLE NAME\tREFS\tFUNCS\tFRICTION\tRISK")
	fmt.Fprintln(w, "----------\t----\t-----\t--------\t------")
	for _, s := range tSum {
		status := getRiskStatus(s.FrictionScore)
		fmt.Fprintf(w, "%s\t%d\t%d\t%d\t%s\n", s.Name, s.TotalRefs, s.UniqueFuncs, s.FrictionScore, status)
	}
	w.Flush()
}

// printLogicSummaries renders a prioritized summary of logic dependencies, highlighting highly-coupled functions and their risk scores.
func printLogicSummaries(w *tabwriter.Writer, fSum []models.Summary) {
	fmt.Println("\n🧠 LOGIC CALL SUMMARY")
	fmt.Fprintln(w, "FUNCTION NAME\tREFS\tCALLERS\tFRICTION\tRISK")
	fmt.Fprintln(w, "-------------\t----\t-------\t--------\t------")
	for _, s := range fSum {
		if s.TotalRefs >= 0 {
			status := getRiskStatus(s.FrictionScore)
			fmt.Fprintf(w, "%s\t%d\t%d\t%d\t%s\n", s.Name, s.TotalRefs, s.UniqueFuncs, s.FrictionScore, status)
		}
	}
	w.Flush()
}

// getRiskStatus evaluates a numerical friction score and returns a formatted severity label.
// Higher scores signify tighter coupling and greater potential impact during migration.
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
