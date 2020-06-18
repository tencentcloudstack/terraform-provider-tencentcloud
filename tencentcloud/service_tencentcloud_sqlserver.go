package tencentcloud

import (
	"context"
	"log"

	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type SqlserverService struct {
	client *connectivity.TencentCloudClient
}

func (me *SqlserverService) DescribeZones(ctx context.Context) (zoneInfoList []*sqlserver.ZoneInfo, errRet error) {
	logId := getLogId(ctx)
	request := sqlserver.NewDescribeZonesRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseSqlserverClient().DescribeZones(request)
	if err != nil {
		errRet = err
		return
	}
	if response != nil && response.Response != nil {
		zoneInfoList = response.Response.ZoneSet
	}
	return
}

func (me *SqlserverService) DescribeProductConfig(ctx context.Context, zone string) (specInfoList []*sqlserver.SpecInfo, errRet error) {
	logId := getLogId(ctx)
	request := sqlserver.NewDescribeProductConfigRequest()
	request.Zone = &zone

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseSqlserverClient().DescribeProductConfig(request)
	if err != nil {
		errRet = err
		return
	}
	if response != nil && response.Response != nil {
		specInfoList = response.Response.SpecInfoList
	}
	return
}
