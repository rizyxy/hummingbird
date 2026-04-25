package report

import (
	"fmt"
	"hummingbird/internal/models"
	"text/tabwriter"
)

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
