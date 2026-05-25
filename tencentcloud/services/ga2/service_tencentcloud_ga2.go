package ga2

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	ga2v20250115 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ga2/v20250115"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func NewGa2Service(client *connectivity.TencentCloudClient) Ga2Service {
	return Ga2Service{client: client}
}

type Ga2Service struct {
	client *connectivity.TencentCloudClient
}

func (me *Ga2Service) DescribeAccelerateAreasByGlobalAcceleratorId(ctx context.Context, globalAcceleratorId string) (ret []*ga2v20250115.AcceleratorAreas, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := ga2v20250115.NewDescribeAccelerateAreasRequest()
	request.GlobalAcceleratorId = helper.String(globalAcceleratorId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var (
		offset uint64 = 0
		limit  uint64 = 100
	)

	for {
		request.Offset = &offset
		request.Limit = &limit

		var response *ga2v20250115.DescribeAccelerateAreasResponse
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			result, e := me.client.UseGa2V20250115Client().DescribeAccelerateAreas(request)
			if e != nil {
				return tccommon.RetryError(e)
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("DescribeAccelerateAreas failed, Response is nil"))
			}

			response = result
			return nil
		})

		if err != nil {
			errRet = err
			return
		}

		if response.Response.AccelerateAreaSet != nil {
			ret = append(ret, response.Response.AccelerateAreaSet...)
		}

		if response.Response.TotalCount == nil || uint64(len(ret)) >= *response.Response.TotalCount {
			break
		}

		offset += limit
	}

	return
}
