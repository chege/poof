package generator

import (
	"crypto/md5"
	"fmt"
)

func NewRowContext(tableName, columnName string, pkValue any) RowContext {
	data := fmt.Sprintf("%s:%v", tableName, pkValue)
	seed := md5.Sum([]byte(data))
	return RowContext{
		TableName:       tableName,
		ColumnName:      columnName,
		PrimaryKeyValue: pkValue,
		Seed:            seed,
	}
}
