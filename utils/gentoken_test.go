package utils

import (
	"strings"
	"testing"
)

var tcGenerateTokens = []struct {
	TestName       string
	MustContains   string
	ExpectedLength int
}{
	{
		TestName:       "Generate token to channel",
		MustContains:   crumb,
		ExpectedLength: sufixLength + len(crumb) + 1,
	},
}

func TestGenerateTokens(t *testing.T) {
	for _, tc := range tcGenerateTokens {
		t.Run(tc.TestName, func(t *testing.T) {
			generatedToken := GenToken()
			if !strings.Contains(generatedToken, tc.MustContains) {
				t.Errorf("got %v / must contains %v",
					generatedToken,
					tc.MustContains)
			}
			if len(generatedToken) != tc.ExpectedLength {
				t.Errorf("got token with length %v / expected %v",
					len(generatedToken),
					tc.ExpectedLength)
			}
		})
	}
}
