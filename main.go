package main // import "github.com/laouji/consul-kv-cli"

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/mitchellh/cli"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

var commands map[string]cli.CommandFactory

func init() {
	ui := &cli.BasicUi{Writer: os.Stdout}
	commands = map[string]cli.CommandFactory{
		"put": func() (cli.Command, error) {
			return &putCommand{UI: ui}, nil
		},
	}
}

func main() {
	c := cli.NewCLI("consul-kv-cli", "1.0.0")
	c.Args = os.Args[1:]
	c.Commands = commands

	exitStatus, err := c.Run()
	if err != nil {
		fmt.Println(err)
	}

	os.Exit(exitStatus)
}

type putCommand struct {
	UI cli.Ui
}

func (c *putCommand) Synopsis() string {
	return "Stash the value of stdout of a command in consul kv store"
}

func (c *putCommand) Help() string {
	helpText := `
Usage: consul-kv-cli put <key-suffix> arg ...

	Executes the subcommand passed via arg (with all following arguments and passes its stdout
	(preserving formatting and linebreaks) to consul's kv store on the local node
`
	return strings.TrimSpace(helpText)
}

func (c *putCommand) Run(args []string) int {
	if len(args) < 2 {
		c.UI.Error("A key suffix and subcommand must be specified")
		c.UI.Error("")
		c.UI.Error(c.Help())
		return 1
	}
	keySuffix := args[0]

	cmd := exec.Command(args[1], args[2:]...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	if err = cmd.Start(); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	errorReader := bufio.NewReader(stderr)
	errorMsg, err := errorReader.ReadString('\n')
	if err != nil && err.Error() != "EOF" {
		c.UI.Error(fmt.Sprintf("can't read from stderr: %v", err))
		return 1
	}

	var crontabContent []byte

	scanner := bufio.NewScanner(stdout)
	scanner.Split(bufio.ScanBytes)
	for scanner.Scan() {
		bytes := scanner.Bytes()
		crontabContent = append(crontabContent, bytes...)
	}

	// Assumes stderr will not be empty if the command fails to complete here
	if err = cmd.Wait(); err != nil {
		c.UI.Error(fmt.Sprintf("%s: %s", err.Error(), errorMsg))
		return 1
	}

	nodeName, err := nodeName()
	if err != nil {
		c.UI.Error(fmt.Sprintf("couldn't get Node name: %v", err))
		return 1
	}

	err = setKey(fmt.Sprintf("%s/%s", nodeName, keySuffix), crontabContent)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	return 0
}

func setKey(keyName string, bytes []byte) error {
	size := len(bytes)

	// kv max size is 512kB: https://www.consul.io/docs/agent/http/kv.html
	if size > (512 * 1024) {
		fmt.Printf("crontab is too large. contents cannot exceed 512kB. current size is %dkB\n", size/1024)
		os.Exit(1)
	}

	keyValue := string(bytes[:size])

	client := &http.Client{}
	req, err := http.NewRequest(
		"PUT",
		fmt.Sprintf("http://localhost:8500/v1/kv/%s", keyName),
		strings.NewReader(keyValue),
	)
	if err != nil {
		return err
	}

	_, err = client.Do(req)
	if err != nil {
		return err
	}

	return nil
}

type consulAgentSelf struct {
	Member struct {
		Name string `json:"Name"`
	} `json:"Member"`
}

func nodeName() (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://localhost:8500/v1/agent/self", nil)
	if err != nil {
		return "", err
	}

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var self consulAgentSelf
	err = json.Unmarshal(body, &self)
	if err != nil {
		return "", err
	}

	return self.Member.Name, nil
}
