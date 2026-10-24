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

// DescribeFwmSecurityGroupRuleById returns the enterprise security group rule identified by
// (groupId, ruleId). It returns (nil, nil) when the rule no longer exists.
func (me *FwmService) DescribeFwmSecurityGroupRuleById(ctx context.Context, groupId, ruleId string) (ret *fwmv20250611.SecGroupRuleResp, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := fwmv20250611.NewDescribeSecurityGroupRulesRequest()
	request.GroupId = helper.String(groupId)
	request.Filters = []*fwmv20250611.CommonFilter{
		{
			Name:         helper.String("RuleId"),
			Values:       []*string{helper.String(ruleId)},
			OperatorType: helper.Int64(9),
		},
	}
	request.Offset = helper.Int64(0)
	request.Limit = helper.Int64(100)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var response *fwmv20250611.DescribeSecurityGroupRulesResponse
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseFwmV20250611Client().DescribeSecurityGroupRulesWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("DescribeSecurityGroupRules failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	for _, item := range response.Response.Rules {
		if item == nil || item.RuleId == nil {
			continue
		}
		if *item.RuleId == ruleId {
			ret = item
			return
		}
	}

	ret = nil
	return
}
