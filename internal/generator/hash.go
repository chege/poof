package generator

import (
	/* #nosec G501 */
	"crypto/md5"
	"fmt"
)

type hashGenerator struct{}

// NewHashGenerator creates a new generator that produces deterministic MD5 hashes.
func NewHashGenerator() Generator {
	return &hashGenerator{}
}

func (g *hashGenerator) Generate(ctx RowContext) (any, error) {
	data := []byte(fmt.Sprintf("%s:%s:%v", ctx.TableName, ctx.ColumnName, ctx.PrimaryKeyValue))
	/* #nosec G401 */
	return fmt.Sprintf("%x", md5.Sum(data)), nil
}
