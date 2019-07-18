package main

import (
  "github.com/hashicorp/terraform/helper/schema"
  "log"
  "time"
  "fmt"
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
  var resource_params map[string]string = conf.Orderresourceparameters
  id := orderFromCatalog(resource_params)
  d.SetId(id)
  log.Printf("Id (%v) of new resourceVM set", id)
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
    log.Printf("Failed to read VM specs, removing %v",d.Id())
    d.SetId("")
    return nil
  }
  name := resp["name"].(string)
  d.Set("name", name)
  log.Printf("VM (%v) name set to %v", d.Id(), name)
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
  if err != nil {
    log.Printf("Failed deleting: %T",err)
  }
  return err
}


func orderFromCatalog(resource_params map[string]string) string {
  path := "/service_catalogs?expand=resources"
  resp, err := apicall(path,"",nil)
  if err != nil {
    log.Printf("Failed to get service_catalogs: %T",err)
    panic(err)
  }
  log.Printf("Type resources: %T", resp["resources"])
  resources := resp["resources"].([]interface{})
  log.Printf("Type resource: %T", resources[0])
  resource := resources[0].(map[string]interface{})
  log.Printf("Type catalog_href: %T", resource["href"])
  catalog_href := resource["href"].(string)
  log.Printf("Type template: %T", resource["service_templates"])
  template := resource["service_templates"].(map[string]interface{})
  log.Printf("Type template_resources: %T", template["resources"])
  template_resources := template["resources"].([]interface{})
  log.Printf("Type template_resource: %T", template_resources[0])
  template_resource := template_resources[0].(map[string]interface{})
  log.Printf("Type service_href: %T", template_resource["href"])
  service_href := template_resource["href"].(string)
  
  resource_params["href"] = service_href
  body := map[string]interface{}{ "action": "order", "resource": resource_params }
  order_href := catalog_href + "/service_templates"
  resp2, err := apicall(order_href, "POST", body)
  if err != nil {
    log.Printf("Failed to orderFromCatalog: %T",err)
    panic(err)
  }
  
  log.Printf("Type results: %T", resp2["results"])
  results := resp2["results"]
  log.Printf("Type resultlist: %T", results)
  resultlist := results.([]interface{})
  log.Printf("Type result: %T", resultlist[0])
  result := resultlist[0].(map[string]interface{})
  log.Printf("Type service_req_href: %T", result["href"])
  service_req_href := result["href"].(string)
  
  path = service_req_href + "?expand=request_tasks"
  var id string
  // we'll loop for half an hour, increasing the timeout
  // in javascript: var j=0;for(var i=0;i<61;i++){j+=i;console.log(j)}
  for i := 0; i < 61; i++ {
    time.Sleep(time.Duration(i) * time.Second)
    resp3, err := apicall(path,"",nil)
    if err != nil {
      panic(err)
    }
    // if resp3["request_state"] == "finished"
    if resp3["request_tasks"] != nil {
      log.Printf("Type request_tasks: %T", resp3["request_tasks"])
      request_tasks := resp3["request_tasks"].([]interface{})
      for _,val := range request_tasks {
        log.Printf("Type val: %T",val)
        task := val.(map[string]interface{})
        log.Printf("Type dt: %T",task["destination_type"])
        dt := task["destination_type"]
        if dt != nil {
          if dt.(string) == "Vm" {
            if task["destination_id"] != nil {
              id = task["destination_id"].(string)
              log.Printf("Found id: %v",id)
              i = 999999 // exits the loop
            }
          }
        } else {
          log.Printf("No destination_type found in on of the request_tasks")
        }
      }
    }
  }
  if id == "" {
    msg := "Got no new ID, please look at this manually"
    log.Printf(msg)
    fmt.Errorf(msg)
  }
  return id
}

