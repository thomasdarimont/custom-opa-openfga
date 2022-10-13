package main

import (
	"fmt"
	"github.com/thomasdarimont/custom-opa/custom-opa-openfga/builtins"
	"github.com/thomasdarimont/custom-opa/custom-opa-openfga/plugins"
	"os"

	"github.com/open-policy-agent/opa/cmd"
)

func main() {
	builtins.Register()
	plugins.Register()

	if err := cmd.RootCommand.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
