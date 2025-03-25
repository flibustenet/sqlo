// Copyright (c) 2025 William Dode
// Licensed under the MIT license. See LICENSE file in the project root for details.
package sqlo

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

type Tx struct {
	Ctx    context.Context
	tx     *sqlx.Tx
	Logger *log.Logger
	DbType int
}

func WrapTx(ctx context.Context, tx *sqlx.Tx) *Tx {
	return &Tx{
		Ctx: ctx,
		tx:  tx,
	}
}
func (x *Tx) Commit() error {
	err := x.tx.Commit()
	if err != nil {
		return fmt.Errorf("Tx Commit: %w", err)
	}
	return nil
}

func (x *Tx) Rollback() error {
	err := x.tx.Rollback()
	if err != nil {
		return fmt.Errorf("Tx Rollback: %w", err)
	}
	return nil
}

func (x *Tx) log(query string, args ...any) {
	if x.Logger == nil {
		return
	}
	x.Logger.Println(sql_fake(x.DbType, query, args...))
}

func (x *Tx) Select(dest any, query string, args ...any) error {
	x.log(query, args...)
	return sqlx.SelectContext(x.Ctx, x.tx, dest, query, args...)
}

func (x *Tx) Get(dest any, query string, args ...any) error {
	x.log(query, args...)
	return sqlx.GetContext(x.Ctx, x.tx, dest, query, args...)
}

func (x *Tx) MustExec(query string, args ...any) sql.Result {
	res, err := x.Exec(query, args...)
	if err != nil {
		panic(err)
	}
	return res
}
func (x *Tx) Exec(query string, args ...any) (sql.Result, error) {
	x.log(query, args...)
	if x.tx == nil {
		return nil, fmt.Errorf("sxc: %T", x.tx)
	}
	return x.tx.ExecContext(x.Ctx, query, args...)
}

func (x *Tx) InsertMap(table string, m map[string]any) (sql.Result, error) {
	s, values := insertSt(x.DbType, table, m)
	res, err := x.Exec(s, values...)
	return res, err
}

// InsertMapReturning will add returning at the end of the statement
// with returning string and call Get to dest
// dest must be a pointer to destination
// returning is the name(s) of the field(s)
func (x *Tx) InsertMapReturning(dest any, returning string, table string, m map[string]any) error {
	s, values := insertSt(x.DbType, table, m)
	s += " returning " + returning
	return x.Get(dest, s, values...)
}

func (x *Tx) UpdateMap(table string, m map[string]any, where string, where_vals ...any) (sql.Result, error) {
	s, values := updateSt(x.DbType, table, m, where, where_vals...)
	res, err := x.Exec(s, values...)

	return res, err
}

// UpdateMapReturning will add returning at the end of the statement
// with returning string and call Get to dest
// dest must be a pointer to destination
// returning is the name(s) of the field(s)
func (x *Tx) UpdateMapReturning(dest any, returning string, table string, m map[string]any, where string, where_vals ...any) error {
	s, values := updateSt(x.DbType, table, m, where, where_vals...)
	s += " returning " + returning
	return x.Get(dest, s, values...)
}
