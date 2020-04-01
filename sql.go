package sqlo

import (
	"database/sql"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/jmoiron/sqlx"
)

type Selecter interface {
	Select(interface{}, string, ...interface{}) error
	Get(interface{}, string, ...interface{}) error
}

type Execer interface {
	Selecter
	Exec(string, ...interface{}) (sql.Result, error)
	MustExec(string, ...interface{}) sql.Result
	NamedExec(string, interface{}) (sql.Result, error)
	InsertMap(string, map[string]interface{}) (sql.Result, error)
	UpdateMap(string, map[string]interface{}, string, ...interface{}) (sql.Result, error)
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

func (x *Sx) log(query string, args ...interface{}) {
	if x.Logger == nil {
		return
	}
	x.Logger.Println(sql_fake(x.DbType, query, args...))
}

func (x *Sx) Select(dest interface{}, query string, args ...interface{}) error {
	x.log(query, args...)
	return sqlx.Select(x.Sx, dest, query, args...)
}

func (x *Sx) Get(dest interface{}, query string, args ...interface{}) error {
	x.log(query, args...)
	return sqlx.Get(x.Sx, dest, query, args...)
}

func (x *Sx) MustExec(query string, args ...interface{}) sql.Result {
	res, err := x.Exec(query, args...)
	if err != nil {
		panic(err)
	}
	return res
}

func (x *Sx) Exec(query string, args ...interface{}) (sql.Result, error) {
	x.log(query, args...)
	return x.Sx.Exec(query, args...)
}

func (x *Sx) NamedExec(query string, arg interface{}) (sql.Result, error) {
	x.log(query, arg)
	return sqlx.NamedExec(x.Sx, query, arg)
}

// renvoi la chaine sql et les valeurs pour un insert
// à partir d'un map
func (x *Sx) insertSt(table string, m map[string]interface{}) (string, []interface{}) {
	fieldols := make([]string, 0)
	values := make([]interface{}, 0)
	fieldnames := make([]string, 0)
	for name, _ := range m {
		fieldnames = append(fieldnames, name)
	}
	sort.Strings(fieldnames)

	for i, name := range fieldnames {
		if x.DbType == DB_ACCESS {
			fieldols = append(fieldols, "?")
		} else {
			fieldols = append(fieldols, fmt.Sprintf("$%d", i+1))
		}
		values = append(values, m[name])
	}
	s := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		table,
		strings.Join(fieldnames, ", "),
		strings.Join(fieldols, ", "))
	return s, values
}
func (x *Sx) InsertMap(table string, m map[string]interface{}) (sql.Result, error) {
	s, values := x.insertSt(table, m)
	res, err := x.Exec(s, values...)
	return res, err
}

// renvoi la chaine sql et les valeurs pour un update
// à partir d'un map
func (x *Sx) updateSt(table string, m map[string]interface{}, where string, where_vals ...interface{}) (string, []interface{}) {
	sets := make([]string, 0)
	num := len(where_vals) + 1
	values := []interface{}{}
	if x.DbType != DB_ACCESS { // si type $1 $2... on met les vals en premier sinon en dernier
		values = where_vals[:]
	}

	fieldnames := make([]string, 0)
	for name, _ := range m {
		fieldnames = append(fieldnames, name)
	}
	sort.Strings(fieldnames)
	for _, name := range fieldnames {
		if x.DbType == DB_ACCESS {
			sets = append(sets, fmt.Sprintf("%s=?", name))
		} else {
			sets = append(sets, fmt.Sprintf("%s=$%d", name, num))
		}
		num += 1
		values = append(values, m[name])
	}
	s := fmt.Sprintf("UPDATE %s SET %s WHERE %s",
		table,
		strings.Join(sets, ", "),
		where)

	if x.DbType == DB_ACCESS {
		values = append(values, where_vals...)
	}
	return s, values
}
func (x *Sx) UpdateMap(table string, m map[string]interface{}, where string, where_vals ...interface{}) (sql.Result, error) {
	s, values := x.updateSt(table, m, where, where_vals...)
	res, err := x.Exec(s, values...)

	return res, err
}
