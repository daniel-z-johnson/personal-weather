package config

import (
	"os"
	"testing"
)

func TestConfig_EnvironmentOverrides(t *testing.T) {
	// Set environment variables
	os.Setenv("WEATHER_API_KEY", "test-api-key")
	os.Setenv("PORT", "8080")
	os.Setenv("DATABASE_PATH", "test.db")
	defer func() {
		os.Unsetenv("WEATHER_API_KEY")
		os.Unsetenv("PORT")
		os.Unsetenv("DATABASE_PATH")
	}()

	// Create a temporary config file
	configContent := `{
		"weatherAPI": {
			"key": "file-api-key"
		},
		"server": {
			"port": "3000"
		},
		"database": {
			"path": "file.db"
		}
	}`
	
	tmpFile, err := os.CreateTemp("", "test-config-*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	
	if _, err := tmpFile.WriteString(configContent); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	// Load config
	config, err := LoadConfig(tmpFile.Name())
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	// Verify environment variables override file values
	if config.WeatherAPI.Key != "test-api-key" {
		t.Errorf("Expected API key 'test-api-key', got '%s'", config.WeatherAPI.Key)
	}
	if config.Server.Port != "8080" {
		t.Errorf("Expected port '8080', got '%s'", config.Server.Port)
	}
	if config.Database.Path != "test.db" {
		t.Errorf("Expected database path 'test.db', got '%s'", config.Database.Path)
	}
}

func TestConfig_DefaultValues(t *testing.T) {
	// Create a minimal config file
	configContent := `{
		"weatherAPI": {
			"key": "test-key"
		}
	}`
	
	tmpFile, err := os.CreateTemp("", "test-config-*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	
	if _, err := tmpFile.WriteString(configContent); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	// Load config
	config, err := LoadConfig(tmpFile.Name())
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	// Verify default values are set
	if config.Server.Port != "1117" {
		t.Errorf("Expected default port '1117', got '%s'", config.Server.Port)
	}
	if config.Database.Path != "w.db" {
		t.Errorf("Expected default database path 'w.db', got '%s'", config.Database.Path)
	}
}