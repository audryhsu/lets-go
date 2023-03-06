package main

import (
	"snippetbox.audryhsu.com/internal/assert"
	"testing"
	"time"
)

// unit test
func TestHumanDate(t *testing.T) {

	// table-driven tests to run multiple test cases by creating a slice of anon struct "cases".
	tests := []struct {
		name     string
		input    time.Time
		expected string
	}{
		{
			name:     "UTC",
			input:    time.Date(2022, 3, 17, 10, 15, 0, 0, time.UTC),
			expected: "17 Mar 2022 at 10:15",
		},
		{
			name:     "Empty time",
			input:    time.Time{},
			expected: "",
		},
		{
			name:     "CET",
			input:    time.Date(2022, 3, 17, 10, 15, 0, 0, time.FixedZone("CET", 1*60*60)),
			expected: "17 Mar 2022 at 09:15",
		},
	}

	// loop over test cases
	for _, test := range tests {
		// run subtests using t.Run() -- takes a subtest name and anon func to run actual test for each case.
		t.Run(test.name, func(t *testing.T) {
			hd := humanDate(test.input)

			// use our custom assert package
			assert.Equal(t, hd, test.expected)
		})
	}
}
