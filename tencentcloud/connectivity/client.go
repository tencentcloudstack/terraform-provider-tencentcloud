package connectivity

import (
	"net/http"

	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	cos "github.com/tencentyun/cos-go-sdk-v5"
)

//client for all TencentCloud service
type TencentCloudClient struct {
	Region    string
	SecretId  string
	SecretKey string
	mysqlConn *cdb.Client
	cosConn   *cos.Client
}

func NewTencentCloudClient(secretId, secretKey, region string) *TencentCloudClient {

	var tencentCloudClient TencentCloudClient

	tencentCloudClient.SecretId,
		tencentCloudClient.SecretKey,
		tencentCloudClient.Region =

		secretId,
		secretKey,
		region

	return &tencentCloudClient
}

// get mysql(cdb) client for service
func (me *TencentCloudClient) UseMysqlClient() *cdb.Client {
	if me.mysqlConn != nil {
		return me.mysqlConn
	}

	credential := common.NewCredential(
		me.SecretId,
		me.SecretKey,
	)

	cpf := profile.NewClientProfile()
	//all request use method POST
	cpf.HttpProfile.ReqMethod = "POST"
	//request timeout
	cpf.HttpProfile.ReqTimeout = 10
	cpf.SignMethod = "HmacSHA1"

	mysqlClient, _ := cdb.NewClient(credential, me.Region, cpf)
	me.mysqlConn = mysqlClient

	return me.mysqlConn
}

// the format of bucketName is {name}-{appid}
// if bucket is null, client can only request GetService
func (me *TencentCloudClient) UseCosClient(bucketName string) (client *cos.Client) {
	if me.cosConn != nil {
		return me.cosConn
	}

	b := &cos.BaseURL{}
	if bucketName != "" {
		b.BucketURL = cos.NewBucketURL(bucketName, me.Region, true)
	}
	client = cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  me.SecretId,
			SecretKey: me.SecretKey,
		},
	})
	return
}
