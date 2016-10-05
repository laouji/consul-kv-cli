package main

import (
	"github.com/mitchellh/cli"
	"testing"
)

func TestPutCommand(t *testing.T) {
	mock := new(cli.MockUi)
	pc := &putCommand{UI: mock}

	if len(pc.Synopsis()) == 0 {
		t.Fatalf("synopsis is empty for command: %s", "put")
	}

	if len(pc.Help()) == 0 {
		t.Fatalf("help is empty for command: %s", "put")
	}
}

func TestDeleteCommand(t *testing.T) {
	mock := new(cli.MockUi)
	dc := &deleteCommand{UI: mock}

	if len(dc.Synopsis()) == 0 {
		t.Fatalf("synopsis is empty for command: %s", "delete")
	}

	if len(dc.Help()) == 0 {
		t.Fatalf("help is empty for command: %s", "delete")
	}
}
