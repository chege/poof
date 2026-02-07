package config

type Config struct {
	Allowlist []string `hcl:"allowlist,optional"`
	Tables    []Table  `hcl:"table,block"`
}

type Table struct {
	Name    string   `hcl:"name,label"`
	PK      string   `hcl:"pk"`
	Source  *Source  `hcl:"source,block"`
	Columns []Column `hcl:"column,block"`
}

type Source struct {
	Type   string            `hcl:"type,label"`
	Name   string            `hcl:"name,optional"`
	SQL    string            `hcl:"sql,optional"`
	Params map[string]string `hcl:"params,optional"`
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
