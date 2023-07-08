package valuegen

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testCase struct {
	url      string
	expected string
}

func TestGenerateValue(t *testing.T) {
	testCases := []testCase{
		{
			url:      `https://youtube.com`,
			expected: "AvdMs6o9",
		},
		{
			url:      `https://sreda.v-a-c.org/en/read-11`,
			expected: "gmjHKtpU",
		},
		{
			url:      `https://ru.wikipedia.org/wiki/%D0%9C%D0%BD%D0%BE%D0%B6%D0%B5%D1%81%D1%82%D0%B2%D0%BE`,
			expected: "gkTF8orL",
		},
	}
	for _, test := range testCases {
		assert.Equal(t, GenerateValue(test.url), test.expected)
	}
}
