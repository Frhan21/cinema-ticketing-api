package config

import (
	"os"
	"testing"
)

func TestLoadUsesRenderPort(t *testing.T) {
	workingDirectory, err := os.Getwd()
	if err != nil {
		t.Fatalf("get working directory: %v", err)
	}
	if err := os.Chdir(t.TempDir()); err != nil {
		t.Fatalf("change working directory: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(workingDirectory); err != nil {
			t.Errorf("restore working directory: %v", err)
		}
	})

	t.Setenv("APP_PORT", ":8080")
	t.Setenv("PORT", "10000")

	cfg := Load()
	if cfg.App.Port != ":10000" {
		t.Fatalf("expected Render port :10000, got %q", cfg.App.Port)
	}
}
