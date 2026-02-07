package config

type Config struct {
	Allowlist []string `hcl:"allowlist,optional"`
	Tables    []Table  `hcl:"table,block"`
}

type Table struct {
	Name    string   `hcl:"name,label"`
	PK      string   `hcl:"pk"`
	Columns []Column `hcl:"column,block"`
}

type Column struct {
	Name string `hcl:"name,label"`
	Gen  Gen    `hcl:"gen,block"`
}

type Gen struct {
	Type     string            `hcl:"type,label"`
	Provider string            `hcl:"provider,optional"`
	Value    string            `hcl:"value,optional"`
	Template string            `hcl:"template,optional"`
	Params   map[string]string `hcl:"params,optional"`
}
