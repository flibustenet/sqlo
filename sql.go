// Copyright (c) 2023 William Dode
// Licensed under the MIT license. See LICENSE file in the project root for details.

package sqlo

import (
	"database/sql"
	"log"

	"github.com/jmoiron/sqlx"
)

type Selecter interface {
	Select(any, string, ...any) error
	Get(any, string, ...any) error
}

type Execer interface {
	Selecter
	Exec(string, ...any) (sql.Result, error)
	MustExec(string, ...any) sql.Result
	//	NamedExec(string, any) (sql.Result, error)
	InsertMap(string, map[string]any) (sql.Result, error)
	InsertMapReturning(any, string, string, map[string]any) error
	UpdateMap(string, map[string]any, string, ...any) (sql.Result, error)
}

type Sx struct {
	Sx     sqlx.Ext
	Logger *log.Logger
	DbType int
}

func New(tx sqlx.Ext) *Sx {
	t := &Sx{Sx: tx}
	return t
}

func (x *Sx) log(query string, args ...any) {
	if x.Logger == nil {
		return
	}
	x.Logger.Println(sql_fake(x.DbType, query, args...))
}

func (x *Sx) Select(dest any, query string, args ...any) error {
	x.log(query, args...)
	return sqlx.Select(x.Sx, dest, query, args...)
}

func (x *Sx) Get(dest any, query string, args ...any) error {
	x.log(query, args...)
	return sqlx.Get(x.Sx, dest, query, args...)
}

func (x *Sx) MustExec(query string, args ...any) sql.Result {
	res, err := x.Exec(query, args...)
	if err != nil {
		panic(err)
	}
	return res
}

func (x *Sx) Exec(query string, args ...any) (sql.Result, error) {
	x.log(query, args...)
	return x.Sx.Exec(query, args...)
}

func (x *Sx) NamedExec(query string, arg any) (sql.Result, error) {
	x.log(query, arg)
	return sqlx.NamedExec(x.Sx, query, arg)
}

func (x *Sx) InsertMap(table string, m map[string]any) (sql.Result, error) {
	s, values := insertSt(x.DbType, table, m)
	res, err := x.Exec(s, values...)
	return res, err
}

// InsertMapReturning will add returning at the end of the statement
// with returning string and call Get to dest
// dest must be a pointer to destination
// returning is the name(s) of the field(s)
func (x *Sx) InsertMapReturning(dest any, returning string, table string, m map[string]any) error {
	s, values := insertSt(x.DbType, table, m)
	s += " returning " + returning
	return x.Get(dest, s, values...)
}

func (x *Sx) UpdateMap(table string, m map[string]any, where string, where_vals ...any) (sql.Result, error) {
	s, values := updateSt(x.DbType, table, m, where, where_vals...)
	res, err := x.Exec(s, values...)

	return res, err
}

// UpdateMapReturning will add returning at the end of the statement
// with returning string and call Get to dest
// dest must be a pointer to destination
// returning is the name(s) of the field(s)
func (x *Sx) UpdateMapReturning(dest any, returning string, table string, m map[string]any, where string, where_vals ...any) error {
	s, values := updateSt(x.DbType, table, m, where, where_vals...)
	s += " returning " + returning
	return x.Get(dest, s, values...)
}
