// Package config package config
package config

import "testing"

func Test_parseConfig(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "parse config file",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parseConfig()
		})
	}
}
