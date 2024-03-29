// Copyright (c) 2023 William Dode
// Licensed under the MIT license. See LICENSE file in the project root for details.

package sqlo

import (
	"database/sql"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/lib/pq"
)

var DB_PG = 0
var DB_ACCESS = 1
var DB_MSSQL = 2

var sql_log_re_question = regexp.MustCompile(`\?`)
var sql_log_re_dollar = regexp.MustCompile(`\$\d+`)
var sql_log_re_sqlserver = regexp.MustCompile(`\@p\d+`)

func sql_fake(db_type int, query string, args ...interface{}) string {
	if len(args) == 0 {
		return query
	}
	rqi := 0

	frq := func(s string) string {
		switch db_type {
		case DB_ACCESS:
			rqi++
			return sql_quoter(db_type, args[rqi-1])
		case DB_MSSQL:
			rqi, _ = strconv.Atoi(s[2:])
		default: // PG
			rqi, _ = strconv.Atoi(s[1:])
		}
		return sql_quoter(db_type, args[rqi-1])
	}
	switch db_type {
	case DB_ACCESS:
		return sql_log_re_question.ReplaceAllStringFunc(query, frq)
	case DB_PG:
		return sql_log_re_dollar.ReplaceAllStringFunc(query, frq)
	case DB_MSSQL:
		return sql_log_re_sqlserver.ReplaceAllStringFunc(query, frq)
	}
	return query
}

func sql_quoter(db_type int, s interface{}) string {
	switch v := s.(type) {
	case Raw:
		return string(v)
	case nil:
		return "null"
	case int:
		return strconv.Itoa(v)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case string:
		return "'" + strings.Replace(v, "'", "''", -1) + "'"
	case time.Time:
		return v.Format("'2006-01-02 15:04:05'")
	case pq.NullTime:
		if v.Valid {
			return v.Time.Format("'2006-01-02 15:04:05'")
		}
		return "null"
	case sql.NullTime:
		if v.Valid {
			return v.Time.Format("'2006-01-02 15:04:05'")
		}
		return "null"
	case *time.Time:
		return v.Format("'2006-01-02 15:04:05'")
	case bool:
		switch db_type {
		case DB_ACCESS:
			switch v {
			case true:
				return "-1"
			case false:
				return "0"
			}
		default:
			switch v {
			case true:
				return "true"
			case false:
				return "false"
			}
		}
	case sql.NullBool:
		if !v.Valid {
			return "null"
		}
		switch db_type {
		case DB_ACCESS:
			switch v.Bool {
			case true:
				return "-1"
			case false:
				return "0"
			}
		default:
			switch v.Bool {
			case true:
				return "true"
			case false:
				return "false"
			}
		}
	case sql.NullInt64:
		if !v.Valid {
			return "null"
		}
		return strconv.Itoa(int(v.Int64))
	case sql.NullInt32:
		if !v.Valid {
			return "null"
		}
		return strconv.Itoa(int(v.Int32))
	case sql.NullInt16:
		if !v.Valid {
			return "null"
		}
		return strconv.Itoa(int(v.Int16))
	case sql.NullFloat64:
		if !v.Valid {
			return "null"
		}
		return strconv.FormatFloat(v.Float64, 'f', -1, 64)

	}
	return fmt.Sprintf("%s", s)
}
