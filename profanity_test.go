package main

import (
	"fmt"
	"testing"
)

func TestProfanityCheck(t *testing.T) {
	cases := []struct {
		input string
		expected string
	}{
		{
			input: "What in the KerFuffle is this?",
			expected: "What in the **** is this?",
		},
		{
			input: "I need you to calm the fornax down.",
			expected: "I need you to calm the **** down.",
		},
		{
			input: "SHARBERT this isn't good",
			expected: "**** this isn't good",
		},
		{
			input: "what the fornax!",
			expected: "what the fornax!",
		},
		{
			input: "Hello, there!",
			expected: "Hello, there!",
		},
	}

	for i, c := range cases{
		t.Run(fmt.Sprintf("Test case %v",i), func(t *testing.T){
			result := replace_profanity(c.input)
			if result != c.expected{
				t.Errorf("Failed Test %d: '%s' and '%s' don't match\n",i+1,result,c.expected)
				return
			}
		})
	}

	
}