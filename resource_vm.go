package main

import (
  "bytes"
  "os"
  "encoding/json"
  "net/http"
// TODO: PathEscape and QueryEscape, for security
//  "net/url"
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
  var resource_params map[string]string // TODO Here you should put your vm values, like vm_memory and tags such as billing/cost center
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

/* Example code:
  client := m.(*MyClient)

  // Attempt to read from an upstream API
  obj, ok := client.Get(d.Id())

  // If the resource does not exist, inform Terraform. We want to immediately
  // return here to prevent further processing.
  if !ok {
    d.SetId("")
    return nil
  }

  d.Set("address", obj.Address)
  return nil
*/
}

func resourceVMDelete(d *schema.ResourceData, m interface{}) error {
/*
https://github.com/ManageIQ/manageiq_docs/blob/master/api/examples/delete_vm.adoc
https://github.com/ManageIQ/manageiq_docs/blob/master/api/reference/vms.adoc#delete-vm
*/
  return nil
}

func apicall(path string, method string, body map[string]interface{} ) (map[string]interface{}, error) {
  if method == "" {
    method = "GET"
  }
  var uri_base string = "https://" + os.Getenv("MANAGEIQ_API_HOSTNAME") + "/api/"
  var username string = os.Getenv("MANAGEIQ_USERNAME")
  var password string = os.Getenv("MANAGEIQ_PASSWORD")

  client := &http.Client{}
  
  var uri string
  // both a relative from /api or a full (get from link) are possible
  if path[0:4] == "http" {
    uri = path
  } else {
    uri = uri_base + path
  }
  jsonbody, err := json.Marshal(body)
  reqbody := bytes.NewBuffer(jsonbody)
  req, err := http.NewRequest(method, uri, reqbody)
  req.SetBasicAuth(username,password)
  resp, err := client.Do(req)
  
  var result map[string]interface{}
  json.NewDecoder(resp.Body).Decode(&result)
  //json.Unmarshal(resp.Body,&result)

  return result, err
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
