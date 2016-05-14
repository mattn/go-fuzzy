package fuzzy_test

import (
	"testing"

	"github.com/mattn/go-fuzzy"
)

func TestSimple(t *testing.T) {
	m, s := fuzzy.Match(`ue`, `UnrealEngine`, nil)
	if !m {
		t.Fatal(`Should be true`)
	}
	if s != 11 {
		t.Fatal(`Should be 11`)
	}
}
