package client

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

type Client struct {
	Debug     bool
	secretId  string
	secretKey string
	region    string
	method    string
}

func NewClient(secretId, secretKey, region string) *Client {
	client := &Client{}
	client.secretId = secretId
	client.secretKey = secretKey
	client.region = region
	return client
}

func (c *Client) SendRequest(mod string, params map[string]string) (response string, err error) {
	secretId := c.secretId
	secretKey := c.secretKey
	region := c.region

	method := "POST"
	host := mod + ".api.qcloud.com"
	path := "/v2/index.php"

	paramValues := url.Values{}
	params["SecretId"] = secretId
	if params["Timestamp"] == "" {
		params["Timestamp"] = fmt.Sprintf("%v", time.Now().Unix())
	}
	if params["Nonce"] == "" {
		rand.Seed(time.Now().UnixNano())
		params["Nonce"] = fmt.Sprintf("%v", rand.Int())
	}
	if params["Region"] == "" {
		params["Region"] = region
	}

	sign, err := sign(method, host, path, params, secretKey)
	paramValues.Add("Signature", sign)

	for k, v := range params {
		paramValues.Add(k, v)
	}

	url := "https://" + host + path

	if c.Debug == true {
		log.Printf("[DEBUG] [tencentcloud-sdk-go] request start: action=%v, url=%v, params=%v", params["Action"], url, paramValues)
	}

	rsp, err := http.PostForm(url, paramValues)

	if err != nil {
		log.Fatal("http post error.", err)
		return "", err
	}

	defer rsp.Body.Close()

	buf, err := ioutil.ReadAll(rsp.Body)
	if err != err {
		return "", err
	}

	if c.Debug == true {
		log.Printf("[DEBUG] [tencentcloud-sdk-go] request ended: action=%v, url=%v, params=%v, response=%v",
			params["Action"], url, paramValues, string(buf))
	}

	return string(buf), nil
}

func getSignText(method string, host string, path string, params map[string]string) (text string, err error) {
	method = strings.ToUpper(method)

	text += method + host + path + "?"

	// sort params
	keys := make([]string, 0, len(params))
	for k, _ := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for i := range keys {
		k := keys[i]
		if method == "POST" && params[k][0] == '@' {
			continue
		}
		text += fmt.Sprintf("%v=%v&", strings.Replace(k, "_", ".", -1), params[k])
	}
	text = text[:len(text)-1]
	return text, nil
}

func sign(method string, host string, path string, params map[string]string, secretKey string) (sign string, err error) {
	var source string
	source, err = getSignText(method, host, path, params)
	if err != nil {
		log.Fatalln("Make PlainText error.", err)
		return "", err
	}

	hashed := hmac.New(sha1.New, []byte(secretKey))
	if sm, ok := params["SignatureMethod"]; ok {
		if sm == "HmacSHA256" {
			hashed = hmac.New(sha256.New, []byte(secretKey))
		}
	}
	hashed.Write([]byte(source))

	sign = base64.StdEncoding.EncodeToString(hashed.Sum(nil))
	return sign, nil
}
