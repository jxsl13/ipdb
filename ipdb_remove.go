package ipdb

import (
	"context"
	"database/sql"
	"fmt"
)

const (
	stmtCleanupUnreferenced = "DELETE FROM range_value rv WHERE NOT EXISTS (SELECT 1 FROM ipv6_range ir WHERE ir.range_value_rowid = rv.rowid);"
	stmtRemoveRange         = "DELETE FROM ipv6_range WHERE ipv6 BETWEEN ? AND ?;"
	stmtRemoveIP            = "DELETE FROM ipv6_range WHERE ipv6 = ?;"
)

func (i *IPDB) Remove(ctx context.Context, r IPv6Range) error {
	if !r.IsValid() {
		return fmt.Errorf("remove failed: invalid ip range: %s", r)
	}

	return i.do(ctx, func(ctx context.Context, tx *sql.Tx) error {
		below, inside, above, err := i.vicinityN(ctx, tx, r, 1)
		if err != nil {
			return err
		}

		if len(inside) > 0 {
			err = i.removeRange(ctx, tx, IPv6Range{lb: inside[0].IP, ub: inside[len(inside)-1].IP})
			if err != nil {
				return err
			}
		}

		low := r.First()
		high := r.Last()

		belowNearest := below[0]
		aboveNearest := above[0]

		belowCut := boundary{
			IP:           low.Prev(),
			BoundaryType: ub,
			ValueRowID:   belowNearest.ValueRowID,
		}

		aboveCut := boundary{
			IP:           high.Next(),
			BoundaryType: lb,
			ValueRowID:   aboveNearest.ValueRowID,
		}

		if belowNearest.IsLowerBoundary() {
			// need to cut below
			if belowNearest.IP != belowCut.IP {
				// can cut
				err = i.insertBoundary(ctx, tx, belowCut)
				if err != nil {
					return err
				}
			} else {
				// cannot cut
				belowNearest.SetDoubleBoundary()
				err = i.insertBoundary(ctx, tx, belowNearest)
				if err != nil {
					return err
				}
			}
		}

		if aboveNearest.IsUpperBoundary() {
			// need to cut above
			if aboveNearest.IP != aboveCut.IP {
				// can cut above
				err = i.insertBoundary(ctx, tx, aboveCut)
				if err != nil {
					return err
				}
			} else {
				// cannot cut above
				aboveNearest.SetDoubleBoundary()
				err = i.insertBoundary(ctx, tx, aboveNearest)
				if err != nil {
					return err
				}

			}
		}

		return nil
	})
}

func (i *IPDB) cleanupUnreferenced(ctx context.Context, tx *sql.Tx) error {
	stmt := tx.StmtContext(ctx, i.stmtCleanupUnreferenced)
	_, err := stmt.ExecContext(ctx)
	return err
}

func (i *IPDB) removeRange(ctx context.Context, tx *sql.Tx, r IPv6Range) error {
	if r.IsDoubleBoundary() {
		return i.removeIP(ctx, tx, r.First())
	}

	stmt := tx.StmtContext(ctx, i.stmtRemoveRange)
	_, err := stmt.ExecContext(ctx, r.lb, r.ub)
	if err != nil {
		return err
	}

	return i.cleanupUnreferenced(ctx, tx)
}

func (i *IPDB) removeIP(ctx context.Context, tx *sql.Tx, ip IPv6) error {
	stmt := tx.StmtContext(ctx, i.stmtRemoveIP)
	_, err := stmt.ExecContext(ctx, ip)
	if err != nil {
		return err
	}

	return i.cleanupUnreferenced(ctx, tx)
}
