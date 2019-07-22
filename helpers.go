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

func loadconfigDEPRECATED() configfile {
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

func getHref(hostname string, path string) string {
  var uri_base string = "https://" + hostname + "/api/"
  if string(path[0]) == "/" {
    // CFME throws an error if you have /api//service_catalogs
    uri_base = uri_base[:len(uri_base)-1]
  }
  var uri string
  // both a relative from /api or a full (get from link) are possible
  if path[0:4] == "http" {
    uri = path
  } else {
    uri = uri_base + path
  }
  return uri
}

func apicall(href string, method string, body interface{} ) (map[string]interface{}, error) {
  if method == "" {
    method = "GET"
  }
  var username string = os.Getenv("MANAGEIQ_USERNAME")
  var password string = os.Getenv("MANAGEIQ_PASSWORD")
  var insecure string = os.Getenv("MANAGEIQ_INSECURE")
  
  client := &http.Client{}
  if strings.ToUpper(insecure) == "TRUE" {
    tr := &http.Transport{ TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, }
    client = &http.Client{Transport: tr}
  }
  
  log.Printf("User %v will do an API call: %v %v",username,method,href)
  jsonbody, err := json.Marshal(body)
  reqbody := bytes.NewBuffer(jsonbody)
  req, err := http.NewRequest(method, href, reqbody)
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
  defer resp.Body.Close()
  respbody, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    log.Printf("Failed reading body: %T", err)
    panic(err)
  }
  log.Printf("Response body: %v",string(respbody))
  
  var result map[string]interface{}
  //json.NewDecoder(resp.Body).Decode(&result)
  json.Unmarshal(respbody,&result)

  log.Printf("Completed API call, returning a %T",result)
  return result, err
}

