// Copyright (c) 2023 William Dode
// Licensed under the MIT license. See LICENSE file in the project root for details.

package sqlo

import (
	"fmt"
	"strings"
)

type Where struct {
	Style string // $ ? @p
	where []string
	Args  []interface{}
}

// remplace %s par ? ou $d
func (w *Where) And(s string, a ...interface{}) {
	if len(s) > 0 {
		arg_nb := []interface{}{} // deviendra $1 $2 $...
		for i := 0; i < len(a); i++ {
			switch w.Style {
			case "?":
				arg_nb = append(arg_nb, "?")
			case "@p":
				arg_nb = append(arg_nb, fmt.Sprintf("@p%d", len(w.Args)+1+i))
			case "", "$":
				arg_nb = append(arg_nb, fmt.Sprintf("$%d", len(w.Args)+1+i))
			}
		}
		w.where = append(w.where, fmt.Sprintf(s, arg_nb...))
	}

	if len(a) > 0 {
		w.Args = append(w.Args, a...)
	}
}

// ajout sous forme de liste
// exemple AppendListf("xyz in (%s)", "a","b","c")
// doit ajouter "xyz in ($1,$2,$3)" avec args "a","b","c"
func (w *Where) AndList(s string, a ...interface{}) {
	q := []string{} // les $1 $2...
	for i := 0; i < len(a); i++ {
		w.Args = append(w.Args, a[i])
		switch w.Style {
		case "?":
			q = append(q, "?")
		case "@p":
			q = append(q, fmt.Sprintf("@p%d", len(w.Args)))
		case "", "$":
			q = append(q, fmt.Sprintf("$%d", len(w.Args)))
		}
	}
	w.where = append(w.where, fmt.Sprintf(s, strings.Join(q, ",")))
}

// renvoi le where avec "where" sauf si vide
func (w *Where) Where() string {
	if len(w.where) == 0 {
		return ""
	}
	ands := strings.Join(w.where, " and ")
	return " where " + ands
}

// renvoi une copie
func (w *Where) Clone() *Where {
	wc := &Where{}
	wc.Args = append(wc.Args, w.Args...)
	wc.where = append(wc.where, w.where...)
	return wc
}
