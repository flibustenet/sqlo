package sqlo

import (
	"database/sql"
	"testing"
	"time"

	"github.com/lib/pq"
)

func Test_sql_quote(t *testing.T) {
	type Tst struct {
		T int
		V interface{}
		S string
	}
	tbl := []Tst{Tst{DB_ACCESS, 5, "5"},
		Tst{DB_ACCESS, 3.1415, "3.1415"},
		Tst{DB_ACCESS, 3., "3"},
		Tst{DB_ACCESS, "abcd", "'abcd'"},
		Tst{DB_ACCESS, "ab'cd", "'ab''cd'"},
		Tst{DB_ACCESS, true, "-1"},
		Tst{DB_ACCESS, false, "0"},
		Tst{DB_ACCESS, nil, "null"},
		Tst{DB_PG, pq.NullTime{}, "null"},
		Tst{DB_PG, sql.NullTime{}, "null"},
		Tst{DB_ACCESS, time.Date(1969, 11, 05, 23, 05, 03, 0, time.Local), "'1969-11-05 23:05:03'"},
		Tst{DB_PG, time.Date(1969, 11, 05, 23, 05, 03, 0, time.Local), "'1969-11-05 23:05:03'"},
		Tst{DB_MSSQL, time.Date(1969, 11, 05, 23, 05, 03, 0, time.Local), "'1969-11-05 23:05:03'"},
		Tst{DB_PG, sql.NullBool{}, "null"},
		Tst{DB_PG, sql.NullBool{true, true}, "true"},
		Tst{DB_ACCESS, sql.NullBool{true, true}, "-1"},
		Tst{DB_ACCESS, sql.NullBool{false, true}, "0"},
		Tst{DB_ACCESS, sql.NullInt64{0, false}, "null"},
		Tst{DB_ACCESS, sql.NullInt64{42, true}, "42"},
		Tst{DB_ACCESS, sql.NullInt32{42, true}, "42"},
		Tst{DB_ACCESS, sql.NullInt16{42, true}, "42"},
		Tst{DB_PG, sql.NullInt64{42, true}, "42"},
		Tst{DB_PG, sql.NullFloat64{42.42, true}, "42.42"},
		Tst{DB_PG, sql.NullFloat64{42.42, false}, "null"},
		Tst{DB_PG, Raw("now()"), "now()"},
	}
	for _, s := range tbl {
		r := sql_quoter(s.T, s.V)
		if r != s.S {
			t.Errorf("attend %s reçoit %s", s.S, r)
		}
	}
}
func Test_sql_quote_query(t *testing.T) {
	type Tst struct {
		T int
		Q string
		V []interface{}
		S string
	}
	tbl := []Tst{
		Tst{DB_ACCESS, "? ? ? ? ? ?", []interface{}{5, "abcd", "e'fg", true, false, nil}, "5 'abcd' 'e''fg' -1 0 null"},
		Tst{DB_ACCESS, "update xyz set a=?, b=? where c=?", []interface{}{5, "abcd", "e'fg"}, "update xyz set a=5, b='abcd' where c='e''fg'"},
		Tst{DB_PG, "$1, $2, $3 $4 $5 $6", []interface{}{5, "abcd", "e'fg", true, false, nil}, "5, 'abcd', 'e''fg' true false null"},
		Tst{DB_PG, "$1 $3 $2 $3", []interface{}{5, "abcd", "e'fg"}, "5 'e''fg' 'abcd' 'e''fg'"},
		Tst{DB_MSSQL, "@p1, @p2, @p3 @p4 @p5 @p6", []interface{}{5, "abcd", "e'fg", true, false, nil}, "5, 'abcd', 'e''fg' true false null"},
		Tst{DB_MSSQL, "@p1 @p3 @p2 @p3", []interface{}{5, "abcd", "e'fg"}, "5 'e''fg' 'abcd' 'e''fg'"},
		Tst{DB_MSSQL, "@p1 @p2 @p3 @p4", []interface{}{
			pq.NullTime{},
			pq.NullTime{Valid: true, Time: time.Date(2019, 1, 2, 0, 0, 0, 0, time.Local)},
			sql.NullTime{},
			sql.NullTime{Valid: true, Time: time.Date(2019, 1, 2, 0, 0, 0, 0, time.Local)},
		}, "null '2019-01-02 00:00:00' null '2019-01-02 00:00:00'"},
		Tst{DB_PG, "$1 $3 $2 $3", []interface{}{5, Raw("now()"), "e'fg"}, "5 'e''fg' now() 'e''fg'"},
	}
	for _, s := range tbl {
		r := sql_fake(s.T, s.Q, s.V...)
		if r != s.S {
			t.Errorf("type %d : %s attend %s reçoit %s", s.T, s.Q, s.S, r)
		}
	}
}
