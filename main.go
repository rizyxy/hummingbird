package main

import (
	"fmt"
	"hummingbird/internal/models"
	"hummingbird/internal/pkg/analyzer"
	"hummingbird/internal/pkg/cli"
	"hummingbird/internal/pkg/report"
	"hummingbird/internal/pkg/scanner"
	"os"
	"path/filepath"
	"time"
)

func main() {
	cfg := cli.ParseConfig()

	start := time.Now()
	fmt.Println("🚀 Hummingbird: Commencing Audit...")

	// --- 1. Discovery Phase ---
	funcs, err := scanner.ScanFunctions(cfg.Args[1])
	if err != nil {
		fmt.Printf("❌ Error scanning functions: %v\n", err)
		return
	}

	tables, err := scanner.ScanTables(cfg.Args[0])
	if err != nil {
		fmt.Printf("❌ Error loading tables: %v\n", err)
		return
	}

	// --- 2. Scan Phase ---
	var matches []models.Match
	err = filepath.WalkDir(cfg.Args[1], func(path string, d os.DirEntry, err error) error {
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

	// --- 3. Analyze & Report Phase ---
	tSum, fSum := analyzer.GenerateSummaries(funcs, tables, matches)

	if cfg.CLI {
		report.PrintCLIReport(tSum, fSum)
	}

	if cfg.Graph {
		report.ExportToMermaid(cfg.GraphDir, matches)
		fmt.Println("🎨 Graphs generated: architecture_logic.mmd, architecture_data.mmd")
	}

	if cfg.Target != "" {
		analyzer.PrintBlastRadius(cfg.Target, matches)
	}

	fmt.Printf("\n✨ Audit complete in %v\n", time.Since(start))
}
