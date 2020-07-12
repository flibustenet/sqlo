package sqlo

import (
	"testing"
)

func TestWhereDol(t *testing.T) {
	type D struct {
		sql   string
		args  []int
		query string
	}
	d := D{"a=%s and b=%s", []int{1, 2}, " where a=$1 and b=$2"}
	where := &Where{}
	where.Appendf("a=%s", 1)
	where.Appendf("b=%s", 2)
	res := where.Where()
	if res != d.query {
		t.Errorf("de %s %v attend %s reçoit %s", d.sql, d.args, d.query, res)
	}
	where = &Where{}
	where.Appendf("a=%s and b=%s", 1, 2)
	res = where.Where()
	if res != d.query {
		t.Errorf("de %s %v attend %s reçoit %s", d.sql, d.args, d.query, res)
	}

	d = D{"a=%s and b=%s", []int{1, 2}, " where x=x and a=$1 and b=$2"}
	where = &Where{}
	where.Appendf("x=x")
	where.Appendf("a=%s and b=%s", 1, 2)
	res = where.Where()
	if res != d.query {
		t.Errorf("de %s %v attend %s reçoit %s", d.sql, d.args, d.query, res)
	}

	d = D{"a=%s and b=%s", []int{1, 2}, " where x=$1 and a=$2 and b=$3"}
	where = &Where{}
	where.Appendf("x=$1")
	where.Appendf("", "x")
	where.Appendf("a=%s and b=%s", 1, 2)
	res = where.Where()
	if res != d.query {
		t.Errorf("de %s %v attend %s reçoit %s", d.sql, d.args, d.query, res)
	}

	d = D{"a=%s and b in (%s,%s) and c=%s", []int{1, 2, 3, 4}, " where a=$1 and b in ($2,$3) and c=$4"}
	where = &Where{}
	where.Appendf("a=%s", 1)
	where.AppendListf("b in (%s)", 2, 3)
	where.Appendf("c=%s", 4)
	res = where.Where()
	if res != d.query {
		t.Errorf("de %s %v attend %s reçoit %s", d.sql, d.args, d.query, res)
	}
}
func TestWhereSQLServer(t *testing.T) {
	type D struct {
		sql   string
		args  []int
		query string
	}
	d := D{"a=%s and b=%s", []int{1, 2}, " where a=@p1 and b=@p2"}
	where := &Where{}
	where.Style = "@p"
	where.Appendf("a=%s", 1)
	where.Appendf("b=%s", 2)
	res := where.Where()
	if res != d.query {
		t.Errorf("de %s %v attend %s reçoit %s", d.sql, d.args, d.query, res)
	}
	where = &Where{}
	where.Style = "@p"
	where.Appendf("a=%s and b=%s", 1, 2)
	res = where.Where()
	if res != d.query {
		t.Errorf("de %s %v attend %s reçoit %s", d.sql, d.args, d.query, res)
	}

	d = D{"a=%s and b=%s", []int{1, 2}, " where x=x and a=@p1 and b=@p2"}
	where = &Where{}
	where.Style = "@p"
	where.Appendf("x=x")
	where.Appendf("a=%s and b=%s", 1, 2)
	res = where.Where()
	if res != d.query {
		t.Errorf("de %s %v attend %s reçoit %s", d.sql, d.args, d.query, res)
	}

	d = D{"a=%s and b=%s", []int{1, 2}, " where x=@p1 and a=@p2 and b=@p3"}
	where = &Where{}
	where.Style = "@p"
	where.Appendf("x=@p1")
	where.Appendf("", "x")
	where.Appendf("a=%s and b=%s", 1, 2)
	res = where.Where()
	if res != d.query {
		t.Errorf("de %s %v attend %s reçoit %s", d.sql, d.args, d.query, res)
	}

	d = D{"a=%s and b in (%s,%s) and c=%s", []int{1, 2, 3, 4}, " where a=@p1 and b in (@p2,@p3) and c=@p4"}
	where = &Where{}
	where.Style = "@p"
	where.Appendf("a=%s", 1)
	where.AppendListf("b in (%s)", 2, 3)
	where.Appendf("c=%s", 4)
	res = where.Where()
	if res != d.query {
		t.Errorf("de %s %v attend %s reçoit %s", d.sql, d.args, d.query, res)
	}
}
func TestWhereQMark(t *testing.T) {
	type D struct {
		sql   string
		args  []int
		query string
	}
	d := D{"a=%s and b=%s", []int{1, 2}, " where a=? and b=?"}
	where := &Where{}
	where.Style = "?"
	where.Appendf("a=%s", 1)
	where.Appendf("b=%s", 2)
	res := where.Where()
	if res != d.query {
		t.Errorf("de %s %v attend %s reçoit %s", d.sql, d.args, d.query, res)
	}
	where = &Where{}
	where.Style = "?"
	where.Appendf("a=%s and b=%s", 1, 2)
	res = where.Where()
	if res != d.query {
		t.Errorf("de %s %v attend %s reçoit %s", d.sql, d.args, d.query, res)
	}

	d = D{"a=%s and b=%s", []int{1, 2}, " where x=x and a=? and b=?"}
	where = &Where{}
	where.Style = "?"
	where.Appendf("x=x")
	where.Appendf("a=%s and b=%s", 1, 2)
	res = where.Where()
	if res != d.query {
		t.Errorf("de %s %v attend %s reçoit %s", d.sql, d.args, d.query, res)
	}

	d = D{"a=%s and b=%s", []int{1, 2}, " where x=? and a=? and b=?"}
	where = &Where{}
	where.Style = "?"
	where.Appendf("x=?")
	where.Appendf("", "x")
	where.Appendf("a=%s and b=%s", 1, 2)
	res = where.Where()
	if res != d.query {
		t.Errorf("de %s %v attend %s reçoit %s", d.sql, d.args, d.query, res)
	}

	d = D{"a=%s and b in (%s,%s) and c=%s", []int{1, 2, 3, 4}, " where a=? and b in (?,?) and c=?"}
	where = &Where{}
	where.Style = "?"
	where.Appendf("a=%s", 1)
	where.AppendListf("b in (%s)", 2, 3)
	where.Appendf("c=%s", 4)
	res = where.Where()
	if res != d.query {
		t.Errorf("de %s %v attend %s reçoit %s", d.sql, d.args, d.query, res)
	}
}
