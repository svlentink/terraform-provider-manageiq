package main

import (
  "bytes"
  "os"
  "encoding/json"
  "net/http"
  "gopkg.in/yaml.v2"
  "io/ioutil"
// TODO: PathEscape and QueryEscape, for security
//  "net/url"
)

type configfile struct {
  api_hostname string
  order_resource_parameters map[string]string
}

func loadconfig() configfile {
  fileloc := os.Getenv("MANAGEIQ_CONFIGFILE")
  yamlfile, err := ioutil.ReadFile(fileloc)
  if err != nil {
    panic(err)
  }

  var conf configfile
  err2 := yaml.Unmarshal(yamlfile, &conf)
  if err2 != nil {
    panic(err2)
  }
  return conf
}


func apicall(path string, method string, body map[string]interface{} ) (map[string]interface{}, error) {
  if method == "" {
    method = "GET"
  }
  conf := loadconfig()
  var uri_base string = "https://" + conf.api_hostname + "/api/"
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

