package ipdb

import (
	"context"
	"database/sql"
)

var (
	stmtGetValue = "SELECT value FROM range_value WHERE rowid = ?;"
)

// value must exist in database when we call this method!
func (i *IPDB) getValue(ctx context.Context, tx *sql.Tx, rowid int64) (value []byte, err error) {
	stmt := tx.StmtContext(ctx, i.stmtGetValue)
	row := stmt.QueryRowContext(ctx, rowid)
	err = row.Scan(&value)
	if err != nil {
		return nil, err
	}
	return value, nil
}
