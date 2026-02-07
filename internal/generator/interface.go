package generator

type RowContext struct {
	TableName        string
	ColumnName       string
	PrimaryKeyValue  any
	Seed             [16]byte // MD5 hash
}

type Generator interface {
	Generate(ctx RowContext) (any, error)
}
