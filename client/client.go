package client

// inspired by https://github.com/spaceapegames/terraform-provider-example/blob/master/api/client/client.go

import (
  "bytes"
  "log"
  "time"
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
	httpClient *http.Client
}

func NewClient(hostname string, username string, password string, insecure bool) *Client {
  client := &http.Client{}
  if insecure {
    tr := &http.Transport{ TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, }
    client = &http.Client{Transport: tr}
  }
  if hostname == "" {
    panic("No hostname specified")
  }
  if username == "" {
    log.Printf("WARNING no username found, will not use BASIC-AUTH")
  }
  if password == "" {
    log.Printf("WARNING no password found, will not use BASIC-AUTH")
  }
	return &Client{
		hostname: hostname,
		username: username,
		password: password,
		httpClient: client,
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
  return c.doCall(path,method,body,5)
}
func (c *Client) retryCall(path string, method string, body interface{}, retry int, errbody string, err error) (map[string]interface{}, error) {
  if retry < 1 {
    if errbody == "" {
      panic(err)
    }
    panic(errbody)
  }
  log.Printf("Waiting 30sec in retryCall %v %v",method,path)
  time.Sleep(time.Duration(30) * time.Second)
  return c.doCall(path,method,body,retry - 1)
}
func (c *Client) doCall(path string, method string, body interface{}, retry int) (map[string]interface{}, error) {
  href := c.getHref(path)
  if method == "" {
    method = "GET"
  }
  
  client := c.httpClient
  
  log.Printf("User %v will do an API call: %v %v",c.username,method,href)
  jsonbody, err := json.Marshal(body)
  reqbody := bytes.NewBuffer(jsonbody)
  req, err := http.NewRequest(method, href, reqbody)
  if err != nil {
    log.Printf("Failed creating NewRequest: %T",err)
    panic(err)
  }
  if c.username != "" {
    if c.password != "" {
      req.SetBasicAuth(c.username,c.password)
    }
  }
  resp, err := client.Do(req)
  if err != nil {
    log.Printf("Failed doing request: %T",err)
    return c.retryCall(path,method,body,retry,"",err)
  }
  defer resp.Body.Close()
  respbody, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    log.Printf("Failed reading body: %T", err)
    return c.retryCall(path,method,body,retry,"",err)
  }
  log.Printf("Response body: %v",string(respbody))

  if resp.StatusCode >= 300 {
    log.Printf("Request did not return 2XX response: %v",resp.StatusCode)
    return c.retryCall(path,method,body,retry,string(respbody),err)
  }
  
  var result map[string]interface{}
  //json.NewDecoder(resp.Body).Decode(&result)
  json.Unmarshal(respbody,&result)

  log.Printf("Completed API call %v %v",method,path)
  return result, err
}


