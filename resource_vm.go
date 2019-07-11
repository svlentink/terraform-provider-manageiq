package main

import (
  "github.com/hashicorp/terraform/helper/schema"
)

func resourceVM() *schema.Resource {
  return &schema.Resource{
/*
"The Create, Read, and Delete functions are required for a resource to be functional."
src: https://www.terraform.io/docs/extend/writing-custom-providers.html
*/
    Create: resourceVMCreate,
    Read:   resourceVMRead,
    Delete: resourceVMDelete,

    Schema: map[string]*schema.Schema{
      "name": &schema.Schema{
        Type:     schema.TypeString,
        Computed: true,
      },
/*
      "vm_memory": &schema.Schema{
        Type:     schema.TypeInt,
        Required: true,
      },
      "disk": &schema.Schema{
        Type:     schema.TypeInt,
        Required: true,
      },
      "vlan": &schema.Schema{
        Type:     schema.TypeList,
        Required: true,
      },
*/
    },
  }
}

func resourceVMCreate(d *schema.ResourceData, m interface{}) error {
/*
https://github.com/ManageIQ/manageiq_docs/blob/master/doc-REST_API/topics/provision_request.adoc
https://github.com/ManageIQ/manageiq_docs/blob/master/api/reference/provision_requests.adoc
https://github.com/ManageIQ/manageiq_docs/blob/master/api/examples/provision_request.adoc

https://github.com/ManageIQ/manageiq_docs/blob/master/api/examples/order_service.adoc
*/
  conf := loadconfig()
  var resource_params map[string]string = conf.order_resource_parameters
  resp, err := orderFromCatalog(resource_params)
  if err != nil {
    return err
  }
  results := resp["results"]
  resultlist := results.([]interface{})
  result := resultlist[0].(map[string]interface{})
  id := result["source_id"].(string)
  d.SetId(id)
  return resourceVMRead(d, m)
}

func resourceVMRead(d *schema.ResourceData, m interface{}) error {
/*
http://manageiq.org/docs/reference/fine/api/examples/queries
https://github.com/ManageIQ/manageiq_docs/blob/master/api/examples/provision_request.adoc
*/
  path := "/vms/" + d.Id() //+ "?expand=tags"
  resp, err := apicall(path, "", nil)
  if err != nil {
    d.SetId("")
    return nil
  }
  name := resp["name"].(string)
  d.Set("name", name)
//  d.Set("vm_memory", vm_memory)
//  d.Set("vlan", vlan)
  return err
}

func resourceVMDelete(d *schema.ResourceData, m interface{}) error {
/*
https://github.com/ManageIQ/manageiq_docs/blob/master/api/examples/delete_vm.adoc
https://github.com/ManageIQ/manageiq_docs/blob/master/api/reference/vms.adoc#delete-vm
*/
  body := map[string]string{"action": "delete"}
  path := "/vms/" + d.Id()
  _, err := apicall(path, "DELETE", body)
  return err
}


func orderFromCatalog(resource_params map[string]string) (map[string]interface{}, error) {
  path := "/service_catalogs?expand=resources"
  resp, err := apicall(path,"",nil)
  if err != nil {
    return resp, err
  }
  resources := resp["resources"].([]interface{})
  resource := resources[0].(map[string]interface{})
  catalog_href := resource["href"].(string)
  template := resource["service_templates"].(map[string]interface{})
  template_resources := template["resources"].([]interface{})
  template_resource := template_resources[0].(map[string]string)
  service_href := template_resource["href"]
  
  resource_params["href"] = service_href
  body := map[string]interface{}{ "action": "order", "resource": resource_params }
  order_href := catalog_href + "/service_templates"
  res2, err2 := apicall(order_href, "POST", body)
  return res2, err2
}
