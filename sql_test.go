// Copyright (c) 2023 William Dode

package sqlo

import (
	"log"
	"testing"
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
