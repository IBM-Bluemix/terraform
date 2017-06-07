package main

import (
	"log"

	"github.com/hashicorp/terraform/builtin/providers/ibm"
	"github.com/hashicorp/terraform/plugin"
)

//Version of this provider plugin
//go build -ldflags "-X main.Version=<version>" github.com/hashicorp/terraform/builtin/bins/provider-ibm will populate the Version
var Version = "dev"

func main() {
	log.Println("IBM Cloud Provider version", Version)
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: ibm.Provider,
	})
}
