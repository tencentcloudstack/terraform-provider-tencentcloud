package tencentcloud

import (
	"context"
	"log"

	audit "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cloudaudit/v20190319"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type AuditService struct {
	client *connectivity.TencentCloudClient
}

func (me *AuditService) DescribeAuditCosRegions(ctx context.Context) (regions []*audit.CosRegionInfo, errRet error) {
	logId := getLogId(ctx)
	request := audit.NewListCosEnableRegionRequest()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseAuditClient().ListCosEnableRegion(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	regions = response.Response.EnableRegions
	return
}

func (me *AuditService) DescribeAuditCmqRegions(ctx context.Context) (regions []*audit.CmqRegionInfo, errRet error) {
	logId := getLogId(ctx)
	request := audit.NewListCmqEnableRegionRequest()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseAuditClient().ListCmqEnableRegion(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	regions = response.Response.EnableRegions
	return
}
