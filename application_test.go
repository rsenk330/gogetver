package main

import (
	"sort"
	"testing"
)

func TestByLength(t *testing.T) {
	expected := []string{
		"1",
		"123",
		"1234",
		"12345",
		"123456",
		"1234567",
		"12345678",
		"123456789",
	}

	testArr := []string{
		"123456",
		"123456789",
		"1234",
		"12345",
		"1",
		"12345678",
		"123",
		"1234567",
	}

	sort.Sort(ByLength(testArr))
	for idx, val := range testArr {
		if val != expected[idx] {
			t.Errorf("Expected ordered list, got %v", testArr)
			break
		}
	}
}

func TestPossibleVersions(t *testing.T) {
	tests := []struct {
		url      string
		versions []string
	}{
		{"github.com/rsenk330/gogetver", []string{"master"}},
		{"github.com/rsenk330/gogetver.develop", []string{"develop"}},
		{"github.com/rsenk330/gogetver.v1", []string{"v1"}},
		{"github.com/rsenk330/gogetver.v1.0", []string{"v1.0", "0"}},
		{"github.com/rsenk330/gogetver.v1.0a", []string{"v1.0a", "0a"}},
	}

	for _, test := range tests {
		versions := PossibleVersions(test.url)

		for idx, eVer := range test.versions {
			if versions[idx] != eVer {
				t.Errorf("Expected versions to equal %v, got %v", test.versions, versions)
				break
			}
		}
	}
}
