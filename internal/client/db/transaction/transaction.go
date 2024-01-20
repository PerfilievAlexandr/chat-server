package transaction

import (
	"chat-server/internal/client/db"
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
)

type key string

const (
	TxKey key = "tx"
)

type manager struct {
	db db.Transactor
}

func New(ctx context.Context, db db.Transactor) db.TxManager {
	return &manager{db}
}

func (m manager) transaction(ctx context.Context, options pgx.TxOptions, fn db.TxHandler) error {
	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	if ok {
		return fn(ctx)
	}

	tx, err := m.db.BeginTx(ctx, options)
	if err != nil {
		return errors.Wrap(err, "can't begin transaction")
	}

	ctx = context.WithValue(ctx, TxKey, tx)

	defer func() {
		if r := recover(); r != nil {
			err = errors.Errorf("panic recovered: %v", r)
		}

		if err != nil {
			if errRollback := tx.Rollback(ctx); errRollback != nil {
				err = errors.Wrapf(err, "errRollback: %v", errRollback)
			}

			return
		}

		err := tx.Commit(ctx)
		if err != nil {
			err = errors.Wrap(err, "tx commit failed")
		}
	}()

	if err = fn(ctx); err != nil {
		err = errors.Wrap(err, "failed executing code inside transaction")
	}

	return err
}

func (m manager) ReadCommitted(ctx context.Context, fn db.TxHandler) error {
	options := pgx.TxOptions{IsoLevel: pgx.ReadCommitted}

	return m.transaction(ctx, options, fn)
}
