package client

// inspired by https://github.com/spaceapegames/terraform-provider-example/blob/master/api/client/client.go

import (
  "bytes"
  "log"
  "encoding/json"
  "net/http"
  "io/ioutil"
  "crypto/tls"  
// TODO: PathEscape and QueryEscape, for security
//  "net/url"
)

type Client struct {
	hostname string
	username string
	password string
	insecure bool
	httpClient *http.Client
}

func NewClient(hostname string, username string, password string, insecure bool) *Client {
	return &Client{
		hostname: hostname,
		username: username,
		password: password,
		insecure: insecure,
		httpClient: &http.Client{},
	}
}

func (c *Client) getHref(path string) string {
  var uri_base string = "https://" + c.hostname + "/api/"
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

func (c *Client) Apicall(path string, method string, body interface{} ) (map[string]interface{}, error) {
  href := c.getHref(path)
  if method == "" {
    method = "GET"
  }
  
  client := &http.Client{}
  if c.insecure {
    tr := &http.Transport{ TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, }
    client = &http.Client{Transport: tr}
  }
  
  log.Printf("User %v will do an API call: %v %v",c.username,method,href)
  jsonbody, err := json.Marshal(body)
  reqbody := bytes.NewBuffer(jsonbody)
  req, err := http.NewRequest(method, href, reqbody)
  if err != nil {
    log.Printf("Failed creating NewRequest: %T",err)
    panic(err)
  }
  req.SetBasicAuth(c.username,c.password)
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


