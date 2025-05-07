package ipdb

import (
	"context"
	"database/sql"
	"errors"
)

func (i *IPDB) do(ctx context.Context, fn func(ctx context.Context, tx *sql.Tx) error) error {
	tx, err := i.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		err = errors.Join(err, tx.Rollback())
	}()

	err = fn(ctx, tx)
	if err != nil {
		return err
	}

	return tx.Commit()
}
