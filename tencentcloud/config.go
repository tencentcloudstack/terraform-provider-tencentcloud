package tencentcloud

import (
	"github.com/zqfan/tencentcloud-sdk-go/client"
	"github.com/zqfan/tencentcloud-sdk-go/services/cvm"
)

type Config struct {
	SecretId  string
	SecretKey string
	Region    string
}

type TencentCloudClient struct {
	commonConn *client.Client
	cvmConn    *cvm.Client
}

func (c *Config) Client() (interface{}, error) {
	var tcClient TencentCloudClient
	tcClient.commonConn = client.NewClient(c.SecretId, c.SecretKey, c.Region)
	tcClient.commonConn.Debug = true
	cvmConn, err := cvm.NewClientWithSecretId(c.SecretId, c.SecretKey, c.Region)
	if err != nil {
		return nil, err
	}
	tcClient.cvmConn = cvmConn
	return &tcClient, nil
}
