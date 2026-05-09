package resolver

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveRealBinary(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(t *testing.T) (pathEnv, skipDir string)
		binary  string
		wantErr bool
	}{
		{
			name:   "next match found",
			binary: "git",
			setup: func(t *testing.T) (string, string) {
				dir := t.TempDir()
				path := filepath.Join(dir, "git")
				if err := os.WriteFile(path, []byte("#!/bin/sh"), 0755); err != nil {
					t.Fatal(err)
				}
				return dir, "/shim"
			},
		},
		{
			name:    "no match found",
			binary:  "git",
			wantErr: true,
			setup: func(t *testing.T) (string, string) {
				return t.TempDir(), "/shim"
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pathEnv, skipDir := tt.setup(t)
			t.Setenv("PATH", pathEnv)

			got, err := ResolveRealBinary(tt.binary, skipDir)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ResolveRealBinary() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got == "" {
				t.Error("ResolveRealBinary() returned empty path, want non-empty")
			}
		})
	}
}
