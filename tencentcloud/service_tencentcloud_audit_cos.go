package tencentcloud

import (
	"context"
	"log"

	auditcos "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cloudaudit/v20190319"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type AuditCosService struct {
	client *connectivity.TencentCloudClient
}

func (me *AuditCosService) DescribeRegions(ctx context.Context) (regions []*auditcos.CosRegionInfo, errRet error) {
	logId := getLogId(ctx)
	request := auditcos.NewListCosEnableRegionRequest()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseAuditCosClient().ListCosEnableRegion(request)
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
