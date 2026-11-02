package ags

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	agsv20250920 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ags/v20250920"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type AgsService struct {
	client *connectivity.TencentCloudClient
}

func (me *AgsService) DescribeSandboxToolById(ctx context.Context, toolId string) (ret *agsv20250920.SandboxTool, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := agsv20250920.NewDescribeSandboxToolListRequest()
	response := agsv20250920.NewDescribeSandboxToolListResponse()
	request.ToolIds = []*string{&toolId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var (
		offset int64 = 0
		limit  int64 = 100
	)

	for {
		request.Offset = &offset
		request.Limit = &limit
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			result, e := me.client.UseAgsV20250920Client().DescribeSandboxToolList(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Describe sandbox_tool list failed, Response is nil."))
			}

			response = result
			return nil
		})

		if err != nil {
			errRet = err
			return
		}

		if response.Response.SandboxToolSet != nil && len(response.Response.SandboxToolSet) > 0 {
			ret = response.Response.SandboxToolSet[0]
			return
		}

		if len(response.Response.SandboxToolSet) < int(limit) {
			break
		}

		offset += limit
	}

	return
}
