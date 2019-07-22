package main

import (
  "github.com/hashicorp/terraform/helper/schema"
)

func Provider() *schema.Provider {
  return &schema.Provider{
    ResourcesMap: map[string]*schema.Resource{
      "manageiq_vm": resourceVM(),
    },
    Schema: map[string]*schema.Schema{
  		"hostname": {
  			Type:     schema.TypeString,
  			Required: true,
  			Description:  "hostname of api endpoint",
  		},
    },
  }
}

