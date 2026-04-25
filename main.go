package main

import (
	"fmt"
	"hummingbird/internal/models"
	"hummingbird/internal/pkg/analyzer"
	"hummingbird/internal/pkg/cli"
	"hummingbird/internal/pkg/db"
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
	var tablesFile, codebasePath string
	if len(cfg.Args) == 1 {
		codebasePath = cfg.Args[0]
	} else {
		tablesFile = cfg.Args[0]
		codebasePath = cfg.Args[1]
	}

	funcs, err := scanner.ScanFunctions(codebasePath)
	if err != nil {
		fmt.Printf("❌ Error scanning functions: %v\n", err)
		return
	}

	var tables []string
	if tablesFile != "" {
		fileTables, err := scanner.ScanTables(tablesFile)
		if err != nil {
			fmt.Printf("❌ Error loading tables: %v\n", err)
			return
		}
		tables = append(tables, fileTables...)
	}

	if cfg.DBDriver != "" && cfg.DBDsn != "" {
		fmt.Printf("🔍 Fetching tables from %s database...\n", cfg.DBDriver)
		dbTables, err := db.FetchTables(cfg.DBDriver, cfg.DBDsn)
		if err != nil {
			fmt.Printf("❌ Error fetching tables from db: %v\n", err)
			return
		}
		tables = append(tables, dbTables...)
		fmt.Printf("✅ Found %d tables in database\n", len(dbTables))
	}

	// --- 2. Scan Phase ---
	var matches []models.Match
	err = filepath.WalkDir(codebasePath, func(path string, d os.DirEntry, err error) error {
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
		withData := len(tables) > 0
		report.ExportToMermaid(cfg.GraphDir, matches, withData)
		if withData {
			fmt.Println("🎨 Graphs generated: architecture_logic.mmd, architecture_data.mmd")
		} else {
			fmt.Println("🎨 Graphs generated: architecture_logic.mmd")
		}
	}

	if cfg.Target != "" {
		analyzer.PrintBlastRadius(cfg.Target, matches)
	}

	fmt.Printf("\n✨ Audit complete in %v\n", time.Since(start))
}
