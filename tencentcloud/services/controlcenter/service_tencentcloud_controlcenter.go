package controlcenter

import (
	"context"
	"fmt"
	"log"

	controlcenterv20230110 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/controlcenter/v20230110"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func NewControlcenterService(client *connectivity.TencentCloudClient) ControlcenterService {
	return ControlcenterService{client: client}
}

type ControlcenterService struct {
	client *connectivity.TencentCloudClient
}

func (me *ControlcenterService) DescribeControlcenterAccountFactoryBaselineConfigById(ctx context.Context) (ret *controlcenterv20230110.GetAccountFactoryBaselineResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := controlcenterv20230110.NewGetAccountFactoryBaselineRequest()
	response := controlcenterv20230110.NewGetAccountFactoryBaselineResponse()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, err := me.client.UseControlcenterV20230110Client().GetAccountFactoryBaseline(request)
		if err != nil {
			return tccommon.RetryError(err)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Get account factory baseline failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	ret = response.Response
	return
}

func (me *ControlcenterService) DescribeControlcenterAccountFactoryBaselineItemsByFilter(ctx context.Context, param map[string]interface{}) (ret []*controlcenterv20230110.AccountFactoryItem, errRet error) {
	var (
		logId    = tccommon.GetLogId(ctx)
		request  = controlcenterv20230110.NewListAccountFactoryBaselineItemsRequest()
		response = controlcenterv20230110.NewListAccountFactoryBaselineItemsResponse()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var (
		offset int64 = 0
		limit  int64 = 20
	)

	for {
		request.Offset = &offset
		request.Limit = &limit

		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			result, err := me.client.UseControlcenterV20230110Client().ListAccountFactoryBaselineItems(request)
			if err != nil {
				return tccommon.RetryError(err)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("List account factory baseline items failed, Response is nil."))
			}

			response = result
			return nil
		})

		if err != nil {
			errRet = err
			return
		}

		if len(response.Response.BaselineItems) < 1 {
			break
		}

		ret = append(ret, response.Response.BaselineItems...)
		if len(response.Response.BaselineItems) < int(limit) {
			break
		}

		offset += limit
	}

	return
}
