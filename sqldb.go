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

type DB struct {
	Ctx    context.Context
	db     *sqlx.DB
	Logger *log.Logger
	DbType int
}

func WrapDB(ctx context.Context, db *sqlx.DB) *DB {
	return &DB{
		Ctx: ctx,
		db:  db,
	}
}

func (x *DB) Begin() (*Tx, error) {
	var err error
	tx, err := x.db.BeginTxx(x.Ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("conn Begin: %w", err)
	}
	return &Tx{
		tx:     tx,
		Logger: x.Logger,
		Ctx:    x.Ctx,
	}, nil
}

func (x *DB) log(query string, args ...any) {
	if x.Logger == nil {
		return
	}
	x.Logger.Println(sql_fake(x.DbType, query, args...))
}

func (x *DB) Select(dest any, query string, args ...any) error {
	x.log(query, args...)
	if x.db == nil {
		return fmt.Errorf("sxc: %T", x.db)
	}
	return sqlx.SelectContext(x.Ctx, x.db, dest, query, args...)
}

func (x *DB) Get(dest any, query string, args ...any) error {
	x.log(query, args...)
	return sqlx.GetContext(x.Ctx, x.db, dest, query, args...)
}

func (x *DB) MustExec(query string, args ...any) sql.Result {
	res, err := x.Exec(query, args...)
	if err != nil {
		panic(err)
	}
	return res
}
func (x *DB) Exec(query string, args ...any) (sql.Result, error) {
	x.log(query, args...)
	return x.db.ExecContext(x.Ctx, query, args...)
}

func (x *DB) InsertMap(table string, m map[string]any) (sql.Result, error) {
	s, values := insertSt(x.DbType, table, m)
	res, err := x.Exec(s, values...)
	return res, err
}

// InsertMapReturning will add returning at the end of the statement
// with returning string and call Get to dest
// dest must be a pointer to destination
// returning is the name(s) of the field(s)
func (x *DB) InsertMapReturning(dest any, returning string, table string, m map[string]any) error {
	s, values := insertSt(x.DbType, table, m)
	s += " returning " + returning
	return x.Get(dest, s, values...)
}

func (x *DB) UpdateMap(table string, m map[string]any, where string, where_vals ...any) (sql.Result, error) {
	s, values := updateSt(x.DbType, table, m, where, where_vals...)
	res, err := x.Exec(s, values...)

	return res, err
}

// UpdateMapReturning will add returning at the end of the statement
// with returning string and call Get to dest
// dest must be a pointer to destination
// returning is the name(s) of the field(s)
func (x *DB) UpdateMapReturning(dest any, returning string, table string, m map[string]any, where string, where_vals ...any) error {
	s, values := updateSt(x.DbType, table, m, where, where_vals...)
	s += " returning " + returning
	return x.Get(dest, s, values...)
}
