package fwm

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	fwmv20250611 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func NewFwmService(client *connectivity.TencentCloudClient) FwmService {
	return FwmService{client: client}
}

type FwmService struct {
	client *connectivity.TencentCloudClient
}

func (me *FwmService) DescribeFwmRuleGroupById(ctx context.Context, groupId string) (ret *fwmv20250611.RuleGroupResp, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := fwmv20250611.NewDescribeRuleGroupsRequest()
	request.Filters = []*fwmv20250611.CommonFilter{
		{
			Name:         helper.String("GroupId"),
			Values:       []*string{helper.String(groupId)},
			OperatorType: helper.Int64(9),
		},
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ruleGroups, err := me.DescribeFwmRuleGroupsByFilter(ctx, map[string]interface{}{
		"Filters": request.Filters,
	})
	if err != nil {
		errRet = err
		return
	}

	for _, ruleGroup := range ruleGroups {
		if ruleGroup == nil || ruleGroup.GroupId == nil {
			continue
		}
		if *ruleGroup.GroupId == groupId {
			ret = ruleGroup
			return
		}
	}

	return
}

func (me *FwmService) DescribeFwmRuleGroupsByFilter(ctx context.Context, param map[string]interface{}) (ret []*fwmv20250611.RuleGroupResp, errRet error) {
	var (
		logId    = tccommon.GetLogId(ctx)
		request  = fwmv20250611.NewDescribeRuleGroupsRequest()
		response = fwmv20250611.NewDescribeRuleGroupsResponse()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "Product" {
			request.Product = v.(*string)
		}
		if k == "Filters" {
			request.Filters = v.([]*fwmv20250611.CommonFilter)
		}
		if k == "Order" {
			request.Order = v.(*string)
		}
		if k == "By" {
			request.By = v.(*string)
		}
	}

	var (
		offset int64 = 0
		limit  int64 = 100
	)
	for {
		request.Offset = helper.Int64(offset)
		request.Limit = helper.Int64(limit)
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			result, e := me.client.UseFwmV20250611Client().DescribeRuleGroupsWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Describe fwm rule groups failed, Response is nil."))
			}

			response = result
			return nil
		})

		if err != nil {
			errRet = err
			return
		}

		ret = append(ret, response.Response.RuleGroups...)
		if len(response.Response.RuleGroups) < int(limit) {
			break
		}

		offset += limit
	}

	return
}
