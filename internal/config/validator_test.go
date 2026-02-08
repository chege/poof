package config

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockGeneratorValidator struct {
	validProviders map[string]bool
}

func (m *mockGeneratorValidator) ValidateGen(_ string, gen Gen, sqlDataType string) error {
	switch gen.Type {
	case "faker":
		if !m.validProviders[gen.Provider] {
			return fmt.Errorf("unknown faker provider %q", gen.Provider)
		}
		if sqlDataType == "integer" && gen.Provider == "email" {
			return fmt.Errorf("type mismatch")
		}
	case "template":
		if gen.Template == "invalid" {
			return fmt.Errorf("invalid template")
		}
	}
	return nil
}

func TestValidateStatic(t *testing.T) {
	gv := &mockGeneratorValidator{
		validProviders: map[string]bool{"email": true},
	}

	tests := []struct {
		cfg     *Config
		name    string
		msg     string
		wantErr bool
	}{
		{
			name: "Valid static config",
			cfg: &Config{
				Tables: []Table{
					{
						Name: "users",
						Columns: []Column{
							{Name: "email", Gen: Gen{Type: "faker", Provider: "email"}},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Unknown faker provider",
			cfg: &Config{
				Tables: []Table{
					{
						Name: "users",
						Columns: []Column{
							{Name: "email", Gen: Gen{Type: "faker", Provider: "unknown_provider"}},
						},
					},
				},
			},
			wantErr: true,
			msg:     "unknown faker provider",
		},
		{
			name: "Invalid template syntax",
			cfg: &Config{
				Tables: []Table{
					{
						Name: "users",
						Columns: []Column{
							{Name: "email", Gen: Gen{Type: "template", Template: "invalid"}},
						},
					},
				},
			},
			wantErr: true,
			msg:     "invalid template",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.ValidateStatic(gv)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.msg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
