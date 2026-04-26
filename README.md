# 🐦 Hummingbird | Migration Intelligence Engine

Hummingbird is a fast, lightweight codebase analysis tool written in Go. It maps database dependencies and logic paths to identify high-friction areas and mitigate risks during system migrations.

By scanning your source code and a provided list of target tables, Hummingbird helps architects and developers visualize dependencies, calculate the impact ("blast radius") of changes, and confidently plan complex system migrations.

## 🚀 Features

- **Dependency Discovery:** Scans codebases (`.go`, `.js`, `.ts`, etc.) to trace logic and database usage.
- **Dynamic DB Integration:** Connects directly to PostgreSQL or MySQL to fetch table names dynamically.
- **Blast Radius Calculation:** Calculates the recursive impact of modifying a specific table.
- **Visual Architecture Graphs:** Generates Mermaid.js diagrams to visualize logic calls and data flow.
- **Strategic Summaries:** Provides prioritized CLI reports detailing the frequency of table usage and logic dependencies.
- **Standalone Binary:** Pure Go implementation that compiles to a single executable without external runtime dependencies.

## 🗄️ Supported Databases

When using the dynamic DB integration feature (`--db-driver`), Hummingbird currently supports fetching tables directly from:

- **PostgreSQL** (`postgres`)
- **MySQL / MariaDB** (`mysql`)

## 📦 Installation

Ensure you have Go installed (version 1.25+ recommended).

Navigate to the source directory and build the binary:

```bash
cd hummingbird
go build -o hummingbird main.go
```

To make it globally accessible, move the executable to your path, or use `go install`:

```bash
go install
```

## 🛠 Usage

Hummingbird takes an optional tables file and a target codebase path.

```bash
hummingbird [flags] [tables_file] <codebase_path>
```

### Arguments:

- `[tables_file]` : _(Optional)_ Path to a `.txt` file containing target table names (one per line).
- `<codebase_path>` : Directory containing the source code to audit.

### Flags:

| Flag                    | Description                                                                | Default    |
| :---------------------- | :------------------------------------------------------------------------- | :--------- |
| `--cli`                 | Print prioritized Strategic Summary and Logic Call tables to the terminal. | `false`    |
| `--graph`               | Generate separated Mermaid JS files for visualization.                     | `false`    |
| `--target <table_name>` | Calculate the recursive "Blast Radius" for a specific table.               | `""`       |
| `--out <dir>`           | Directory to save generated Mermaid files.                                 | `diagrams` |
| `--db-driver`           | Database driver to connect to (`postgres`, `mysql`).                       | `""`       |
| `--db-dsn`              | Database connection string (DSN) to fetch tables dynamically.              | `""`       |

## 💡 Examples

### 1. Standard Audit

Analyze a codebase against a list of tables and print a prioritized summary to the terminal:

```bash
hummingbird --cli tables.txt ./src
```

### 2. Generate Visualizations

Generate Mermaid architecture graphs (`architecture_logic.mmd`, `architecture_data.mmd`) for review:

```bash
hummingbird --graph tables.txt ./src
```

_Note: The generated files will be saved in the `diagrams/` folder by default. Use `--out` to change the destination._

### 3. Calculate Blast Radius

Identify all functions and downstream logic affected if a critical table is modified:

```bash
hummingbird --target TBL_USER_01_ALT tables.txt ./src
```

### 4. Codebase Scanning Only

You can also run Hummingbird without a tables file to just map internal function calls:

```bash
hummingbird --cli ./src
```

### 5. Dynamic Database Table Fetching

Fetch table names directly from a live database instead of using a local file:

```bash
hummingbird --cli --db-driver postgres --db-dsn "postgres://user:pass@localhost:5432/mydb?sslmode=disable" ./src
```

## 📝 Example `tables.txt`

```text
TBL_USER_01_ALT
ord_mstr_final
GLOBAL_SETTINGS_LOGS_DATA
LOG_AUDIT_X7
```
