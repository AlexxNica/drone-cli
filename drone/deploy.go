package main

import (
	"fmt"
	"html/template"
	"os"
	"strconv"

	"github.com/codegangsta/cli"
	"github.com/drone/drone-go/drone"
)

var deployCmd = cli.Command{
	Name:   "deploy",
	Usage:  "deploy code",
	Action: deploy,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "format",
			Usage: "format output",
			Value: tmplDeployInfo,
		},
	},
}

func deploy(c *cli.Context) error {
	repo := c.Args().First()
	owner, name, err := parseRepo(repo)
	if err != nil {
		return err
	}
	number, err := strconv.Atoi(c.Args().Get(1))
	if err != nil {
		return err
	}

	client, err := newClient(c)
	if err != nil {
		return err
	}

	build, err := client.Build(owner, name, number)
	if err != nil {
		return err
	}
	if build.Event == drone.EventPull {
		return fmt.Errorf("Cannot deploy a pull request")
	}
	env := c.Args().Get(2)
	if env == "" {
		return fmt.Errorf("Please specify the target environment (ie production)")
	}

	deploy, err := client.Deploy(owner, name, number, env)
	if err != nil {
		return err
	}

	tmpl, err := template.New("_").Parse(c.String("format"))
	if err != nil {
		return err
	}
	return tmpl.Execute(os.Stdout, deploy)
}

// template for deployment information
var tmplDeployInfo = `Number: {{ .Number }}
Status: {{ .Status }}
Commit: {{ .Commit }}
Branch: {{ .Branch }}
Ref: {{ .Ref }}
Message: {{ .Message }}
Author: {{ .Author }}
Target: {{ .Deploy }}
`
