package main

type configfile struct {
  Apihostname string `yaml:"api_hostname"`
  Orderresourceparameters map[string]string `yaml:"order_resource_parameters"`
}
