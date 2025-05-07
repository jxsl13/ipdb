package ipdb

import (
	"context"
	"database/sql"
)

const (
	stmtInside = "SELECT ipv6, boundary_type, range_value_rowid FROM ipv6_range WHERE ipv6 BETWEEN ? AND ? ORDER BY ipv6 ASC;"
)

func (i *IPDB) inside(ctx context.Context, tx *sql.Tx, r IPv6Range) ([]boundary, error) {

	stmt := tx.StmtContext(ctx, i.stmtInside)
	rows, err := stmt.QueryContext(ctx, r.lb, r.ub)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	boundaries := []boundary{}
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
	if err = rows.Err(); err != nil && !isNoRows(err) {
		return nil, err
	}
	return boundaries, nil
}
