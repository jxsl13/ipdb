package ipdb

import (
	"context"
	"database/sql"
	"fmt"
)

func (i *IPDB) Insert(ctx context.Context, r IPv6Range, value []byte) error {
	if !r.IsValid() {
		return fmt.Errorf("insert failed: invalid ip range: %s", r)
	}

	return i.do(ctx, func(ctx context.Context, tx *sql.Tx) error {

		rowid, err := i.insertValue(ctx, tx, value)
		if err != nil {
			return err
		}

		belowN, inside, aboveN, err := i.vicinityN(ctx, tx, r, 1)
		if err != nil {
			return err
		}

		// remove all boundaries inside
		if len(inside) > 0 {
			if len(inside) == 1 {
				err = i.removeIP(ctx, tx, inside[0].IP)
				if err != nil {
					return err
				}
			} else {
				insideRange := IPv6Range{
					lb: inside[0].IP,             // first
					ub: inside[len(inside)-1].IP, // last
				}

				err = i.removeRange(ctx, tx, insideRange)
				if err != nil {
					return err
				}
			}
		}

		low := boundary{
			IP:           r.lb,
			BoundaryType: lb,
			ValueRowID:   rowid,
		}

		high := boundary{
			IP:           r.ub,
			BoundaryType: ub,
			ValueRowID:   rowid,
		}

		belowNearest := belowN[0]
		aboveNearest := aboveN[0]

		belowCut := boundary{
			IP:           r.lb.Prev(),
			BoundaryType: ub,
			ValueRowID:   rowid,
		}

		aboveCut := boundary{
			IP:           r.ub.Next(),
			BoundaryType: lb,
			ValueRowID:   rowid,
		}

		insertLowerBound := true
		insertUpperBound := true

		if belowNearest.IsLowerBoundary() {
			// need to cut below
			if belowNearest.IP != belowCut.IP {
				// can cut below |----
				if belowNearest.ValueRowID != rowid {
					// only insert if reasons differ
					err = i.insertBoundary(ctx, tx, belowCut)
					if err != nil {
						return err
					}
				} else {
					// extend range towards belowNearest
					insertLowerBound = false
				}
			} else {
				// cannot cut below
				if belowNearest.ValueRowID != rowid {
					// if reasons differ, make beLowNearest a single bound
					belowNearest.SetDoubleBoundary()
					err = i.insertBoundary(ctx, tx, belowNearest)
					if err != nil {
						return err
					}
				} else {
					insertLowerBound = false
				}
			}
		} else if belowNearest.IsDoubleBoundary() && belowNearest.IP == belowCut.IP && belowNearest.ValueRowID == rowid {
			// one IP below we have a single boundary range with the same reason
			belowNearest.SetLowerBoundary()
			err = i.insertBoundary(ctx, tx, belowNearest)
			if err != nil {
				return err
			}
		}

		if aboveNearest.IsUpperBoundary() {
			// need to cut above
			if aboveNearest.IP != aboveCut.IP {
				// can cut above -----|
				if aboveNearest.ValueRowID != rowid {
					// insert if reasons differ
					err = i.insertBoundary(ctx, tx, aboveCut)
					if err != nil {
						return err
					}
				} else {
					// don't insert, because extends range
					// to upperbound above
					insertUpperBound = false
				}

			} else {
				// cannot cut above
				if aboveNearest.ValueRowID != rowid {
					aboveNearest.SetDoubleBoundary()
					err = i.insertBoundary(ctx, tx, aboveNearest)
					if err != nil {
						return err
					}
				} else {
					insertUpperBound = false
				}
			}
		} else if aboveNearest.IsDoubleBoundary() && aboveNearest.IP == aboveCut.IP && aboveNearest.ValueRowID == rowid {
			// one IP above we have a single boundary range with the same reason
			aboveNearest.SetUpperBoundary()
			err = i.insertBoundary(ctx, tx, aboveNearest)
			if err != nil {
				return err
			}
		}

		if r.IsDoubleBoundary() && insertLowerBound && insertUpperBound {
			doubleBoundary := boundary{
				IP:           r.lb,
				BoundaryType: db,
				ValueRowID:   rowid,
			}
			err = i.insertBoundary(ctx, tx, doubleBoundary)
			if err != nil {
				return err
			}

		} else if insertLowerBound && insertUpperBound {
			err = i.insertBoundary(ctx, tx, low)
			if err != nil {
				return err
			}
			err = i.insertBoundary(ctx, tx, high)
			if err != nil {
				return err
			}
		} else if insertLowerBound {
			err = i.insertBoundary(ctx, tx, low)
			if err != nil {
				return err
			}
		} else if insertUpperBound {
			err = i.insertBoundary(ctx, tx, high)
			if err != nil {
				return err
			}
		}

		return nil
	})
}
