package main

import (
  "bytes"
  "os"
  "fmt"
  "encoding/json"
  "net/http"
  "gopkg.in/yaml.v2"
  "io/ioutil"
// TODO: PathEscape and QueryEscape, for security
//  "net/url"
)

func loadconfig() configfile {
  fileloc := os.Getenv("MANAGEIQ_CONFIGFILE")
  yamlfile, err := ioutil.ReadFile(fileloc)
  if err != nil {
    fmt.Println("Loading config failed, please supply a valid MANAGEIQ_CONFIGFILE. %T",err)
    panic(err)
  }

  conf := configfile{}
  err2 := yaml.Unmarshal(yamlfile, &conf)
  if err2 != nil {
    fmt.Println("Failed parsing MANAGEIQ_CONFIGFILE. %T", err2)
    panic(err2)
  }
  fmt.Println("Loaded config. %T", conf)
  return conf
}


func apicall(path string, method string, body interface{} ) (map[string]interface{}, error) {
  if method == "" {
    method = "GET"
  }
  conf := loadconfig()
  var uri_base string = "https://" + conf.api_hostname + "/api/"
  var username string = os.Getenv("MANAGEIQ_USERNAME")
  var password string = os.Getenv("MANAGEIQ_PASSWORD")
  fmt.Println("User %v will do an API call. %v %v",username,method,path)

  //tr := &http.Transport{ TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, }
  client := &http.Client{} //Transport: tr}
  
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
  if err != nil {
    fmt.Println("Failed creating NewRequest. %T",err)
    panic(err)
  }
  req.SetBasicAuth(username,password)
  resp, err := client.Do(req)
  if err != nil {
    fmt.Println("Failed doing request. %T",err)
    panic(err)
  }
  
  var result map[string]interface{}
  json.NewDecoder(resp.Body).Decode(&result)
  //json.Unmarshal(resp.Body,&result)

  fmt.Println("Completed API call. %T",result)
  return result, err
}

