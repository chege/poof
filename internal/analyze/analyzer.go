// Package analyze provides heuristic analysis of database schemas to suggest masking rules.
package analyze

import (
	"context"
	"fmt"

	"github.com/christopher/poof/internal/db"
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

	var suggestions []Suggestion
	for _, tableName := range tables {
		columns, err := a.DB.GetTableColumns(ctx, tableName)
		if err != nil {
			return nil, fmt.Errorf("getting columns for %s: %w", tableName, err)
		}

		for _, col := range columns {
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

	return suggestions, nil
}
