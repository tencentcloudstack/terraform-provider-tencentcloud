package ccs

import (
	"github.com/zqfan/tencentcloud-sdk-go/common"
)

type Client struct {
	common.Client
}

func NewClientWithSecretId(secretId, secretKey, region string) (client *Client, err error) {
	client = &Client{}
	client.Init(region).WithSecretId(secretId, secretKey)
	return
}
