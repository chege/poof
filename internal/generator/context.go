// Package generator provides the masking data generation logic.
package generator

import (
	"crypto/md5" // #nosec G501 -- MD5 is used for deterministic seeding, not for cryptographic security.
	"fmt"
)

// NewRowContext creates a new context for generating masked values for a specific row and column.
func NewRowContext(tableName, columnName, locale, salt string, pkValue, originalValue any, seedBy string, attempt int) RowContext {
	identifier := fmt.Sprintf("%s:%v", tableName, pkValue)
	if seedBy == "value" {
		identifier = fmt.Sprintf("%v", originalValue)
	}

	data := fmt.Sprintf("%s:%s:%d", salt, identifier, attempt)
	// #nosec G401 -- MD5 is used for deterministic seeding, not for cryptographic security.
	seed := md5.Sum([]byte(data))
	return RowContext{
		TableName:       tableName,
		ColumnName:      columnName,
		Locale:          locale,
		Salt:            salt,
		Attempt:         attempt,
		PrimaryKeyValue: pkValue,
		OriginalValue:   originalValue,
		Seed:            seed,
	}
}
