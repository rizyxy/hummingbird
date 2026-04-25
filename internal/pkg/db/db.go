package db

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

// FetchTables connects to the database and retrieves a list of all table names
func FetchTables(driver, dsn string) ([]string, error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open db connection: %w", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	var query string
	switch strings.ToLower(driver) {
	case "postgres":
		query = "SELECT table_name FROM information_schema.tables WHERE table_schema='public';"
	case "mysql":
		query = "SELECT table_name FROM information_schema.tables WHERE table_schema=DATABASE();"
	default:
		return nil, fmt.Errorf("unsupported database driver: %s. Supported drivers are 'postgres' and 'mysql'", driver)
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query tables: %w", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, fmt.Errorf("failed to scan table name: %w", err)
		}
		tables = append(tables, tableName)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return tables, nil
}
