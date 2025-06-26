package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {

	cases := []struct {
		input    string
		expected []string
	}{
		{input: "hello world",
			expected: []string{"hello", "world"},
		},
		{
			input:    "le test",
			expected: []string{"le", "test"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("el largo no coincide")
			return
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]

			if word != expectedWord {
				t.Errorf("las palabras no coinciden")
			}
		}

	}

}
