package balancer

import (
	"testing"
	"time"
)

func Test_IsWithin(t *testing.T) {
	now := time.Now()
	nowUnix := now.Unix()
	tests := []struct {
		name              string
		start             int64
		end               int64
		durationFromStart time.Duration
		expected          bool
	}{
		{
			"2 hours within a time 4 hours from now should be false",
			nowUnix,
			now.Add(4 * time.Hour).Unix(),
			2 * time.Hour,
			false,
		},
		{
			"6 hours within a time 4 hours from now should be true",
			nowUnix,
			now.Add(4 * time.Hour).Unix(),
			6 * time.Hour,
			true,
		},
		{
			"a time 2 minutes from the initial time within 6 hours should be true",
			1689631378,
			1689631497,
			6 * time.Hour,
			true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := IsWithin(test.start, test.end, test.durationFromStart)
			if test.expected != actual {
				t.Fatalf("test.expected(%v) != actual(%v)", test.expected, actual)
			}
		})
	}
}
