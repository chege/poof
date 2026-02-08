package engine

import (
	"fmt"
	"strings"
)

// buildBatchUpdateQuery generates an UPDATE ... FROM (VALUES ...) query for PostgreSQL.
// The query expects the primary key to be the first element in each row's values.
func (e *Engine) buildBatchUpdateQuery(tableName, pk string, columns []string, columnTypes []string, batchSize int) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("UPDATE %s AS t SET ", tableName))

	for i, col := range columns {
		if i > 0 {
			sb.WriteString(", ")
		}
		// Offset by 1 because pk is at index 0 in the source values
		sb.WriteString(fmt.Sprintf("%s = v.col%d", col, i+1))
	}

	sb.WriteString(" FROM (VALUES ")

	numCols := len(columns) + 1 // +1 for the PK
	for r := 0; r < batchSize; r++ {
		if r > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString("(")
		for c := 0; c < numCols; c++ {
			if c > 0 {
				sb.WriteString(", ")
			}
			// Placeholder index starts at 1
			placeholder := fmt.Sprintf("$%d", r*numCols+c+1)

			// Add explicit casting for the first row to help Postgres infer types
			if r == 0 && len(columnTypes) > c && columnTypes[c] != "" {
				sb.WriteString(fmt.Sprintf("CAST(%s AS %s)", placeholder, columnTypes[c]))
			} else {
				sb.WriteString(placeholder)
			}
		}
		sb.WriteString(")")
	}

	sb.WriteString(") AS v(pk")
	for i := range columns {
		sb.WriteString(fmt.Sprintf(", col%d", i+1))
	}
	sb.WriteString(") WHERE t." + pk + " = v.pk")

	return sb.String()
}
