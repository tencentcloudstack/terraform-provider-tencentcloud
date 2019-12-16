package connectivity

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	as "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/as/v20180419"
	cam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116"
	cbs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cbs/v20170312"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	cfs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cfs/v20190719"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	dc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dc/v20180410"
	gaap "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/gaap/v20180529"
	mongodb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mongodb/v20180408"
	redis "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412"
	scf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/scf/v20180416"
	tag "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tag/v20180813"
	tcaplusdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcaplusdb/v20190823"
	tke "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	ssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/wss/v20180426"
)

// TencentCloudClient is client for all TencentCloud service
type TencentCloudClient struct {
	Region     string
	Credential *common.Credential

	cosConn     *s3.S3
	mysqlConn   *cdb.Client
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
	camConn     *cam.Client
	gaapConn    *gaap.Client
	sslConn     *ssl.Client
	cfsConn     *cfs.Client
	scfConn     *scf.Client
	tcaplusConn *tcaplusdb.Client
}

// NewTencentCloudClient returns a new TencentCloudClient
func NewTencentCloudClient(secretId, secretKey, securityToken, region string) *TencentCloudClient {
	return &TencentCloudClient{
		Region: region,
		Credential: common.NewTokenCredential(
			secretId,
			secretKey,
			securityToken,
		),
	}
}

// newTencentCloudClientProfile returns a new ClientProfile
func newTencentCloudClientProfile(timeout int) *profile.ClientProfile {
	cpf := profile.NewClientProfile()

	// all request use method POST
	cpf.HttpProfile.ReqMethod = "POST"
	// request timeout
	cpf.HttpProfile.ReqTimeout = timeout
	// default language
	cpf.Language = "en-US"

	return cpf
}

// UseCosClient returns cos client for service
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

	creds := credentials.NewStaticCredentials(me.Credential.SecretId, me.Credential.SecretKey, me.Credential.Token)
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials:      creds,
		Region:           aws.String(me.Region),
		EndpointResolver: endpoints.ResolverFunc(resolver),
	}))

	return s3.New(sess)
}

// UseMysqlClient returns mysql(cdb) client for service
func (me *TencentCloudClient) UseMysqlClient() *cdb.Client {
	if me.mysqlConn != nil {
		return me.mysqlConn
	}

	cpf := newTencentCloudClientProfile(300)
	me.mysqlConn, _ = cdb.NewClient(me.Credential, me.Region, cpf)
	me.mysqlConn.WithHttpTransport(&LogRoundTripper{})

	return me.mysqlConn
}

// UseRedisClient returns redis client for service
func (me *TencentCloudClient) UseRedisClient() *redis.Client {
	if me.redisConn != nil {
		return me.redisConn
	}

	cpf := newTencentCloudClientProfile(300)
	me.redisConn, _ = redis.NewClient(me.Credential, me.Region, cpf)
	me.redisConn.WithHttpTransport(&LogRoundTripper{})

	return me.redisConn
}

// UseAsClient returns as client for service
func (me *TencentCloudClient) UseAsClient() *as.Client {
	if me.asConn != nil {
		return me.asConn
	}

	cpf := newTencentCloudClientProfile(300)
	me.asConn, _ = as.NewClient(me.Credential, me.Region, cpf)
	me.asConn.WithHttpTransport(&LogRoundTripper{})

	return me.asConn
}

// UseVpcClient returns vpc client for service
func (me *TencentCloudClient) UseVpcClient() *vpc.Client {
	if me.vpcConn != nil {
		return me.vpcConn
	}

	cpf := newTencentCloudClientProfile(300)
	me.vpcConn, _ = vpc.NewClient(me.Credential, me.Region, cpf)
	me.vpcConn.WithHttpTransport(&LogRoundTripper{})

	return me.vpcConn
}

// UseCbsClient returns cbs client for service
func (me *TencentCloudClient) UseCbsClient() *cbs.Client {
	if me.cbsConn != nil {
		return me.cbsConn
	}

	cpf := newTencentCloudClientProfile(300)
	me.cbsConn, _ = cbs.NewClient(me.Credential, me.Region, cpf)
	me.cbsConn.WithHttpTransport(&LogRoundTripper{})

	return me.cbsConn
}

// UseDcClient returns dc client for service
func (me *TencentCloudClient) UseDcClient() *dc.Client {
	if me.dcConn != nil {
		return me.dcConn
	}

	cpf := newTencentCloudClientProfile(300)
	me.dcConn, _ = dc.NewClient(me.Credential, me.Region, cpf)
	me.dcConn.WithHttpTransport(&LogRoundTripper{})

	return me.dcConn
}

