package utils

import (
	"strings"
	"testing"
)

var tcGenerateTokens = []struct {
	TestName     string
	ChannelName  string
	MustContains string
}{
	{
		TestName:     "Generate token from channel with 2 words name",
		ChannelName:  "Foo Bar",
		MustContains: "foobar-whatsapp-demo-",
	},
	{
		TestName:     "Generate token from channel with 3 words name",
		ChannelName:  "Foo Bar Baz",
		MustContains: "foobarbaz-whatsapp-demo-",
	},
}

func TestGenerateTokens(t *testing.T) {
	for _, tc := range tcGenerateTokens {
		t.Run(tc.TestName, func(t *testing.T) {
			generatedToken := GenToken(tc.ChannelName)
			if !strings.Contains(generatedToken, tc.MustContains) {
				t.Errorf("got %v / must contains %v",
					generatedToken,
					tc.MustContains)
			}
		})
	}
}
