package cli

import (
	"flag"
	"fmt"
	"hummingbird/internal/models"
	"os"
)

func PrintCustomUsage() {
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

func ParseConfig() *models.Config {
	c := &models.Config{}

	flag.BoolVar(&c.CLI, "cli", false, "Print prioritized Strategic Summary and Logic Call tables to terminal")
	flag.BoolVar(&c.Graph, "graph", false, "Generate separated Mermaid JS files")
	flag.StringVar(&c.Target, "target", "", "Calculate the recursive 'Blast Radius' for a table")
	flag.StringVar(&c.GraphDir, "out", "diagrams", "Directory to save generated Mermaid files")

	flag.Usage = PrintCustomUsage
	flag.Parse()

	c.Args = flag.Args()
	if len(c.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	return c
}
