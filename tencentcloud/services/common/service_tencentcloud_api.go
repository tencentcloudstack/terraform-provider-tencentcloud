package common

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	api "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/api/v20201106"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type APIService struct {
	client *connectivity.TencentCloudClient
}

func (me *APIService) DescribeZonesWithProduct(ctx context.Context, product string) (zones []*api.ZoneInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := api.NewDescribeZonesRequest()
	request.Product = common.StringPtr(product)

	ratelimit.Check(request.GetAction())
	// API: https://cloud.tencent.com/document/product/1278/55254
	response, err := me.client.UseApiClient().DescribeZones(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	zones = response.Response.ZoneSet
	return
}
