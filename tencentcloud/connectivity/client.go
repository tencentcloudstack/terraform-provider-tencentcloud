package connectivity

import (
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	redis "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412"
)

//client for all TencentCloud service
type TencentCloudClient struct {
	Region    string
	SecretId  string
	SecretKey string
	mysqlConn *cdb.Client
	redisConn *redis.Client
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
