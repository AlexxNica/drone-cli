package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/drone/drone-go/drone"

	"github.com/codegangsta/cli"
	"github.com/jackspirou/syscerts"
)

func newClient(c *cli.Context) (drone.Client, error) {
	var token = c.GlobalString("token")
	var server = c.GlobalString("server")

	// if no server url is provided we can default
	// to the hosted Drone service.
	if len(server) == 0 {
		return nil, fmt.Errorf("Error: you must provide the Drone server address.")
	}
	if len(token) == 0 {
		return nil, fmt.Errorf("Error: you must provide your Drone access token.")
	}

	// attempt to find system CA certs
	certs := syscerts.SystemRootsPool()
	tlsConfig := &tls.Config{RootCAs: certs}

	// create the drone client with TLS options
	return drone.NewClientTokenTLS(server, token, tlsConfig), nil
}

func parseRepo(str string) (user, repo string, err error) {
	var parts = strings.Split(str, "/")
	if len(parts) != 2 {
		err = fmt.Errorf("Error: Invalid or missing repository. eg octocat/hello-world.")
		return
	}
	user = parts[0]
	repo = parts[1]
	return
}

func readInput(in string) ([]byte, error) {
	if in == "-" {
		return ioutil.ReadAll(os.Stdin)
	}
	return ioutil.ReadFile(in)
}
