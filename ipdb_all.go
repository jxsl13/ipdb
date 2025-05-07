package ipdb

import (
	"context"
	"database/sql"
	"fmt"
)

const (
	stmtAll = "SELECT ipv6, boundary_type, range_value_rowid FROM ipv6_range ORDER BY ipv6 ASC"
)

func (i *IPDB) all(ctx context.Context, tx *sql.Tx) (_ []boundary, err error) {
	stmt := tx.StmtContext(ctx, i.stmtAll)
	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	boundaries := make([]boundary, 0, 2)
	for rows.Next() {
		var b boundary
		err = rows.Scan(
			&b.IP,
			&b.BoundaryType,
			&b.ValueRowID,
		)
		if err != nil {
			return nil, err
		}
		boundaries = append(boundaries, b)
	}

	// must always contain at least two elements in the database
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return boundaries, nil
}

func (i *IPDB) ValidateConsistency(ctx context.Context) error {
	return i.do(ctx, func(ctx context.Context, tx *sql.Tx) error {
		return i.consistent(ctx, tx)
	})
}

func (i *IPDB) consistent(ctx context.Context, tx *sql.Tx, ipRange ...IPv6Range) error {

	boundaries, err := i.all(ctx, tx)
	if err != nil {
		return err
	}

	if len(ipRange) > 0 {
		r := ipRange[0]
		low := r.First()
		high := r.Last()

		foundLow, foundHigh := false, false
		for _, b := range boundaries {
			if b.IP == low && b.IsLowerBoundary() {
				foundLow = true
			}

			if b.IP == high && b.IsUpperBoundary() {
				foundHigh = true
			}
		}
		if !foundLow || !foundHigh {
			if !foundLow && !foundHigh {
				return fmt.Errorf("did neither find inserted LOWERBOUND neither UPPERBOUND")
			} else if !foundLow {
				return fmt.Errorf("did not find inserted LOWERBOUND")
			}
			return fmt.Errorf("did not find inserted UPPERBOUND")
		}
	}

	cnt := 0
	state := lb
	for idx, b := range boundaries {

		if b.IsDoubleBoundary() {
			if state != ub {
				return fmt.Errorf("database inconsistent: double boundary: idx=%d state=%d, expected state=%s", idx, state, ub)
			}

			cnt += 2
		} else if b.IsLowerBoundary() {
			if state != ub {
				return fmt.Errorf("database inconsistent: lower boundary: idx=%d state=%d, expected state=%s", idx, state, ub)
			}
			cnt++
			state = boundaryType(cnt % 2)
		} else if b.IsUpperBoundary() {
			if state != lb {
				return fmt.Errorf("database inconsistent: upper boundary: idx=%d state=%d, expected state=%s", idx, state, lb)
			}

			if idx > 0 {
				currentValueID := b.ValueRowID
				prevValueID := boundaries[idx-1].ValueRowID

				// reasons consistent
				if currentValueID != prevValueID {
					return fmt.Errorf("reason mismatch: idx=%4d rowid=%d - idx=%4d rowid=%d", idx-1, prevValueID, idx, currentValueID)
				}
			}

			cnt++
			state = boundaryType(cnt % 2)
		}
	}

	if state != lb {
		return fmt.Errorf("database inconsistent: final boundary is supposed to be a lower boundary")
	}

	return nil
}
