package ipdb

import (
	"context"
	"database/sql"
	"fmt"
)

func (i *IPDB) vicinityN(ctx context.Context, tx *sql.Tx, r IPv6Range, n uint) (below, inside, above []boundary, err error) {
	inside, err = i.inside(ctx, tx, r)
	if err != nil {
		return nil, nil, nil, err
	}

	below, err = i.belowN(ctx, tx, r.lb, n)
	if err != nil {
		return nil, nil, nil, err
	}

	if n > 0 && len(below) == 0 {
		return nil, nil, nil, fmt.Errorf("database inconsistent: no below values in vicinity of %s with n=%d", r, n)
	}

	above, err = i.aboveN(ctx, tx, r.ub, n)
	if err != nil {
		return nil, nil, nil, err
	}

	if n > 0 && len(above) == 0 {
		return nil, nil, nil, fmt.Errorf("database inconsistent: no above values in vicinity of %s with n=%d", r, n)
	}

	return below, inside, above, nil
}
