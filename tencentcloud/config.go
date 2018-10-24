package tencentcloud

import (
	"github.com/zqfan/tencentcloud-sdk-go/client"
	cbs "github.com/zqfan/tencentcloud-sdk-go/services/cbs/unversioned"
	ccs "github.com/zqfan/tencentcloud-sdk-go/services/ccs/unversioned"
	cvm "github.com/zqfan/tencentcloud-sdk-go/services/cvm/v20170312"
	lb "github.com/zqfan/tencentcloud-sdk-go/services/lb/unversioned"
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
	cbsConn    *cbs.Client
	ccsConn    *ccs.Client
	lbConn     *lb.Client
	vpcConn    *vpc.Client
}

func (c *Config) Client() (interface{}, error) {
	var tcClient TencentCloudClient
	userAgent := "TF_TC_1.2.2"
	tcClient.commonConn = client.NewClient(c.SecretId, c.SecretKey, c.Region)
	tcClient.commonConn.Debug = true

	cvmConn, err := cvm.NewClientWithSecretId(c.SecretId, c.SecretKey, c.Region)
	cvmConn.WithUserAgent(userAgent)
	if err != nil {
		return nil, err
	}
	tcClient.cvmConn = cvmConn

	vpcConn, err := vpc.NewClientWithSecretId(c.SecretId, c.SecretKey, c.Region)
	vpcConn.WithUserAgent(userAgent)
	if err != nil {
		return nil, err
	}
	tcClient.vpcConn = vpcConn

	cbsConn, err := cbs.NewClientWithSecretId(c.SecretId, c.SecretKey, c.Region)
	cbsConn.WithUserAgent(userAgent)
	if err != nil {
		return nil, err
	}
	tcClient.cbsConn = cbsConn

	ccsConn, err := ccs.NewClientWithSecretId(c.SecretId, c.SecretKey, c.Region)
	ccsConn.WithUserAgent(userAgent)
	if err != nil {
		return nil, err
	}
	tcClient.ccsConn = ccsConn

	lbConn, err := lb.NewClientWithSecretId(c.SecretId, c.SecretKey, c.Region)
	lbConn.WithUserAgent(userAgent)
	if err != nil {
		return nil, err
	}
	tcClient.lbConn = lbConn

	return &tcClient, nil
}
