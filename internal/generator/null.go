package generator

type nullGenerator struct{}

func NewNullGenerator() Generator {
	return &nullGenerator{}
}

func (g *nullGenerator) Generate(ctx RowContext) (any, error) {
	return nil, nil
}
