package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/terraform"
//	"github.com/terraform-providers/terraform-provider-manageiq/manageiq"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
//		ProviderFunc: manageiq.Provider
    ProviderFunc: func() terraform.ResourceProvider {
      return Provider()
    },
	})
}

