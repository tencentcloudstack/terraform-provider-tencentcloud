package tencentcloud

import (
	"context"
	"log"

	redis "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/connectivity"
)

type RedisService struct {
	client *connectivity.TencentCloudClient
}

func (me *RedisService) DescribeRedisZoneConfig(ctx context.Context) (sellConfigures []*redis.RegionConf, errRet error) {

	logId := GetLogId(ctx)
	request := redis.NewDescribeProductInfoRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	response, err := me.client.UseRedisClient().DescribeProductInfo(request)
	if err != nil {
		errRet = err
		return
	}
	sellConfigures = response.Response.RegionSet
	return
}
