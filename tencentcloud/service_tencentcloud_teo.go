package tencentcloud

import (
	"context"
	"log"

	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220106"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type TeoService struct {
	client *connectivity.TencentCloudClient
}

func (me *TeoService) DescribeTeoZone(ctx context.Context, zoneId string) (zone *teo.DescribeZoneDetailsResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = teo.NewDescribeZoneDetailsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()
	request.Id = &zoneId

	response, err := me.client.UseTeoClient().DescribeZoneDetails(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	zone = response.Response
	return
}

func (me *TeoService) DeleteTeoZoneById(ctx context.Context, zoneId string) (errRet error) {
	logId := getLogId(ctx)

	request := teo.NewDeleteZoneRequest()
	request.Id = &zoneId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTeoClient().DeleteZone(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}
