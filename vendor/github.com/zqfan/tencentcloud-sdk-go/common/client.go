package common

import (
	"log"
	"net/http"
)

type Client struct {
	region     string
	httpClient *http.Client
	credential Credential
	signMethod string
	userAgent  string
	debug      bool
}

func (c *Client) Send(request Request, response Response) (err error) {
	if request.GetDomain() == "" {
		domain := GetServiceDomain(request.GetService())
		request.SetDomain(domain)
	}
	err = ConstructParams(request)
	if err != nil {
		return
	}
	CompleteCommonParams(request, c)
	err = Sign(request, c.credential, c.signMethod)
	if err != nil {
		return
	}
	httpRequest, err := http.NewRequest(request.GetHttpMethod(), request.GetUrl(), request.GetBodyReader())
	if err != nil {
		return
	}
	if request.GetHttpMethod() == POST {
		httpRequest.Header["Content-Type"] = []string{"application/x-www-form-urlencoded"}
	}
	log.Printf("[DEBUG] http request=%v", httpRequest)
	var httpResponse *http.Response
	httpResponse, err = c.httpClient.Do(httpRequest)
	if err != nil {
		return
	}
	err = ParseFromHttpResponse(httpResponse, response)
	return
}

func (c *Client) GetRegion() string {
	return c.region
}

func (c *Client) GetUserAgent() string {
	return c.userAgent
}

func (c *Client) Init(region string) *Client {
	c.httpClient = &http.Client{}
	c.region = region
	c.signMethod = "HmacSHA256"
	c.debug = true
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	return c
}

func (c *Client) WithSecretId(secretId, secretKey string) *Client {
	c.credential = NewBasicCredential(secretId, secretKey)
	return c
}

func (c *Client) WithSignatureMethod(method string) *Client {
	c.signMethod = method
	return c
}

func (c *Client) WithUserAgent(agent string) *Client {
	c.userAgent = agent
	return c
}

func NewClientWithSecretId(secretId, secretKey, region string) (client *Client, err error) {
	client = &Client{}
	client.Init(region).WithSecretId(secretId, secretKey)
	return
}
