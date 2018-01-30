package tencentcloud

import (
	"github.com/zqfan/tencentcloud-sdk-go/client"
	cbs "github.com/zqfan/tencentcloud-sdk-go/services/cbs/unversioned"
	cvm "github.com/zqfan/tencentcloud-sdk-go/services/cvm/v20170312"
	vpc "github.com/zqfan/tencentcloud-sdk-go/services/vpc/unversioned"
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
	vpcConn, err := vpc.NewClientWithSecretId(c.SecretId, c.SecretKey, c.Region)
	if err != nil {
		return nil, err
	}
	tcClient.vpcConn = vpcConn
	cbsConn, err := cbs.NewClientWithSecretId(c.SecretId, c.SecretKey, c.Region)
	if err != nil {
		return nil, err
	}
	tcClient.cbsConn = cbsConn
	return &tcClient, nil
}