// UseMongodbClient returns mongodb client for service
func (me *TencentCloudClient) UseMongodbClient() *mongodb.Client {
	if me.mongodbConn != nil {
		return me.mongodbConn
	}

	cpf := newTencentCloudClientProfile(300)
	me.mongodbConn, _ = mongodb.NewClient(me.Credential, me.Region, cpf)
	me.mongodbConn.WithHttpTransport(&LogRoundTripper{})

	return me.mongodbConn
}

// UseClbClient returns clb client for service
func (me *TencentCloudClient) UseClbClient() *clb.Client {
	if me.clbConn != nil {
		return me.clbConn
	}

	cpf := newTencentCloudClientProfile(300)
	me.clbConn, _ = clb.NewClient(me.Credential, me.Region, cpf)
	me.clbConn.WithHttpTransport(&LogRoundTripper{})

	return me.clbConn
}

// UseCvmClient returns cvm client for service
func (me *TencentCloudClient) UseCvmClient() *cvm.Client {
	if me.cvmConn != nil {
		return me.cvmConn
	}

	cpf := newTencentCloudClientProfile(300)
	me.cvmConn, _ = cvm.NewClient(me.Credential, me.Region, cpf)
	me.cvmConn.WithHttpTransport(&LogRoundTripper{})

	return me.cvmConn
}

// UseTagClient returns tag client for service
func (me *TencentCloudClient) UseTagClient() *tag.Client {
	if me.tagConn != nil {
		return me.tagConn
	}

	cpf := newTencentCloudClientProfile(300)
	me.tagConn, _ = tag.NewClient(me.Credential, me.Region, cpf)
	me.tagConn.WithHttpTransport(&LogRoundTripper{})

	return me.tagConn
}

// UseTkeClient returns tke client for service
func (me *TencentCloudClient) UseTkeClient() *tke.Client {
	if me.tkeConn != nil {
		return me.tkeConn
	}

	cpf := newTencentCloudClientProfile(300)
	me.tkeConn, _ = tke.NewClient(me.Credential, me.Region, cpf)
	me.tkeConn.WithHttpTransport(&LogRoundTripper{})

	return me.tkeConn
}

// UseGaapClient returns gaap client for service
func (me *TencentCloudClient) UseGaapClient() *gaap.Client {
	if me.gaapConn != nil {
		return me.gaapConn
	}

	cpf := newTencentCloudClientProfile(300)
	me.gaapConn, _ = gaap.NewClient(me.Credential, me.Region, cpf)
	me.gaapConn.WithHttpTransport(&LogRoundTripper{})

	return me.gaapConn
}

// UseSslClient returns ssl client for service
func (me *TencentCloudClient) UseSslClient() *ssl.Client {
	if me.sslConn != nil {
		return me.sslConn
	}

	cpf := newTencentCloudClientProfile(300)
	me.sslConn, _ = ssl.NewClient(me.Credential, me.Region, cpf)
	me.sslConn.WithHttpTransport(&LogRoundTripper{})

	return me.sslConn
}

// UseCamClient returns cam client for service
func (me *TencentCloudClient) UseCamClient() *cam.Client {
	if me.camConn != nil {
		return me.camConn
	}

	cpf := newTencentCloudClientProfile(300)
	me.camConn, _ = cam.NewClient(me.Credential, me.Region, cpf)
	me.camConn.WithHttpTransport(&LogRoundTripper{})

	return me.camConn
}

// UseCfsClient returns cfs client for service
func (me *TencentCloudClient) UseCfsClient() *cfs.Client {
	if me.cfsConn != nil {
		return me.cfsConn
	}

	cpf := newTencentCloudClientProfile(300)
	me.cfsConn, _ = cfs.NewClient(me.Credential, me.Region, cpf)
	me.cfsConn.WithHttpTransport(&LogRoundTripper{})

	return me.cfsConn
}

// UseScfClient returns scf client for service
func (me *TencentCloudClient) UseScfClient() *scf.Client {
	if me.scfConn != nil {
		return me.scfConn
	}

	cpf := newTencentCloudClientProfile(300)
	me.scfConn, _ = scf.NewClient(me.Credential, me.Region, cpf)
	me.scfConn.WithHttpTransport(&LogRoundTripper{})

	return me.scfConn
}

// UseTcaplusClient returns tcaplush client for service
func (me *TencentCloudClient) UseTcaplusClient() *tcaplusdb.Client {
	if me.tcaplusConn != nil {
		return me.tcaplusConn
	}

	cpf := newTencentCloudClientProfile(300)
	me.tcaplusConn, _ = tcaplusdb.NewClient(me.Credential, me.Region, cpf)
	me.tcaplusConn.WithHttpTransport(&LogRoundTripper{})

	return me.tcaplusConn
}
