package ipdb

import (
	"context"
	"database/sql"
	"errors"

	_ "modernc.org/sqlite"
)

type IPDB struct {
	db *sql.DB

	stmtAboveN              *sql.Stmt
	stmtAll                 *sql.Stmt
	stmtBelowN              *sql.Stmt
	stmtGetValue            *sql.Stmt
	stmtInsertBoundary      *sql.Stmt
	stmtSelectValue         *sql.Stmt
	stmtInsertValue         *sql.Stmt
	stmtInside              *sql.Stmt
	stmtCleanupUnreferenced *sql.Stmt
	stmtRemoveRange         *sql.Stmt
	stmtRemoveIP            *sql.Stmt
}

type Option func(*options) error

type options struct{}

func (i *IPDB) Close() (err error) {
	defer func() {
		err = errors.Join(err, i.db.Close())
	}()

	return errors.Join(
		i.stmtAboveN.Close(),
		i.stmtAll.Close(),
		i.stmtBelowN.Close(),
		i.stmtGetValue.Close(),
		i.stmtInsertBoundary.Close(),
		i.stmtSelectValue.Close(),
		i.stmtInsertValue.Close(),
		i.stmtInside.Close(),
		i.stmtCleanupUnreferenced.Close(),
		i.stmtRemoveRange.Close(),
		i.stmtRemoveIP.Close(),
	)
}

func Open(ctx context.Context, dsnURI string, opts ...Option) (_ *IPDB, err error) {
	db, err := sql.Open("sqlite", dsnURI)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			err = errors.Join(err, db.Close())
		}
	}()

	options := &options{}

	for _, opt := range opts {
		if err := opt(options); err != nil {
			return nil, err
		}
	}

	ipdb := &IPDB{
		db: db,
	}

	err = ipdb.prepare(ctx)
	if err != nil {
		return nil, err
	}

	const stmt = `
PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS range_value (
	rowid 	INTEGER NOT NULL PRIMARY KEY,
	value 	BLOB UNIQUE
);

-- boundary 0 is a lower boundary, 1 is an upper boundary, 2 is lower and upper boundary = double boundary (single ip range)
CREATE TABLE IF NOT EXISTS ipv6_range (
	ipv6	 			BLOB UNIQUE NOT NULL,
	boundary_type 		INTEGER NOT NULL,
	range_value_rowid 	INTEGER,
	FOREIGN KEY (range_value_rowid)
		REFERENCES
			range_value(rowid) ON DELETE CASCADE
);

INSERT OR REPLACE
	INTO range_value (value)
	VALUES
		('-inf'),
		('+inf')
;

INSERT OR REPLACE
	INTO ipv6_range (ipv6, boundary_type, range_value_rowid)
	VALUES
		(?, ?, 1),
		(?, ?, 2)
;
`
	_, err = db.ExecContext(
		ctx, stmt,
		first, ub, // -inf upper bondary, first ipv6
		last, lb, // +inf lower bondary, last ipv6
	)
	if err != nil {
		return nil, err
	}

	return ipdb, nil
}

func (i *IPDB) prepare(ctx context.Context) (err error) {
	i.stmtAboveN, err = i.db.PrepareContext(ctx, stmtAboveN)
	if err != nil {
		return err
	}
	i.stmtAll, err = i.db.PrepareContext(ctx, stmtAll)
	if err != nil {
		return err
	}
	i.stmtBelowN, err = i.db.PrepareContext(ctx, stmtBelowN)
	if err != nil {
		return err
	}
	i.stmtGetValue, err = i.db.PrepareContext(ctx, stmtGetValue)
	if err != nil {
		return err
	}
	i.stmtInsertBoundary, err = i.db.PrepareContext(ctx, stmtInsertBoundary)
	if err != nil {
		return err
	}
	i.stmtSelectValue, err = i.db.PrepareContext(ctx, stmtSelectValue)
	if err != nil {
		return err
	}
	i.stmtInsertValue, err = i.db.PrepareContext(ctx, stmtInsertValue)
	if err != nil {
		return err
	}
	i.stmtInside, err = i.db.PrepareContext(ctx, stmtInside)
	if err != nil {
		return err
	}
	i.stmtCleanupUnreferenced, err = i.db.PrepareContext(ctx, stmtCleanupUnreferenced)
	if err != nil {
		return err
	}
	i.stmtRemoveRange, err = i.db.PrepareContext(ctx, stmtRemoveRange)
	if err != nil {
		return err
	}
	i.stmtRemoveIP, err = i.db.PrepareContext(ctx, stmtRemoveIP)
	if err != nil {
		return err
	}

	return nil
}
