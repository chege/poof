// Package analyze provides heuristic analysis of database schemas to suggest masking rules.
package analyze

import (
	"context"
	"fmt"

	"github.com/christopher/poof/internal/db"
	"github.com/christopher/poof/internal/ui"
)

// Suggestion represents a recommended masking configuration for a column.
type Suggestion struct {
	TableName  string
	ColumnName string
	Generator  string
	Provider   string
	Reason     string
}

// Analyzer performs heuristic analysis on a database schema.
type Analyzer struct {
	DB db.DB
}

// NewAnalyzer creates a new analyzer for the given database.
func NewAnalyzer(database db.DB) *Analyzer {
	return &Analyzer{DB: database}
}

// Analyze inspects the database and returns masking suggestions.
func (a *Analyzer) Analyze(ctx context.Context) ([]Suggestion, error) {
	tables, err := a.DB.GetAllTables(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting tables: %w", err)
	}

	// 1. Heuristic PII Analysis
	var suggestions []Suggestion
	for _, tableName := range tables {
		columns, err := a.DB.GetTableColumns(ctx, tableName)
		if err != nil {
			return nil, fmt.Errorf("getting columns for %s: %w", tableName, err)
		}

		for _, col := range columns {
			// Never suggest masking the primary key or foreign keys
			if col.IsPrimaryKey || col.IsForeignKey {
				continue
			}

			for _, rule := range DefaultRules {
				if rule.Regex.MatchString(col.Name) {
					suggestions = append(suggestions, Suggestion{
						TableName:  tableName,
						ColumnName: col.Name,
						Generator:  rule.Generator,
						Provider:   rule.Provider,
						Reason:     fmt.Sprintf("Matched rule '%s'", rule.Name),
					})
					break // Only one rule per column
				}
			}
		}
	}

	// 2. Circular Dependency Detection
	a.detectCycles(ctx, tables)

	return suggestions, nil
}

func (a *Analyzer) detectCycles(ctx context.Context, tables []string) {
	adj := make(map[string][]string)
	for _, table := range tables {
		fks, err := a.DB.GetForeignKeys(ctx, table)
		if err != nil {
			continue
		}
		for _, fk := range fks {
			adj[table] = append(adj[table], fk.ReferencedTableName)
		}
	}

	visited := make(map[string]int) // 0: unvisited, 1: visiting, 2: visited
	for _, table := range tables {
		if visited[table] == 0 {
			if a.hasCycle(table, adj, visited) {
				ui.Warning("Circular dependency detected involving table %q. Review masking order.", table)
			}
		}
	}
}

func (a *Analyzer) hasCycle(u string, adj map[string][]string, visited map[string]int) bool {
	visited[u] = 1
	for _, v := range adj[u] {
		if visited[v] == 1 {
			return true
		}
		if visited[v] == 0 && a.hasCycle(v, adj, visited) {
			return true
		}
	}
	visited[u] = 2
	return false
}
