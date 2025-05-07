package ipdb

import (
	"context"
	"database/sql"
)

const (
	stmtSelectValue = "SELECT rowid FROM range_value WHERE value = ?;"
	stmtInsertValue = "INSERT INTO range_value (value) VALUES (?);"
)

// only insert value if it does not already exist
func (i *IPDB) insertValue(ctx context.Context, tx *sql.Tx, value []byte) (rowid int64, err error) {

	// TODO: might speed this up with INSERT OR ABORT or INSERT OR FAIL
	stmt := tx.StmtContext(ctx, i.stmtSelectValue)
	row := stmt.QueryRowContext(ctx, value)
	err = row.Scan(&rowid)
	if err == nil {
		// already exists, return rowid of existing value
		return rowid, nil
	}
	if !isNoRows(err) {
		// unexpected error
		return 0, err
	}

	// not found, insert new value
	stmt = tx.StmtContext(ctx, i.stmtInsertValue)
	res, err := stmt.ExecContext(ctx, value)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}
