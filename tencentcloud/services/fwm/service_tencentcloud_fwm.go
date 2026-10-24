package fwm

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	fwmv20250611 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type FwmService struct {
	client *connectivity.TencentCloudClient
}

func (me *FwmService) DescribeFwmNatAclRuleById(ctx context.Context, groupId, ruleId string, direction int64) (rule *fwmv20250611.NatAclRuleResp, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := fwmv20250611.NewDescribeNatAclRulesRequest()
	request.GroupId = helper.String(groupId)
	request.Direction = helper.Int64(direction)
	limit := int64(100)
	offset := int64(0)
	request.Limit = helper.Int64(limit)
	request.Offset = helper.Int64(offset)

	filter := &fwmv20250611.CommonFilter{
		Name:         helper.String("RuleId"),
		Values:       []*string{helper.String(ruleId)},
		OperatorType: helper.Int64(1),
	}
	request.Filters = []*fwmv20250611.CommonFilter{filter}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseFwmV20250611Client().DescribeNatAclRulesWithContext(ctx, request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil || len(response.Response.Rules) < 1 {
		return
	}

	rule = response.Response.Rules[0]
	return
}
