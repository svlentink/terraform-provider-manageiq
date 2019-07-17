package main

import (
  "bytes"
  "os"
  "log"
  "strings"
  "encoding/json"
  "net/http"
  "gopkg.in/yaml.v2"
  "io/ioutil"
  "crypto/tls"  
// TODO: PathEscape and QueryEscape, for security
//  "net/url"
)

func loadconfig() configfile {
  fileloc := os.Getenv("MANAGEIQ_CONFIGFILE")
  yamlfile, err := ioutil.ReadFile(fileloc)
  if err != nil {
    log.Printf("Loading config failed, please supply a valid MANAGEIQ_CONFIGFILE: %T",err)
    panic(err)
  }

  var conf configfile
  err2 := yaml.Unmarshal(yamlfile, &conf)
  if err2 != nil {
    log.Printf("Failed parsing MANAGEIQ_CONFIGFILE: %T", err2)
    panic(err2)
  }
  log.Printf("Loaded %T from %v", conf, fileloc)
  return conf
}


func apicall(path string, method string, body interface{} ) (map[string]interface{}, error) {
  if method == "" {
    method = "GET"
  }
  conf := loadconfig()
  var uri_base string = "https://" + conf.Apihostname + "/api/"
  var username string = os.Getenv("MANAGEIQ_USERNAME")
  var password string = os.Getenv("MANAGEIQ_PASSWORD")
  var insecure string = os.Getenv("MANAGEIQ_INSECURE")
  log.Printf("User %v will do an API call: %v %v",username,method,path)

  client := &http.Client{}
  if strings.ToUpper(insecure) == "TRUE" {
    tr := &http.Transport{ TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, }
    client = &http.Client{Transport: tr}
  }
  
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
    log.Printf("Failed creating NewRequest: %T",err)
    panic(err)
  }
  req.SetBasicAuth(username,password)
  resp, err := client.Do(req)
  if err != nil {
    log.Printf("Failed doing request: %T",err)
    panic(err)
  }
  log.Printf("Request body: %v",resp.Body)
  
  var result map[string]interface{}
  json.NewDecoder(resp.Body).Decode(&result)
  //json.Unmarshal(resp.Body,&result)

  log.Printf("Completed API call, returning a %T",result)
  return result, err
}

