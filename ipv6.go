package ipdb

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"net/netip"
)

var (
	_ driver.Valuer = IPv6{}
	_ sql.Scanner   = (*IPv6)(nil)

	first = IPv6From16([16]byte{
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	})
	last = IPv6From16([16]byte{
		0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff,
	})
)

func IPv6From16(ip [16]byte) IPv6 {
	return IPv6(netip.AddrFrom16(ip))
}

func MustParseIPv6(s string) IPv6 {
	ip, err := ParseIPv6(s)
	if err != nil {
		panic(err)
	}
	return ip
}

func ParseIPv6(s string) (IPv6, error) {
	ip, err := netip.ParseAddr(s)
	if err != nil {
		return IPv6{}, fmt.Errorf("invalid IPv4/IPv6 address: %q", s)
	}

	if ip.Is4() {
		ip = netip.AddrFrom16(ip.As16())
	}

	return IPv6(ip), nil
}

func (i IPv6) IsFirst() bool {
	return i.Compare(first) == 0
}

func (i IPv6) IsLast() bool {
	return i.Compare(last) == 0
}

func MustIPv6FromSlice(data []byte) IPv6 {
	ip, err := IPv6FromSlice(data)
	if err != nil {
		panic(err)
	}
	return ip
}

func IPv6FromSlice(data []byte) (IPv6, error) {
	ip, ok := netip.AddrFromSlice(data)
	if !ok {
		return IPv6{}, fmt.Errorf("invalid IPv6 address: %X", data)
	}

	if ip.Is4() {
		ip = netip.AddrFrom16(ip.As16())
	}

	return IPv6(ip), nil
}

type IPv6 netip.Addr

func (i IPv6) Value() (driver.Value, error) {
	return i.MarshalBinary()
}

func (i *IPv6) Scan(value any) error {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case []byte:
		return i.UnmarshalBinary(v)
	default:
		return fmt.Errorf("unsupported type %T", value)
	}
}

func (i IPv6) Next() IPv6 {
	return IPv6(netip.Addr(i).Next())
}

func (i IPv6) Prev() IPv6 {
	return IPv6(netip.Addr(i).Prev())
}
func (i IPv6) IsValid() bool {
	return netip.Addr(i).IsValid()
}

func (i IPv6) Is4() bool {
	return netip.Addr(i).Is4()
}

func (i IPv6) Is6() bool {
	return netip.Addr(i).Is6()
}

func (i IPv6) As16() [16]byte {
	return netip.Addr(i).As16()
}

func (i IPv6) As4() [4]byte {
	return netip.Addr(i).As4()
}

func (i IPv6) AsSlice() []byte {
	return netip.Addr(i).AsSlice()
}

func (i IPv6) IsLinkLocalMulticast() bool {
	return netip.Addr(i).IsLinkLocalMulticast()
}

func (i IPv6) IsLinkLocalUnicast() bool {
	return netip.Addr(i).IsLinkLocalUnicast()
}

func (i IPv6) IsLoopback() bool {
	return netip.Addr(i).IsLoopback()
}

func (i IPv6) IsPrivate() bool {
	return netip.Addr(i).IsPrivate()
}

func (i IPv6) IsGlobalUnicast() bool {
	return netip.Addr(i).IsGlobalUnicast()
}

func (i IPv6) IsMulticast() bool {
	return netip.Addr(i).IsMulticast()
}

func (i IPv6) Compare(other IPv6) int {
	return netip.Addr(i).Compare(netip.Addr(other))
}

func (i IPv6) String() string {
	return netip.Addr(i).String()
}

func (i IPv6) StringExpanded() string {
	return netip.Addr(i).StringExpanded()
}

func (i IPv6) Is4In6() bool {
	return netip.Addr(i).Is4In6()
}

func (i IPv6) Zone() string {
	return netip.Addr(i).Zone()
}

func (i IPv6) WithZone(zone string) IPv6 {
	return IPv6(netip.Addr(i).WithZone(zone))
}

func (i IPv6) MarshalBinary() ([]byte, error) {
	return netip.Addr(i).MarshalBinary()
}

func (i *IPv6) UnmarshalBinary(data []byte) error {
	var result netip.Addr
	if err := result.UnmarshalBinary(data); err != nil {
		return err
	}
	*i = IPv6(result)
	return nil
}

func (i IPv6) IsUnspecified() bool {
	return netip.Addr(i).IsUnspecified()
}

func (i IPv6) Less(other IPv6) bool {
	return netip.Addr(i).Less(netip.Addr(other))
}

func (i IPv6) Prefix(b int) (netip.Prefix, error) {
	return netip.Addr(i).Prefix(b)
}

func (i IPv6) MarshalText() ([]byte, error) {
	return netip.Addr(i).MarshalText()
}

func (i *IPv6) UnmarshalText(data []byte) error {
	var result netip.Addr
	if err := result.UnmarshalText(data); err != nil {
		return err
	}
	*i = IPv6(result)
	return nil
}
