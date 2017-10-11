package join

import (
	"flag"
	"fmt"

	"github.com/hashicorp/consul/command/flags"
	"github.com/mitchellh/cli"
)

var synopsis = "Tell Consul agent to join cluster"
var usage = `Usage: consul join [options] address ...

  Tells a running Consul agent (with "consul agent") to join the cluster
  by specifying at least one existing member.`

func New(ui cli.Ui) *cmd {
	c := &cmd{UI: ui}
	c.initFlags()
	return c
}

type cmd struct {
	UI     cli.Ui
	flags  *flag.FlagSet
	client *flags.HTTPClient
	wan    bool
}

func (c *cmd) initFlags() {
	c.flags = flag.NewFlagSet("", flag.ContinueOnError)
	c.flags.BoolVar(&c.wan, "wan", false,
		"Joins a server to another server in the WAN pool.")

	c.client = &flags.HTTPClient{}
	flags.Merge(c.flags, c.client.Flags())
}

func (c *cmd) Synopsis() string {
	return synopsis
}

func (c *cmd) Help() string {
	return flags.Usage(usage, c.flags, c.client.Flags(), nil)
}

func (c *cmd) Run(args []string) int {
	if err := c.flags.Parse(args); err != nil {
		return 1
	}

	addrs := c.flags.Args()
	if len(addrs) == 0 {
		c.UI.Error("At least one address to join must be specified.")
		c.UI.Error("")
		c.UI.Error(c.Help())
		return 1
	}

	client, err := flags.NewAPIClient(c.client)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error connecting to Consul agent: %s", err))
		return 1
	}

	joins := 0
	for _, addr := range addrs {
		err := client.Agent().Join(addr, c.wan)
		if err != nil {
			c.UI.Error(fmt.Sprintf("Error joining address '%s': %s", addr, err))
		} else {
			joins++
		}
	}

	if joins == 0 {
		c.UI.Error("Failed to join any nodes.")
		return 1
	}

	c.UI.Output(fmt.Sprintf("Successfully joined cluster by contacting %d nodes.", joins))
	return 0
}
