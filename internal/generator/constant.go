package generator

type constantGenerator struct {
	value string
}

func NewConstantGenerator(value string) Generator {
	return &constantGenerator{value: value}
}

func (g *constantGenerator) Generate(ctx RowContext) (any, error) {
	return g.value, nil
}
