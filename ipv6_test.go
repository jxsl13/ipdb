package ipdb

import "testing"

func TestIPv6(t *testing.T) {
	table := []struct {
		name     string
		input    string
		expected IPv6
	}{
		{
			name:  "ipv4",
			input: "127.0.0.1",
			expected: MustIPv6FromSlice([]byte{
				0, 0, 0, 0,
				0, 0, 0, 0,
				0, 0, 0xff, 0xff,
				127, 0, 0, 1,
			}),
		},
	}

	for _, tc := range table {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ip, err := ParseIPv6(tc.input)
			if err != nil {
				t.Fatal(err)
			}

			if ip != tc.expected {
				t.Fatalf("expected %s, got %s", tc.expected, ip)
			}
		})
	}
}
