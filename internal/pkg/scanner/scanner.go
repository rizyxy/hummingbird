package scanner

import (
	"bufio"
	"hummingbird/internal/models"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// ScanFunctions traverses the codebase directory to discover all unique function definitions.
// It recursively walks through the directory, skipping non-source files, and extracts function names.
func ScanFunctions(root string) ([]string, error) {
	funcMap := make(map[string]struct{})

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() || !isSourceFile(path) {
			return err
		}

		contentBytes, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		cleanContent := stripComments(string(contentBytes))
		lines := strings.Split(cleanContent, "\n")

		for _, line := range lines {
			if m := funcRegex.FindStringSubmatch(line); m != nil {
				name := extractFunctionName(m)
				if name != "" {
					funcMap[name] = struct{}{}
				}
			}
		}
		return nil
	})

	funcs := make([]string, 0, len(funcMap))
	for f := range funcMap {
		funcs = append(funcs, f)
	}
	return funcs, err
}

// ScanFileContent analyzes a single file to identify references to database tables and other functions.
// It maintains a scope of the current function being analyzed to correctly attribute matches.
func ScanFileContent(path string, tables []string, functions []string) []models.Match {
	contentBytes, err := os.ReadFile(path)
	if err != nil {
		return nil
	}

	rawContent := string(contentBytes)
	cleanContent := stripComments(rawContent)

	originalLines := strings.Split(rawContent, "\n")
	cleanLines := strings.Split(cleanContent, "\n")

	matches := make([]models.Match, 0)
	currentFunc := "GlobalScope"

	for i, cleanLine := range cleanLines {
		if strings.TrimSpace(cleanLine) == "" {
			continue
		}

		// 1. Update Current Scope
		if m := funcRegex.FindStringSubmatch(cleanLine); m != nil {
			currentFunc = extractFunctionName(m)
		}

		// 2. Scan Tables
		for _, t := range tables {
			re := regexp.MustCompile(`\b` + regexp.QuoteMeta(t) + `\b`)
			if re.MatchString(cleanLine) {
				matches = append(matches, newMatch(t, "table", path, currentFunc, i, originalLines))
			}
		}

		// 3. Scan Functions
		for _, f := range functions {
			re := regexp.MustCompile(`\b` + regexp.QuoteMeta(f) + `\b`)
			// Ensure it's a call, not the definition itself
			if re.MatchString(cleanLine) && !isDefinition(cleanLine, f) {
				matches = append(matches, newMatch(f, "function", path, currentFunc, i, originalLines))
			}
		}
	}
	return matches
}

// ScanTables reads a file containing a list of table names and returns them as a slice of strings.
// It handles basic cleanup, removing quotes and common delimiters, and ignores comment lines.
func ScanTables(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var tables []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Clean quotes and commas from SQL/JSON exports
		line = strings.Trim(line, `", '`)

		if line != "" && !strings.HasPrefix(line, "#") {
			tables = append(tables, line)
		}
	}
	return tables, scanner.Err()
}
