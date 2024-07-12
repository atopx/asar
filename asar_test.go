package asar_test

import (
	"testing"

	"github.com/atopx/asar"
)


func TestUnpack(t *testing.T) {
	if err := asar.Unpack("node_modules.asar", "node"); err != nil {
		t.Fatal(err)
	}
}

func TestPack(t *testing.T) {
	if err := asar.Pack("node", "node.asar"); err != nil {
		t.Fatal(err)
	}
}