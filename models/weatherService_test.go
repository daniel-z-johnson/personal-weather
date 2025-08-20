package models

import (
	"database/sql"
	"log/slog"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestWeatherService_SaveAndDeleteLocation(t *testing.T) {
	// Create in-memory database for testing
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}
	defer db.Close()

	// Create table
	_, err = db.Exec(`CREATE TABLE locations (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		city TEXT NOT NULL,
		state TEXT NOT NULL DEFAULT '',
		country TEXT NOT NULL DEFAULT '',
		latitude REAL NOT NULL DEFAULT 0.0,
		longitude REAL NOT NULL DEFAULT 0.0,
		temp REAL NOT NULL DEFAULT 0.0,
		expires TEXT NOT NULL DEFAULT 0.0
	)`)
	if err != nil {
		t.Fatalf("Failed to create test table: %v", err)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	ws := &WeatherService{DB: db, Logger: logger}

	// Test SaveLocation
	err = ws.SaveLocation("Test City", "TX", "US", 30.0, -97.0)
	if err != nil {
		t.Errorf("SaveLocation failed: %v", err)
	}

	// Test GetLocation
	location, err := ws.GetLocation("Test City", "TX", "US")
	if err != nil {
		t.Errorf("GetLocation failed: %v", err)
	}
	if location == nil {
		t.Error("Expected location, got nil")
	}
	if location.Name != "Test City" {
		t.Errorf("Expected city 'Test City', got '%s'", location.Name)
	}

	// Test DeleteLocation
	err = ws.DeleteLocation("Test City", "TX", "US")
	if err != nil {
		t.Errorf("DeleteLocation failed: %v", err)
	}

	// Verify deletion
	location, err = ws.GetLocation("Test City", "TX", "US")
	if err != nil {
		t.Errorf("GetLocation after delete failed: %v", err)
	}
	if location != nil {
		t.Error("Expected nil location after deletion, but got location")
	}
}

func TestWeatherService_InputValidation(t *testing.T) {
	// Create in-memory database for testing
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}
	defer db.Close()

	// Create table
	_, err = db.Exec(`CREATE TABLE locations (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		city TEXT NOT NULL,
		state TEXT NOT NULL DEFAULT '',
		country TEXT NOT NULL DEFAULT '',
		latitude REAL NOT NULL DEFAULT 0.0,
		longitude REAL NOT NULL DEFAULT 0.0,
		temp REAL NOT NULL DEFAULT 0.0,
		expires TEXT NOT NULL DEFAULT 0.0
	)`)
	if err != nil {
		t.Fatalf("Failed to create test table: %v", err)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	ws := &WeatherService{DB: db, Logger: logger}

	// Test with valid coordinates
	err = ws.SaveLocation("Valid City", "CA", "US", 37.7749, -122.4194)
	if err != nil {
		t.Errorf("SaveLocation with valid coordinates failed: %v", err)
	}

	// Test deletion of non-existent location
	err = ws.DeleteLocation("Non-existent City", "XX", "XX")
	if err == nil {
		t.Error("Expected error when deleting non-existent location, but got nil")
	}
}