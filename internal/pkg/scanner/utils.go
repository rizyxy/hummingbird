package scanner

import (
	"hummingbird/internal/models"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	funcRegex    = regexp.MustCompile(`(?i)(?:func|def|function|fn|fun)\s+([a-zA-Z0-9_]+)|(?:const|let|var)\s+([a-zA-Z0-9_]+)\s*=\s*(?:\([^)]*\)|[a-zA-Z0-9_]+)\s*=>`)
	validExts    = map[string]bool{".go": true, ".js": true, ".py": true, ".cs": true, ".ts": true, ".tsx": true, ".sql": true}
	commentRegex = regexp.MustCompile(`(?s)//.*?\n|/\*.*?\*/`)
)

func isValidName(name string) bool {
	if len(name) <= 3 {
		return false
	}

	for _, r := range name {
		if r < '0' || r > '9' {
			return true
		}
	}
	return false
}

func isDefinition(line, funcName string) bool {
	if m := funcRegex.FindStringSubmatch(line); m != nil {
		return extractFunctionName(m) == funcName
	}
	return false
}

func isSourceFile(path string) bool {
	return validExts[filepath.Ext(path)]
}

func extractFunctionName(matches []string) string {
	if len(matches) < 3 {
		return ""
	}
	if matches[1] != "" {
		return matches[1]
	}
	return matches[2]
}

func newMatch(name, cat, path, curFunc string, idx int, lines []string) models.Match {
	return models.Match{
		Name:         name,
		Category:     cat,
		FileName:     filepath.Base(path),
		FunctionName: curFunc,
		LineNumber:   idx + 1,
		Snippet:      captureSnippet(lines, idx, 1),
	}
}

func captureSnippet(lines []string, idx, radius int) string {
	start := max(0, idx-radius)
	end := min(len(lines)-1, idx+radius)

	var sb strings.Builder
	for i := start; i <= end; i++ {
		if i == idx {
			sb.WriteString("> ")
		} else {
			sb.WriteString("  ")
		}
		sb.WriteString(strings.TrimSpace(lines[i]))
		if i < end {
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

func stripComments(content string) string {
	return commentRegex.ReplaceAllStringFunc(content, func(s string) string {
		res := make([]byte, len(s))
		for i := 0; i < len(s); i++ {
			if s[i] == '\n' {
				res[i] = '\n'
			} else {
				res[i] = ' '
			}
		}
		return string(res)
	})
}
