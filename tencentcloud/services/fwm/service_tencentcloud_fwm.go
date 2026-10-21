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

func (me *FwmService) DescribeFwmEdgeAclRuleById(ctx context.Context, groupId, ruleId string, direction int64) (fwmEdgeAclRule *fwmv20250611.EdgeAclRuleResp, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := fwmv20250611.NewDescribeEdgeAclRulesRequest()
	response := fwmv20250611.NewDescribeEdgeAclRulesResponse()

	request.GroupId = helper.String(groupId)
	request.Direction = helper.Int64(direction)
	request.Limit = helper.Int64(1000)
	request.Filters = []*fwmv20250611.CommonFilter{
		{
			Name:         helper.String("RuleId"),
			Values:       []*string{helper.String(ruleId)},
			OperatorType: helper.Int64(1),
		},
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseFwmV20250611Client().DescribeEdgeAclRulesWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe fwm edge acl rule failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	if len(response.Response.Rules) == 0 {
		return nil, nil
	}

	fwmEdgeAclRule = response.Response.Rules[0]
	return
}
