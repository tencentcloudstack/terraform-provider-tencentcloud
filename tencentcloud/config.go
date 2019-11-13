package tencentcloud

import (
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/connectivity"
)

type Config struct {
	SecretId      string
	SecretKey     string
	SecurityToken string
	Region        string
}

type TencentCloudClient struct {
	apiV3Conn *connectivity.TencentCloudClient
}

func (c *Config) Client() (interface{}, error) {
	var tcClient TencentCloudClient
	tcClient.apiV3Conn = connectivity.NewTencentCloudClient(c.SecretId, c.SecretKey, c.SecurityToken, c.Region)

	return &tcClient, nil
}
