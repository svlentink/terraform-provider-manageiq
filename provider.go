package main

import (
  "github.com/hashicorp/terraform/helper/schema"
//  "github.com/hashicorp/terraform/terraform"
  "client"
  "strings"
  "os"
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
    ConfigureFunc: providerConfigure,
  }
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	hostname := d.Get("hostname").(string)
	var username string = os.Getenv("MANAGEIQ_USERNAME")
  var password string = os.Getenv("MANAGEIQ_PASSWORD")
  var insecure bool
  if strings.ToUpper(os.Getenv("MANAGEIQ_INSECURE")) == "TRUE" {
    insecure = true
  }
	return client.NewClient(hostname,username,password,insecure), nil
}
