package parser

import (
	"bufio"
	"os"
	"strings"
)

func LoadTables(filename string) ([]string, error) {
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
