package sqlo

import (
	"fmt"
	"sort"
	"strings"
)

type Raw string

// renvoi la chaine sql et les valeurs pour un insert
// à partir d'un map
func insertSt(dbType int, table string, m map[string]any) (string, []any) {
	fieldols := make([]string, 0)
	values := make([]any, 0)
	fieldnames := make([]string, 0)
	for name := range m {
		fieldnames = append(fieldnames, name)
	}
	sort.Strings(fieldnames)

	i := 0
	for _, name := range fieldnames {
		if v, ok := m[name].(Raw); ok {
			fieldols = append(fieldols, string(v))
			continue
		}
		switch dbType {
		case DB_ACCESS:
			fieldols = append(fieldols, "?")
		case DB_MSSQL:
			fieldols = append(fieldols, fmt.Sprintf("@p%d", i+1))
		default: //pg
			fieldols = append(fieldols, fmt.Sprintf("$%d", i+1))
		}
		i++
		values = append(values, m[name])
	}
	s := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		table,
		strings.Join(fieldnames, ", "),
		strings.Join(fieldols, ", "))
	return s, values
}

// renvoi la chaine sql et les valeurs pour un update
// à partir d'un map
func updateSt(dbType int, table string, m map[string]any, where string, where_vals ...any) (string, []any) {
	sets := make([]string, 0)
	num := len(where_vals) + 1
	values := []any{}
	if dbType != DB_ACCESS { // si type $1 $2... on met les vals en premier sinon en dernier
		values = where_vals[:]
	}

	fieldnames := make([]string, 0)
	for name := range m {
		fieldnames = append(fieldnames, name)
	}
	sort.Strings(fieldnames)
	for _, name := range fieldnames {
		if _, ok := m[name].(Raw); ok {
			sets = append(sets, fmt.Sprintf("%s=%s", name, m[name]))
			continue
		}
		switch dbType {
		case DB_ACCESS:
			sets = append(sets, fmt.Sprintf("%s=?", name))
		case DB_MSSQL:
			sets = append(sets, fmt.Sprintf("%s=@p%d", name, num))
		default:
			sets = append(sets, fmt.Sprintf("%s=$%d", name, num))
		}
		num += 1
		values = append(values, m[name])
	}
	s := fmt.Sprintf("UPDATE %s SET %s WHERE %s",
		table,
		strings.Join(sets, ", "),
		where)

	if dbType == DB_ACCESS {
		values = append(values, where_vals...)
	}
	return s, values
}
