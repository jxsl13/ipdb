package ipdb

import (
	"context"
	"database/sql"
	"fmt"
)

func (i *IPDB) Contains(ctx context.Context, ip IPv6) (value []byte, found bool, err error) {
	if ip.IsFirst() || ip.IsLast() {
		return nil, false, fmt.Errorf("contains failed: ip %s is not valid", ip)
	}

	err = i.do(ctx, func(ctx context.Context, tx *sql.Tx) error {
		below, inside, above, err := i.vicinityN(ctx, tx, IPv6Range{ip, ip}, 1)
		if err != nil {
			return err
		}

		if len(inside) == 1 {
			// set value return variable
			value, err = i.getValue(ctx, tx, inside[0].ValueRowID)
			if err != nil {
				return nil
			}
			found = true
			return nil
		}

		belowNearest := below[0]
		aboveNearest := above[0]

		if belowNearest.IsLowerBoundary() && aboveNearest.IsUpperBoundary() {
			if belowNearest.ValueRowID == aboveNearest.ValueRowID {
				value, err = i.getValue(ctx, tx, belowNearest.ValueRowID)
				if err != nil {
					return nil
				}
				found = true
				return nil
			}
			return fmt.Errorf("database values inconsistent: belowNearest %s (rowid=%d), aboveNearest %s (rowid=%d)",
				belowNearest.IP,
				belowNearest.ValueRowID,
				aboveNearest.IP,
				aboveNearest.ValueRowID,
			)
		}

		// not found at this point
		found = false
		return nil
	})
	if err != nil {
		return nil, false, err
	}

	return value, found, nil
}
