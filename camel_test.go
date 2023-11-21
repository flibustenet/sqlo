// Copyright (c) 2023 William Dode
// Licensed under the MIT license. See LICENSE file in the project root for details.

package sqlo

import (
	"log"
	"testing"
)

func Test_camel(t *testing.T) {
	type S struct {
		a string
		b string
	}
	tbl := []S{
		S{"Camel", "camel"},
		S{"CamelSnake", "camel_snake"},
	}
	for _, t := range tbl {
		opt := CamelToSnake(t.a)
		if opt != t.b {
			log.Fatalf("camel %s want %s got %s", t.a, t.b, opt)
		}
	}
}
