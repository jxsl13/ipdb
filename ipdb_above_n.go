package ipdb

import (
	"context"
	"database/sql"
)

const (
	stmtAboveN = "SELECT ipv6, boundary_type, range_value_rowid FROM ipv6_range WHERE ipv6 > ? ORDER BY ipv6 ASC LIMIT ?;"
)

func (i *IPDB) aboveN(ctx context.Context, tx *sql.Tx, ip IPv6, n uint) ([]boundary, error) {
	stmt := tx.StmtContext(ctx, i.stmtAboveN)
	rows, err := stmt.QueryContext(ctx, ip, n)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	boundaries := make([]boundary, 0, min(64, n))
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

	// we must always have at least the last IPv6 available containing only ones
	// which is why at least one row is always returned.
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return boundaries, nil
}
