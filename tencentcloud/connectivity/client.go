package connectivity

import (
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	as "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/as/v20180419"
	cbs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cbs/v20170312"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	dc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dc/v20180410"
	gaap "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/gaap/v20180529"
	mongodb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mongodb/v20180408"
	redis "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412"
	tag "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tag/v20180813"
	tke "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
)

// client for all TencentCloud service
type TencentCloudClient struct {
	Region      string
	SecretId    string
	SecretKey   string
	mysqlConn   *cdb.Client
	cosConn     *s3.S3
	redisConn   *redis.Client
	asConn      *as.Client
	vpcConn     *vpc.Client
	cbsConn     *cbs.Client
	cvmConn     *cvm.Client
	clbConn     *clb.Client
	dcConn      *dc.Client
	tagConn     *tag.Client
	mongodbConn *mongodb.Client
	tkeConn     *tke.Client
	gaapCoon    *gaap.Client
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
	// all request use method POST
	cpf.HttpProfile.ReqMethod = "POST"
	// request timeout
	cpf.HttpProfile.ReqTimeout = 300
	// cpf.SignMethod = "HmacSHA1"
	cpf.Language = "en-US"

	var round LogRoundTripper

	mysqlClient, _ := cdb.NewClient(credential, me.Region, cpf)
	me.mysqlConn = mysqlClient
	me.mysqlConn.WithHttpTransport(&round)

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
	cpf.Language = "en-US"

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
	cpf.Language = "en-US"

	asConn, _ := as.NewClient(credential, me.Region, cpf)
	me.asConn = asConn

	return me.asConn
}

// get vpc client for service
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
	cpf.Language = "en-US"

	vpcConn, _ := vpc.NewClient(credential, me.Region, cpf)

	var round LogRoundTripper

	vpcConn.WithHttpTransport(&round)

	me.vpcConn = vpcConn

	return me.vpcConn
}

func (me *TencentCloudClient) UseCbsClient() *cbs.Client {
	if me.cbsConn != nil {
		return me.cbsConn
	}
	credential := common.NewCredential(
		me.SecretId,
		me.SecretKey,
	)

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.ReqMethod = "POST"
	cpf.HttpProfile.ReqTimeout = 300
	cpf.Language = "en-US"

	cbsConn, _ := cbs.NewClient(credential, me.Region, cpf)
	me.cbsConn = cbsConn

	return me.cbsConn
}

func (me *TencentCloudClient) UseDcClient() *dc.Client {
	if me.dcConn != nil {
		return me.dcConn
	}

	credential := common.NewCredential(
		me.SecretId,
		me.SecretKey,
	)

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.ReqMethod = "POST"
	cpf.HttpProfile.ReqTimeout = 300
	cpf.Language = "en-US"

	dcConn, _ := dc.NewClient(credential, me.Region, cpf)

	var round LogRoundTripper

	dcConn.WithHttpTransport(&round)

	me.dcConn = dcConn

	return me.dcConn

}

func (me *TencentCloudClient) UseMongodbClient() *mongodb.Client {
	if me.mongodbConn != nil {
		return me.mongodbConn
	}

	credential := common.NewCredential(
		me.SecretId,
		me.SecretKey,
	)

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.ReqMethod = "POST"
	cpf.HttpProfile.ReqTimeout = 300
	cpf.Language = "en-US"

	mongodbConn, _ := mongodb.NewClient(credential, me.Region, cpf)
	var round LogRoundTripper
	mongodbConn.WithHttpTransport(&round)
	me.mongodbConn = mongodbConn

	return me.mongodbConn
}

func (me *TencentCloudClient) UseClbClient() *clb.Client {
	if me.clbConn != nil {
		return me.clbConn
	}

	credential := common.NewCredential(
		me.SecretId,
		me.SecretKey,
	)

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.ReqMethod = "POST"
	cpf.HttpProfile.ReqTimeout = 300
	cpf.Language = "en-US"

	clbConn, _ := clb.NewClient(credential, me.Region, cpf)
	var round LogRoundTripper

	clbConn.WithHttpTransport(&round)
	me.clbConn = clbConn

	return me.clbConn
}

func (me *TencentCloudClient) UseCvmClient() *cvm.Client {
	if me.cvmConn != nil {
		return me.cvmConn
	}

	credential := common.NewCredential(
		me.SecretId,
		me.SecretKey,
	)

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.ReqMethod = "POST"
	cpf.HttpProfile.ReqTimeout = 300
	cpf.Language = "en-US"

	cvmConn, _ := cvm.NewClient(credential, me.Region, cpf)
	var round LogRoundTripper
	cvmConn.WithHttpTransport(&round)
	me.cvmConn = cvmConn
	return me.cvmConn
}

func (me *TencentCloudClient) UseTagClient() *tag.Client {
	if me.tagConn != nil {
		return me.tagConn
	}
	credential := common.NewCredential(
		me.SecretId,
		me.SecretKey,
	)

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.ReqMethod = "POST"
	cpf.HttpProfile.ReqTimeout = 300
	cpf.Language = "en-US"

	tagConn, _ := tag.NewClient(credential, me.Region, cpf)
	var round LogRoundTripper
	tagConn.WithHttpTransport(&round)
	me.tagConn = tagConn
	return me.tagConn
}

func (me *TencentCloudClient) UseTkeClient() *tke.Client {
	if me.tkeConn != nil {
		return me.tkeConn
	}

	credential := common.NewCredential(
		me.SecretId,
		me.SecretKey,
	)

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.ReqMethod = "POST"
	cpf.HttpProfile.ReqTimeout = 300
	cpf.Language = "en-US"

	tkeConn, _ := tke.NewClient(credential, me.Region, cpf)
	var round LogRoundTripper

	tkeConn.WithHttpTransport(&round)
	me.tkeConn = tkeConn

	return me.tkeConn
}

func (me *TencentCloudClient) UseGaapClient() *gaap.Client {
	if me.gaapCoon != nil {
		return me.gaapCoon
	}

	credential := common.NewCredential(
		me.SecretId,
		me.SecretKey,
	)

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.ReqMethod = http.MethodPost
	cpf.HttpProfile.ReqTimeout = 300
	cpf.Language = "en-US"

	gaapConn, _ := gaap.NewClient(credential, me.Region, cpf)
	var round LogRoundTripper
	gaapConn.WithHttpTransport(&round)

	me.gaapCoon = gaapConn

	return me.gaapCoon
}
