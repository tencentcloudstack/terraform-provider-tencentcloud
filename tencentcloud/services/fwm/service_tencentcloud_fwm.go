package fwm

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
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

func (me *FwmService) DescribeFwmVpcAclRuleById(ctx context.Context, groupId, ruleId string) (ret *fwmv20250611.VpcAclRuleResp, allTotalCount *int64, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := fwmv20250611.NewDescribeVpcAclRulesRequest()
	request.GroupId = helper.String(groupId)
	request.Filters = []*fwmv20250611.CommonFilter{
		{
			Name:         helper.String("RuleId"),
			Values:       []*string{helper.String(ruleId)},
			OperatorType: helper.Int64(1),
		},
	}
	request.Limit = helper.Int64(0)
	request.Offset = helper.Int64(0)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseFwmV20250611Client().DescribeVpcAclRulesWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe fwm vpc_acl_rule failed, Response is nil."))
		}

		if len(result.Response.Rules) == 0 {
			return nil
		}

		ret = result.Response.Rules[0]
		allTotalCount = result.Response.AllTotalCount
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	return
}

func (me *FwmService) DeleteFwmVpcAclRuleById(ctx context.Context, groupId, ruleId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := fwmv20250611.NewDeleteVpcAclRuleRequest()
	request.GroupId = helper.String(groupId)
	request.RuleIds = []*string{helper.String(ruleId)}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseFwmV20250611Client().DeleteVpcAclRuleWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		if sdkErr, ok := err.(*errors.TencentCloudSDKError); ok && sdkErr.GetCode() == "ResourceNotFound" {
			return nil
		}
		errRet = err
	}

	return
}
