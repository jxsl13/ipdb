package ipdb

import (
	"errors"
	"fmt"
	"iter"
	"net/netip"
	"strings"

	"github.com/gaissmai/extnetip"
)

var (
	errInvalidRange = errors.New("invalid range")
)

func MustParseIPv6Range(s string) IPv6Range {
	r, err := ParseIPv6Range(s)
	if err != nil {
		panic(err)
	}
	return r
}

func ParseIPv6Range(s string) (r IPv6Range, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("%w: %w", errInvalidRange, err)
		}
	}()

	if strings.Contains(s, "-") {
		parts := strings.SplitN(s, "-", 2)
		if len(parts) != 2 {
			return r, errors.New("contains '-' but not exactly two parts")
		}

		lb, err := ParseIPv6(strings.TrimSpace(parts[0]))
		if err != nil {
			return r, err
		}

		ub, err := ParseIPv6(strings.TrimSpace(parts[1]))
		if err != nil {
			return r, err
		}
		cmp := lb.Compare(ub)
		if cmp == 0 {
			return IPv6Range{
				lb: lb,
				ub: ub,
			}, nil
		} else if 0 < cmp {
			return IPv6Range{
				lb: ub,
				ub: lb,
			}, nil
		} else {
			return IPv6Range{
				lb: lb,
				ub: ub,
			}, nil
		}

	} else if strings.Contains(s, "/") {
		prefix, err := netip.ParsePrefix(s)
		if err != nil {
			return r, err
		}
		lb, ub := extnetip.Range(prefix)
		return IPv6Range{
			lb: IPv6(lb),
			ub: IPv6(ub),
		}, nil
	}

	// single ip range
	ip, err := ParseIPv6(s)
	if err != nil {
		return r, err
	}
	return IPv6Range{
		lb: ip,
		ub: ip,
	}, nil
}

func RangeFromPrefix(cidr netip.Prefix) IPv6Range {
	first, last := extnetip.Range(cidr)
	return IPv6Range{
		lb: IPv6(first),
		ub: IPv6(last),
	}
}

func RangeFromIPv6(lb, ub IPv6) IPv6Range {
	if lb.Compare(ub) > 0 {
		return IPv6Range{
			lb: ub,
			ub: lb,
		}
	}
	return IPv6Range{
		lb: lb,
		ub: ub,
	}
}

type IPv6Range struct {
	lb IPv6
	ub IPv6
}

func (r IPv6Range) IsValid() bool {
	return r.lb.Less(r.ub) && !(r.lb.IsFirst() && r.lb.IsLast() && r.ub.IsFirst() && r.ub.IsLast())
}

func (i IPv6Range) IPs() iter.Seq[IPv6] {
	current := i.lb
	stop := i.ub.Next() // exclusive
	return func(yield func(IPv6) bool) {
		for ; current.Compare(stop) < 0; current = current.Next() {
			if !yield(current) {
				return
			}
		}
	}
}

func (i IPv6Range) String() string {
	return i.lb.String() + " - " + i.ub.String()
}

func (i IPv6Range) StringExpanded() string {
	return i.lb.StringExpanded() + " - " + i.ub.StringExpanded()
}

func (i IPv6Range) Contains(ip IPv6) bool {
	return i.lb.Compare(ip) <= 0 && i.ub.Compare(ip) >= 0
}

func (i IPv6Range) IsDoubleBoundary() bool {
	return i.lb.Compare(i.ub) == 0
}

func (i IPv6Range) First() IPv6 {
	return i.lb
}

func (i IPv6Range) Last() IPv6 {
	return i.ub
}
