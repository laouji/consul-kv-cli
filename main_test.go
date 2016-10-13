package main

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/mitchellh/cli"
)

func setupConsulMock(t *testing.T) {
	agentSelfResponse := &consulAgentSelf{
		Member: struct {
			Name string `json:"Name"`
		}{
			Name: "testNode",
		},
	}

	httpmock.RegisterResponder("PUT",
		fmt.Sprintf("http://127.0.0.1:8500/v1/kv/%s/echo", agentSelfResponse.Member.Name),
		httpmock.NewStringResponder(200, "true"),
	)

	httpmock.RegisterResponder("GET",
		"http://127.0.0.1:8500/v1/agent/self",
		func(req *http.Request) (*http.Response, error) {
			res, err := httpmock.NewJsonResponse(200, agentSelfResponse)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return res, nil
		},
	)

	httpmock.Activate()
}

func TestPutCommand(t *testing.T) {
	setupConsulMock(t)
	defer httpmock.DeactivateAndReset()

	uiMock := new(cli.MockUi)
	pc := &putCommand{UI: uiMock}

	if len(pc.Synopsis()) == 0 {
		t.Fatalf("synopsis is empty for command: %s", "put")
	}

	if len(pc.Help()) == 0 {
		t.Fatalf("help is empty for command: %s", "put")
	}

	args := []string{"echo", "hello"}
	exitStatus := pc.Run(args)
	if exitStatus != 0 {
		t.Fatalf("putCommand.Run() exited with an abnormal status: %d, %s", exitStatus, uiMock.ErrorWriter.String())
	}
}

func TestDeleteCommand(t *testing.T) {
	setupConsulMock(t)
	defer httpmock.DeactivateAndReset()

	uiMock := new(cli.MockUi)
	dc := &deleteCommand{UI: uiMock}

	if len(dc.Synopsis()) == 0 {
		t.Fatalf("synopsis is empty for command: %s", "delete")
	}

	if len(dc.Help()) == 0 {
		t.Fatalf("help is empty for command: %s", "delete")
	}
}
