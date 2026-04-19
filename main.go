package main

import (
	"flag"
	"fmt"
	"hummingbird/internal/models"
	"hummingbird/pkg/parser"
	"hummingbird/pkg/report"
	"hummingbird/pkg/scanner"
	"os"
	"path/filepath"
	"time"
)

func main() {
	// 1. Flag Definitions
	cli := flag.Bool("cli", false, "Print prioritized Strategic Summary and Logic Call tables to terminal")
	graph := flag.Bool("graph", false, "Generate separated Mermaid JS files (architecture_logic.mmd, architecture_data.mmd)")
	target := flag.String("target", "", "Calculate the recursive 'Blast Radius' (direct/indirect impact) for a table")

	// 2. Comprehensive Help Override
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\n🐦 HUMMINGBIRD | Migration Intelligence Engine\n")
		fmt.Fprintf(os.Stderr, "--------------------------------------------------\n")
		fmt.Fprintf(os.Stderr, "Hummingbird maps database dependencies and logic paths to identify\n")
		fmt.Fprintf(os.Stderr, "high-friction areas and mitigate risks during system migrations.\n\n")

		fmt.Fprintf(os.Stderr, "USAGE:\n")
		fmt.Fprintf(os.Stderr, "  hummingbird [flags] <tables_file> <codebase_path>\n\n")

		fmt.Fprintf(os.Stderr, "ARGUMENTS:\n")
		fmt.Fprintf(os.Stderr, "  <tables_file>    Path to a .txt file containing target table names (one per line)\n")
		fmt.Fprintf(os.Stderr, "  <codebase_path>  Directory containing the source code to audit (.go, .js, .ts, etc.)\n\n")

		fmt.Fprintf(os.Stderr, "FLAGS:\n")
		flag.PrintDefaults()

		fmt.Fprintf(os.Stderr, "\nEXAMPLES:\n")
		fmt.Fprintf(os.Stderr, "  # Standard audit with CLI output\n")
		fmt.Fprintf(os.Stderr, "  hummingbird --cli tables.txt ./src\n\n")

		fmt.Fprintf(os.Stderr, "  # Generate visualization graphs for architect review\n")
		fmt.Fprintf(os.Stderr, "  hummingbird --graph tables.txt ./src\n\n")

		fmt.Fprintf(os.Stderr, "  # Calculate risk impact for a specific critical table\n")
		fmt.Fprintf(os.Stderr, "  hummingbird --target TBL_USER_AUTH tables.txt ./src\n")
		fmt.Fprintf(os.Stderr, "--------------------------------------------------\n\n")
	}

	flag.Parse()
	args := flag.Args()

	// 3. Validation Logic
	if len(args) < 2 {
		flag.Usage()
		return
	}

	start := time.Now()
	fmt.Println("🚀 Hummingbird: Commencing Audit...")

	// --- 1. Discovery Phase ---
	funcs, err := scanner.SurveyFunctions(args[1])
	if err != nil {
		fmt.Printf("❌ Error surveying functions: %v\n", err)
		return
	}

	tables, err := parser.LoadTables(args[0])
	if err != nil {
		fmt.Printf("❌ Error loading tables: %v\n", err)
		return
	}

	// --- 2. Scan Phase ---
	var matches []models.Match
	err = filepath.WalkDir(args[1], func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		matches = append(matches, scanner.ScanFileContent(path, tables, funcs)...)
		return nil
	})

	if err != nil {
		fmt.Printf("❌ Error during scan: %v\n", err)
		return
	}

	// --- 3. Reporting Phase ---
	tSum, fSum := report.GenerateSummaries(matches)

	if *cli {
		report.PrintCLIReport(tSum, fSum)
	}

	if *graph {
		// Ensure your report package has this function name updated
		report.ExportToMermaid(matches)
		fmt.Println("🎨 Graphs generated: architecture_logic.mmd, architecture_data.mmd")
	}

	if *target != "" {
		radius := report.CalculateBlastRadius(*target, matches)
		fmt.Printf("\n☢️  BLAST RADIUS: %s\n", *target)
		fmt.Printf("   Directly Impacted:   %d functions\n", len(radius.DirectImpact))
		fmt.Printf("   Indirectly Impacted: %d functions (callers of callers)\n", len(radius.IndirectImpact))
		fmt.Printf("   Total Risk Score:    %d\n", radius.TotalRiskScore)
	}

	fmt.Printf("\n✨ Audit complete in %v\n", time.Since(start))
}
