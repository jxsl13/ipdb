package ipdb

import (
	"context"
	"database/sql"
)

const (
	stmtInsertBoundary = "INSERT INTO ipv6_range (ipv6, boundary_type, range_value_rowid) VALUES (?, ?, ?);"
)

func (i *IPDB) insertBoundary(ctx context.Context, tx *sql.Tx, b boundary) error {
	// Insert the boundary into the database
	stmt := tx.StmtContext(ctx, i.stmtInsertBoundary)
	_, err := stmt.ExecContext(ctx, b.IP, b.BoundaryType, b.ValueRowID)
	return err
}
