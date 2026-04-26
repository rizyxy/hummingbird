package report

import (
	"fmt"
	"hummingbird/internal/models"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"
)

// PrintCLIReport formats and outputs the strategic table and logic summaries to standard output.
func PrintCLIReport(tSum, fSum []models.Summary) {
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 3, ' ', 0)

	printTableSummaries(w, tSum)

	printLogicSummaries(w, fSum)
}

// ExportToMermaid generates Mermaid.js diagram files visualizing the architecture's logic and data flows.
// The generated files are written to the specified target directory.
func ExportToMermaid(targetDirectory string, matches []models.Match, withData bool) error {

	if targetDirectory == "" {
		targetDirectory = "diagrams"
	}

	if err := os.MkdirAll(targetDirectory, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	fmt.Printf("🐦 Hummingbirds: Generating diagrams in %s...\n", targetDirectory)

	var logicSB strings.Builder
	var dataSB strings.Builder

	// Headers
	logicSB.WriteString("graph LR\n    classDef function fill:#bbf,stroke:#333\n")
	if withData {
		dataSB.WriteString("graph LR\n    classDef function fill:#bbf,stroke:#333\n    classDef table fill:#f9f,stroke:#333\n")
	}

	drawnLogic := make(map[string]bool)
	drawnData := make(map[string]bool)

	for _, m := range matches {
		switch m.Category {
		case "function":
			edge := fmt.Sprintf("    %s([%s]) --> %s([%s]):::function", m.FunctionName, m.FunctionName, m.Name, m.Name)
			if !drawnLogic[edge] {
				logicSB.WriteString(edge + "\n")
				drawnLogic[edge] = true
			}
		case "table":
			if withData {
				// Data Graph: Function -> Table
				edge := fmt.Sprintf("    %s([%s]) --> %s[(%s)]:::table", m.FunctionName, m.FunctionName, m.Name, m.Name)
				if !drawnData[edge] {
					dataSB.WriteString(edge + "\n")
					drawnData[edge] = true
				}
			}
		}
	}

	logicFilePath := filepath.Join(targetDirectory, "architecture_logic.mmd")

	if err := os.WriteFile(logicFilePath, []byte(logicSB.String()), 0644); err != nil {
		return err
	}

	if withData {
		dataFilePath := filepath.Join(targetDirectory, "architecture_data.mmd")
		return os.WriteFile(dataFilePath, []byte(dataSB.String()), 0644)
	}

	return nil
}
