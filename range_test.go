package ipdb

import (
	"testing"
)

func TestRange(t *testing.T) {
	table := []struct {
		name     string
		input    string
		expected IPv6Range
	}{
		{
			name:  "single ipv6 address",
			input: "2001:db8::1",
			expected: IPv6Range{
				lb: MustParseIPv6("2001:db8::1"),
				ub: MustParseIPv6("2001:db8::1"),
			},
		},
		{
			name:  "small iv6 range",
			input: "2001:db8::1 - 2001:db8::2",
			expected: IPv6Range{
				lb: MustParseIPv6("2001:db8::1"),
				ub: MustParseIPv6("2001:db8::2"),
			},
		},
		{
			name:  "small reverted iv6 range",
			input: "2001:db8::2 - 2001:db8::1",
			expected: IPv6Range{
				lb: MustParseIPv6("2001:db8::1"),
				ub: MustParseIPv6("2001:db8::2"),
			},
		},
		{
			name:  "single ipv4 address",
			input: "127.0.0.1",
			expected: IPv6Range{
				lb: MustParseIPv6("127.0.0.1"),
				ub: MustParseIPv6("127.0.0.1"),
			},
		},
		{
			name:  "small iv4 range",
			input: "127.0.0.1 - 127.0.0.2",
			expected: IPv6Range{
				lb: MustParseIPv6("127.0.0.1"),
				ub: MustParseIPv6("127.0.0.2"),
			},
		},
		{
			name:  "small reverted iv4 range",
			input: "127.0.0.2 - 127.0.0.1",
			expected: IPv6Range{
				lb: MustParseIPv6("127.0.0.1"),
				ub: MustParseIPv6("127.0.0.2"),
			},
		},
	}

	for _, tc := range table {
		t.Run(tc.name, func(t *testing.T) {

			r, err := ParseIPv6Range(tc.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if r != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, r)
			}
		})
	}
}

func TestIter(t *testing.T) {
	t.Parallel()

	start := MustParseIPv6("127.0.0.1")
	stop := MustParseIPv6("127.0.0.10")

	r := MustParseIPv6Range("127.0.0.1-127.0.0.10")

	cnt := -1
	var last IPv6
	seen := map[IPv6]struct{}{}
	for ip := range r.IPs() {
		if _, ok := seen[ip]; ok {
			t.Errorf("ip %s already seen", ip)
		}
		seen[ip] = struct{}{}
		last = ip
		cnt++

		if cnt == 0 {
			if ip != start {
				t.Errorf("first expected %s, got %s", start, ip)
			}
		}

	}

	if last != stop {
		t.Errorf("last expected %s, got %s", stop, last)
	}
}
