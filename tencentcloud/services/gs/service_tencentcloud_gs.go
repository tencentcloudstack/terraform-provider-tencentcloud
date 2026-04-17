package gs

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	gsv20191118 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/gs/v20191118"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type GsService struct {
	client *connectivity.TencentCloudClient
}

func (me *GsService) DescribeGsAndroidInstancesByFilter(ctx context.Context, param map[string]interface{}) (ret []*gsv20191118.AndroidInstance, errRet error) {
	var (
		logId    = tccommon.GetLogId(ctx)
		request  = gsv20191118.NewDescribeAndroidInstancesRequest()
		response = gsv20191118.NewDescribeAndroidInstancesResponse()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "AndroidInstanceIds" {
			request.AndroidInstanceIds = v.([]*string)
		}
		if k == "AndroidInstanceRegion" {
			request.AndroidInstanceRegion = v.(*string)
		}
		if k == "AndroidInstanceZone" {
			request.AndroidInstanceZone = v.(*string)
		}
		if k == "LabelSelector" {
			request.LabelSelector = v.([]*gsv20191118.LabelRequirement)
		}
		if k == "Filters" {
			request.Filters = v.([]*gsv20191118.Filter)
		}
	}

	var (
		offset int64 = 0
		limit  int64 = 100
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			result, e := me.client.UseGsV20191118Client().DescribeAndroidInstancesWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("DescribeAndroidInstances failed, Response is nil."))
			}

			response = result
			return nil
		})

		if err != nil {
			errRet = err
			return
		}

		ret = append(ret, response.Response.AndroidInstances...)
		if len(response.Response.AndroidInstances) < int(limit) {
			break
		}

		offset += limit
	}

	return
}
