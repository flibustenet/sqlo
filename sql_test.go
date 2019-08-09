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
	q, _ := tx.insertSt("mytable", fs)
	if q != "INSERT INTO mytable (ok, yes) VALUES ($1, $2)" {
		log.Fatal(q)
	}
}

func Test_updateSt(t *testing.T) {
	tx := &Sx{}
	fs := map[string]interface{}{}
	fs["ok"] = "coral"
	fs["yes"] = "no"
	q, _ := tx.updateSt("mytable", fs, "ok=$1", "ok")
	if q != "UPDATE mytable SET ok=$2, yes=$3 WHERE ok=$1" {
		log.Fatal(q)
	}
}
