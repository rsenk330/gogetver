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
