package main

import (
	"os"
	"testing"
)

func setupEnv() {
	testEnv := []struct {
		name string
		val  string
	}{
		{"HOST", "localhost"},
		{"DEBUG", "true"},
		{"DEBUGG", "1234"},
	}
	// Set initial environment
	for _, e := range testEnv {
		os.Setenv(e.name, e.val)
	}
}

func TestGetenv(t *testing.T) {
	setupEnv()

	tests := []struct {
		name     string
		def      string
		expected string
	}{
		{"undef", "def", "def"},
		{"HOST", "127.0.0.1", "localhost"},
		{"host", "127.0.0.1", "127.0.0.1"},
	}

	for _, test := range tests {
		name := test.name
		expected := test.expected
		val := Getenv(name, test.def)
		if val != expected {
			t.Errorf("Exected %v=%v, got %v.", name, expected, val)
		}
	}
}

func TestGetenvBool(t *testing.T) {
	setupEnv()

	tests := []struct {
		name     string
		def      bool
		expected bool
	}{
		{"undef", false, false},
		{"DEBUG", false, true},
		{"DEBUGG", false, false},
		{"debug", true, true},
	}

	for _, test := range tests {
		name := test.name
		expected := test.expected
		val := GetenvBool(name, test.def)
		if val != expected {
			t.Errorf("Exected %v=%v, got %v.", name, expected, val)
		}
	}
}

func TestNewAppConfig(t *testing.T) {
	cfg := NewAppConfig()

	if cfg.Hostname != "gogetver.com" {
		t.Errorf("Expected Hostname=gogetver.com, got %v", cfg.Hostname)
	}
	if cfg.IP != "127.0.0.1" {
		t.Errorf("Expected IP=127.0.0.1, got %v", cfg.IP)
	}
	if cfg.Port != "5000" {
		t.Errorf("Expected Port=5000, got %v", cfg.Port)
	}
	if !cfg.Debug {
		t.Errorf("Expected Debug=false, got %v", cfg.Debug)
	}
	if cfg.TemplatesDir != "./templates" {
		t.Errorf("Expected TemplatesDir=./templates, got %v", cfg.TemplatesDir)
	}
	if cfg.GoogleAnalyticsID != "" {
		t.Errorf("Expected GoogleAnalyticsID=\"\", got %v", cfg.GoogleAnalyticsID)
	}
}
