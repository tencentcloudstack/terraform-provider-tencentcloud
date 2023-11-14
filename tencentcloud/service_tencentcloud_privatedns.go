package tencentcloud

import (
	privatedns "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/privatedns/v20201028"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
	"log"
)

type PrivatednsService struct {
	client *connectivity.TencentCloudClient
}

func (me *PrivatednsService) DescribePrivatednsPrivateZoneVpcById(ctx context.Context, zoneId string) (privateZoneVpc *privatedns.PrivateZone, errRet error) {
	logId := getLogId(ctx)

	request := privatedns.NewDescribePrivateZoneRequest()
	request.ZoneId = &zoneId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UsePrivatednsClient().DescribePrivateZone(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.PrivateZone) < 1 {
		return
	}

	privateZoneVpc = response.Response.PrivateZone[0]
	return
}

func (me *PrivatednsService) DeletePrivatednsPrivateZoneVpcById(ctx context.Context, zoneId string) (errRet error) {
	logId := getLogId(ctx)

	request := privatedns.NewDeleteSpecifyPrivateZoneVpcRequest()
	request.ZoneId = &zoneId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UsePrivatednsClient().DeleteSpecifyPrivateZoneVpc(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}
