package ags

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	ags "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ags/v20250920"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func NewAgsService(client *connectivity.TencentCloudClient) AgsService {
	return AgsService{client: client}
}

type AgsService struct {
	client *connectivity.TencentCloudClient
}

func (me *AgsService) DescribeSandboxToolById(ctx context.Context, toolId string) (tool *ags.SandboxTool, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := ags.NewDescribeSandboxToolListRequest()
	request.ToolIds = []*string{&toolId}
	request.Limit = helper.Int64(1)

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseAgsClient().DescribeSandboxToolListWithContext(ctx, request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil || len(response.Response.SandboxToolSet) == 0 {
		return
	}

	tool = response.Response.SandboxToolSet[0]
	return
}
