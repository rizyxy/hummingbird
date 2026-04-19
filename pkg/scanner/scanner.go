package scanner

import (
	"bufio"
	"hummingbird/internal/models"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var traditionalFuncRegex = regexp.MustCompile(`(?i)(func|def|function|fn|fun)\s+([a-zA-Z0-9_]+)`)
var arrowFuncRegex = regexp.MustCompile(`(?i)(const|let|var)\s+([a-zA-Z0-9_]+)\s*=\s*(\([^)]*\)|[a-zA-Z0-9_]+)\s*=>`)

func SurveyFunctions(root string) ([]string, error) {
	funcMap := make(map[string]bool)
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() || !isSourceFile(path) {
			return err
		}
		content, _ := os.ReadFile(path)
		code := string(content)

		// Find traditional functions
		tradMatches := traditionalFuncRegex.FindAllStringSubmatch(code, -1)
		for _, m := range tradMatches {
			if len(m) > 2 && isValidName(m[2]) {
				funcMap[m[2]] = true
			}
		}

		// Find arrow functions
		arrowMatches := arrowFuncRegex.FindAllStringSubmatch(code, -1)
		for _, m := range arrowMatches {
			if len(m) > 2 && isValidName(m[2]) {
				funcMap[m[2]] = true
			}
		}

		return nil
	})
	var funcs []string
	for f := range funcMap {
		funcs = append(funcs, f)
	}
	return funcs, err
}

func ScanFileContent(path string, tables []string, functions []string) []models.Match {
	file, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	var matches []models.Match
	currentFunc := "Global Scope"

	for i, line := range lines {
		// Update context using both regexes
		if m := traditionalFuncRegex.FindStringSubmatch(line); len(m) > 2 {
			currentFunc = m[2]
		} else if m := arrowFuncRegex.FindStringSubmatch(line); len(m) > 2 {
			currentFunc = m[2]
		}

		// Match Tables
		for _, t := range tables {
			if strings.Contains(line, t) {
				matches = append(matches, models.Match{
					Name: t, Category: "table", FileName: filepath.Base(path),
					FunctionName: currentFunc, LineNumber: i + 1, Snippet: captureSnippet(lines, i, 1),
				})
			}
		}

		// Match Function Calls
		for _, f := range functions {
			// Avoid matching the line where the function is defined
			if strings.Contains(line, f) && !isDefinition(line, f) {
				matches = append(matches, models.Match{
					Name: f, Category: "function", FileName: filepath.Base(path),
					FunctionName: currentFunc, LineNumber: i + 1, Snippet: captureSnippet(lines, i, 1),
				})
			}
		}
	}
	return matches
}

// Helpers
func isValidName(name string) bool {
	return len(name) > 3 && !isNumeric(name)
}

func isNumeric(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func isDefinition(line, funcName string) bool {
	return strings.Contains(line, "func "+funcName) ||
		strings.Contains(line, "def "+funcName) ||
		strings.Contains(line, "function "+funcName) ||
		(strings.Contains(line, funcName) && strings.Contains(line, "=>"))
}

func captureSnippet(lines []string, idx, radius int) string {
	start, end := idx-radius, idx+radius
	if start < 0 {
		start = 0
	}
	if end >= len(lines) {
		end = len(lines) - 1
	}
	var res []string
	for i := start; i <= end; i++ {
		prefix := "  "
		if i == idx {
			prefix = "> "
		}
		res = append(res, prefix+strings.TrimSpace(lines[i]))
	}
	return strings.Join(res, "\n")
}

func isSourceFile(path string) bool {
	ext := filepath.Ext(path)
	valid := map[string]bool{".go": true, ".js": true, ".py": true, ".cs": true, ".ts": true, ".tsx": true, ".sql": true}
	return valid[ext]
}
