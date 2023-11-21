// Copyright (c) 2023 William Dode
// Licensed under the MIT license. See LICENSE file in the project root for details.

package sqlo

import (
	"log"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func Test_insertSt(t *testing.T) {
	tx := &Sx{}
	fs := map[string]interface{}{}
	fs["ok"] = "coral"
	fs["yes"] = "no"
	fs["raw"] = Raw("now()")
	q, args := tx.insertSt("mytable", fs)
	if q != "INSERT INTO mytable (ok, raw, yes) VALUES ($1, now(), $2)" {
		t.Error(q)
		t.Error(args)
	}
}

func Test_updateSt(t *testing.T) {
	tx := &Sx{}
	fs := map[string]interface{}{}
	fs["ok"] = "coral"
	fs["yes"] = "no"
	fs["raw"] = Raw("now()")
	q, args := tx.updateSt("mytable", fs, "ok=$1", "ok")
	if q != "UPDATE mytable SET ok=$2, raw=now(), yes=$3 WHERE ok=$1" {
		t.Error(q)
		t.Error(args)
	}
}

func Test_insertStMssql(t *testing.T) {
	tx := &Sx{}
	tx.DbType = DB_MSSQL
	fs := map[string]interface{}{}
	fs["ok"] = "coral"
	fs["yes"] = "no"
	q, _ := tx.insertSt("mytable", fs)
	if q != "INSERT INTO mytable (ok, yes) VALUES (@p1, @p2)" {
		log.Fatal(q)
	}
}

func Test_updateStMssql(t *testing.T) {
	tx := &Sx{}
	fs := map[string]interface{}{}
	fs["ok"] = "coral"
	fs["yes"] = "no"
	tx.DbType = DB_MSSQL
	q, _ := tx.updateSt("mytable", fs, "ok=@p1", "ok")
	if q != "UPDATE mytable SET ok=@p2, yes=@p3 WHERE ok=@p1" {
		log.Fatal(q)
	}
}

func Test_insertStAccess(t *testing.T) {
	tx := &Sx{}
	tx.DbType = DB_ACCESS

	fs := map[string]interface{}{}
	fs["ok"] = "coral"
	fs["yes"] = "no"
	q, _ := tx.insertSt("mytable", fs)
	if q != "INSERT INTO mytable (ok, yes) VALUES (?, ?)" {
		log.Fatalf("insert access : %s", q)
	}
}

func Test_updateStAccess(t *testing.T) {
	tx := &Sx{}
	tx.DbType = DB_ACCESS
	fs := map[string]interface{}{}
	fs["ok"] = "coral"
	fs["yes"] = "no"
	q, _ := tx.updateSt("mytable", fs, "ok=?", "ok")
	if q != "UPDATE mytable SET ok=?, yes=? WHERE ok=?" {
		log.Fatalf("update access : %s", q)
	}
}

func Test_insert(t *testing.T) {
	dbTest := os.Getenv("SQLO_DBTEST")
	if dbTest == "" {
		return
	}
	db, err := sqlx.Open("postgres", dbTest)
	if err != nil {
		t.Fatalf("Open dbTest: %v", err)
	}
	x := New(db)
	one := 0
	x.Get(&one, "select 1")
	if one != 1 {
		t.Errorf("one should be 1 is %d", one)
	}

	_, err = x.Exec("drop table if exists sqlo_test")
	if err != nil {
		t.Fatalf("drop sqlo_test")
	}

	_, err = x.Exec("create table sqlo_test (dedef text default 'defval', vl text)")
	if err != nil {
		t.Errorf("create table test : %v", err)
	}
	dedef := ""
	row := db.QueryRow("insert into sqlo_test (vl) values ('value') returning dedef")
	err = row.Scan(&dedef)
	if err != nil {
		t.Fatalf("queryrow test dedef error: %v", err)
	}

	_, err = x.InsertMap("sqlo_test", map[string]any{"dedef": "xxx"})
	if err != nil {
		t.Fatalf("InsertMap test dedef error: %v", err)
	}
	err = x.Get(&dedef, "select dedef from sqlo_test")
	if err != nil {
		t.Fatalf("Select dedef after insert test dedef error: %v", err)
	}
	if dedef != "defval" {
		t.Errorf("select dedef should be xxx is: %s", dedef)
	}

	_, err = x.Exec("delete from sqlo_test")
	if err != nil {
		t.Fatalf("delete test error: %v", err)
	}

	err = x.Get(&dedef, "insert into sqlo_test (vl) values ('value') returning dedef")
	if err != nil {
		t.Fatalf("insertget test error: %v", err)
	}

	err = x.InsertMapReturning(&dedef, "dedef", "sqlo_test", map[string]any{"vl": "value"})
	if err != nil {
		t.Fatalf("insertmap test error: %v", err)
	}
	if dedef != "defval" {
		t.Errorf("insertmapreturning should return dedef is: %s", dedef)
	}

}

func Test_update(t *testing.T) {
	dbTest := os.Getenv("SQLO_DBTEST")
	if dbTest == "" {
		return
	}
	db, err := sqlx.Open("postgres", dbTest)
	if err != nil {
		t.Fatalf("Open dbtest: %v", err)
	}
	x := New(db)
	one := 0
	x.Get(&one, "select 1")
	if one != 1 {
		t.Errorf("one should be 1 is %d", one)
	}

	_, err = x.Exec("drop table if exists sqlo_test")
	if err != nil {
		t.Fatalf("drop sqlo_test")
	}
	_, err = x.Exec("create table sqlo_test (dedef text default 'defval', vl text)")
	if err != nil {
		t.Errorf("create table sqlo_test : %v", err)
	}
	dedef := ""
	_, err = x.InsertMap("sqlo_test", map[string]any{"vl": "xxx"}) // vl=xxx dedef=defval
	if err != nil {
		t.Fatalf("InsertMap sqlo_test dedef error: %v", err)
	}

	_, err = x.UpdateMap("sqlo_test", map[string]any{"vl": "vvv"}, "dedef=$1", "defval") // vl=vvv dedef=defval
	if err != nil {
		t.Fatalf("UpdateMap sqlo_test dedef error: %v", err)
	}

	err = x.Get(&dedef, "select dedef from sqlo_test where vl='vvv'")
	if err != nil {
		t.Fatalf("Get sqlo_test dedef error: %v", err)
	}
	if dedef != "defval" {
		t.Errorf("dedef should be ddd is %s", dedef)
	}

	err = x.UpdateMapReturning(&dedef, "dedef", "sqlo_test", map[string]any{"vl": "vvv", "dedef": "ddd"}, "1=1")
	if err != nil {
		t.Fatalf("UpdateMap sqlo_test dedef error: %v", err)
	}
	if dedef != "ddd" {
		t.Errorf("dedef should be ddd is %s", dedef)
	}

	err = x.UpdateMapReturning(&dedef, "dedef", "sqlo_test", map[string]any{"vl": "vvv", "dedef": Raw("DEFAULT")}, "1=1")
	if err != nil {
		t.Fatalf("UpdateMapReturing sqlo_test dedef error: %v", err)
	}
	if dedef != "defval" {
		t.Errorf("dedef should be ddd is %s", dedef)
	}

}
