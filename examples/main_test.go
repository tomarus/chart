package main

import (
	"testing"
)

func TestSin(t *testing.T) {
	a := mksin(10, .25, 1, -2)
	t.Logf("ent=%#v", a)
}
