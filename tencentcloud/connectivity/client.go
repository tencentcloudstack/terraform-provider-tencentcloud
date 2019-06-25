package connectivity

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	as "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/as/v20180419"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	redis "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
)

//client for all TencentCloud service
type TencentCloudClient struct {
	Region    string
	SecretId  string
	SecretKey string
	mysqlConn *cdb.Client
	cosConn   *s3.S3
	redisConn *redis.Client
	asConn    *as.Client
	vpcConn   *vpc.Client
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
	cpf.HttpProfile.ReqTimeout = 300
	//cpf.SignMethod = "HmacSHA1"

	mysqlClient, _ := cdb.NewClient(credential, me.Region, cpf)
	me.mysqlConn = mysqlClient

	return me.mysqlConn
}

// get cos client for service
func (me *TencentCloudClient) UseCosClient() *s3.S3 {
	if me.cosConn != nil {
		return me.cosConn
	}

	resolver := func(service, region string, optFns ...func(*endpoints.Options)) (endpoints.ResolvedEndpoint, error) {
		if service == endpoints.S3ServiceID {
			return endpoints.ResolvedEndpoint{
				URL:           fmt.Sprintf("http://cos.%s.myqcloud.com", region),
				SigningRegion: region,
			}, nil
		}
		return endpoints.DefaultResolver().EndpointFor(service, region, optFns...)
	}
	creds := credentials.NewStaticCredentials(me.SecretId, me.SecretKey, "")

	sess := session.Must(session.NewSession(&aws.Config{
		Credentials:      creds,
		Region:           aws.String(me.Region),
		EndpointResolver: endpoints.ResolverFunc(resolver),
	}))
	return s3.New(sess)
}

// get redis client for service
func (me *TencentCloudClient) UseRedisClient() *redis.Client {
	if me.redisConn != nil {
		return me.redisConn
	}
	credential := common.NewCredential(
		me.SecretId,
		me.SecretKey,
	)

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.ReqMethod = "POST"
	cpf.HttpProfile.ReqTimeout = 300

	redisConn, _ := redis.NewClient(credential, me.Region, cpf)
	me.redisConn = redisConn

	return me.redisConn
}

func (me *TencentCloudClient) UseAsClient() *as.Client {
	if me.asConn != nil {
		return me.asConn
	}
	credential := common.NewCredential(
		me.SecretId,
		me.SecretKey,
	)

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.ReqMethod = "POST"
	cpf.HttpProfile.ReqTimeout = 300

	asConn, _ := as.NewClient(credential, me.Region, cpf)
	me.asConn = asConn

	return me.asConn
}

//get vpc client for service
func (me *TencentCloudClient) UseVpcClient() *vpc.Client {
	if me.vpcConn != nil {
		return me.vpcConn
	}
	credential := common.NewCredential(
		me.SecretId,
		me.SecretKey,
	)

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.ReqMethod = "POST"
	cpf.HttpProfile.ReqTimeout = 300

	vpcConn, _ := vpc.NewClient(credential, me.Region, cpf)
	me.vpcConn = vpcConn

	return me.vpcConn
}
