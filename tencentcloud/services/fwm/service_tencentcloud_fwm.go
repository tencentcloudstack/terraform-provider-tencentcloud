package fwm

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	fwmv20250611 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type FwmService struct {
	client *connectivity.TencentCloudClient
}

func NewFwmService(client *connectivity.TencentCloudClient) *FwmService {
	return &FwmService{client: client}
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

func (me *FwmService) DescribeFwmOrganMemberById(ctx context.Context, appId, memberId string) (organMember *fwmv20250611.OrganMemberItem, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := fwmv20250611.NewDescribeOrganMembersRequest()
	response := fwmv20250611.NewDescribeOrganMembersResponse()

	operatorType := int64(1)
	request.Filters = []*fwmv20250611.CommonFilter{
		{
			Name:         helper.String("AppId"),
			Values:       helper.Strings([]string{appId}),
			OperatorType: helper.Int64(operatorType),
		},
	}
	request.Limit = helper.Int64(100)
	request.Offset = helper.Int64(0)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var iacExtInfo connectivity.IacExtInfo
	iacExtInfo.InstanceId = memberId
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseFwmClient(iacExtInfo).DescribeOrganMembersWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil || result.Response.Members == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe fwm organ members failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	if len(response.Response.Members) < 1 {
		return
	}

	for _, member := range response.Response.Members {
		if member.MemberId != nil && *member.MemberId == memberId {
			organMember = member
			return
		}
	}

	return
}

func (me *FwmService) DescribeFwmAccountGroupById(ctx context.Context, groupId string) (accountGroup *fwmv20250611.AccountGroup, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := fwmv20250611.NewDescribeAccountGroupsRequest()
	request.Filters = []*fwmv20250611.CommonFilter{
		{
			Name:         helper.String("GroupId"),
			Values:       []*string{helper.String(groupId)},
			OperatorType: helper.Int64(1),
		},
	}
	request.Limit = helper.Int64(20)
	request.Offset = helper.Int64(0)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	response := fwmv20250611.NewDescribeAccountGroupsResponse()
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseFwmV20250611Client().DescribeAccountGroupsWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe fwm account group failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	if len(response.Response.AccountGroups) < 1 {
		return
	}

	accountGroup = response.Response.AccountGroups[0]
	return
}

func (me *FwmService) DeleteFwmAccountGroupById(ctx context.Context, groupId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := fwmv20250611.NewDeleteAccountGroupRequest()
	request.GroupId = helper.String(groupId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseFwmV20250611Client().DeleteAccountGroupWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete fwm account group failed, Response is nil."))
		}

		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	return
}

func (me *FwmService) DescribeFwmMemberGroup(ctx context.Context, userMemberId string) (member *fwmv20250611.OrganMemberItem, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := fwmv20250611.NewDescribeOrganMembersRequest()

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
		ref := []*fwmv20250611.OrganMemberItem{}
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			result, e := me.client.UseFwmV20250611Client().DescribeOrganMembers(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil || result.Response.Members == nil {
				return resource.NonRetryableError(fmt.Errorf("Describe organ members failed, Response is nil."))
			}

			ref = result.Response.Members
			return nil
		})

		if err != nil {
			errRet = err
			return
		}

		for _, m := range ref {
			if m != nil && m.MemberId != nil && *m.MemberId == userMemberId {
				member = m
				return
			}
		}

		if len(ref) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *FwmService) ModifyFwmMemberGroup(ctx context.Context, userMemberId, groupId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := fwmv20250611.NewModifyMemberGroupRequest()
	request.UserMemberId = helper.String(userMemberId)
	request.GroupId = helper.String(groupId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseFwmV20250611Client().ModifyMemberGroupWithContext(ctx, request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *FwmService) FwmMemberGroupIdToAppIdMemberId(id string) (appId, memberId string, err error) {
	parts := strings.Split(id, tccommon.FILED_SP)
	if len(parts) != 2 {
		err = fmt.Errorf("invalid fwm_member_group id: %s", id)
		return
	}
	appId = parts[0]
	memberId = parts[1]
	return
}

func (me *FwmService) DescribeFwmRuleGroupById(ctx context.Context, groupId string) (ret *fwmv20250611.RuleGroupResp, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := fwmv20250611.NewDescribeRuleGroupsRequest()
	request.Filters = []*fwmv20250611.CommonFilter{
		{
			Name:         helper.String("GroupId"),
			Values:       []*string{helper.String(groupId)},
			OperatorType: helper.Int64(1),
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

func (me *FwmService) DescribeFwmSecurityGroupRuleById(ctx context.Context, groupId, ruleId string) (ret *fwmv20250611.SecGroupRuleResp, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := fwmv20250611.NewDescribeSecurityGroupRulesRequest()
	request.GroupId = helper.String(groupId)
	request.Filters = []*fwmv20250611.CommonFilter{
		{
			Name:         helper.String("RuleId"),
			Values:       []*string{helper.String(ruleId)},
			OperatorType: helper.Int64(1),
		},
	}
	request.Offset = helper.Int64(0)
	request.Limit = helper.Int64(100)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseFwmV20250611Client().DescribeSecurityGroupRulesWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe fwm_security_group_rule failed, Response is nil."))
		}

		if len(result.Response.Rules) == 0 {
			return nil
		}

		ret = result.Response.Rules[0]
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	return
}

func (me *FwmService) DeleteFwmSecurityGroupRuleById(ctx context.Context, groupId, ruleId string) error {
	logId := tccommon.GetLogId(ctx)

	request := fwmv20250611.NewDeleteSecurityGroupRuleRequest()
	request.GroupId = helper.String(groupId)
	request.RuleIds = []*string{helper.String(ruleId)}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseFwmV20250611Client().DeleteSecurityGroupRuleWithContext(ctx, request)
		if e != nil {
			if ee, ok := e.(*sdkErrors.TencentCloudSDKError); ok {
				errCode := ee.GetCode()
				if strings.Contains(errCode, "ResourceNotFound") {
					return nil
				}
			}
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete fwm_security_group_rule failed, Response is nil."))
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s delete fwm_security_group_rule failed, reason:%+v", logId, err)
		return err
	}

	return nil
}

func (me *FwmService) DescribeFwmNatAclRuleById(ctx context.Context, groupId, ruleId string, direction int64) (rule *fwmv20250611.NatAclRuleResp, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := fwmv20250611.NewDescribeNatAclRulesRequest()
	request.GroupId = helper.String(groupId)
	request.Direction = helper.Int64(direction)
	request.Limit = helper.Int64(100)
	request.Offset = helper.Int64(0)
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
		response, e := me.client.UseFwmV20250611Client().DescribeNatAclRulesWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || response.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe fwm_nat_acl_rule failed, Response is nil."))
		}

		if len(response.Response.Rules) == 0 {
			return nil
		}

		rule = response.Response.Rules[0]
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	return
}

func (me *FwmService) DescribeFwmVpcAclRuleById(ctx context.Context, groupId, ruleId string) (ret *fwmv20250611.VpcAclRuleResp, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := fwmv20250611.NewDescribeVpcAclRulesRequest()
	request.GroupId = helper.String(groupId)
	request.Limit = helper.Int64(100)
	request.Offset = helper.Int64(0)
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
			if ee, ok := e.(*sdkErrors.TencentCloudSDKError); ok && ee.GetCode() == "ResourceNotFound" {
				return nil
			}
			return tccommon.RetryError(e)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

		return nil
	})

	if err != nil {
		errRet = err
	}

	return
}

// fwmDispatchStrategyStatusBegin marks "start dispatch". The DispatchStrategy
// API carries an inverse pair in its Status field (1=start, 2=abort), but the
// server only implements the start case; any other value falls through to
// InvalidParameter, so the provider pins it to 1.
const fwmDispatchStrategyStatusBegin int64 = 1

// FwmDispatchStrategyStatus enum from DescribeStrategyDispatchStatus.Response.Status:
// 0=no change / pending, 1=dispatching, 2=succeeded, 3=failed. The provider
// treats 0 and 1 as in-progress and waits, 2 as success, 3 as failure.
const (
	fwmDispatchStatusNoChange  int64 = 0
	fwmDispatchStatusDispatching int64 = 1
	fwmDispatchStatusSucceeded    int64 = 2
	fwmDispatchStatusFailed     int64 = 3
)

// DispatchFwmStrategy triggers a one-time strategy dispatch for the given
// product via the DispatchStrategy API. The cloud-side operator Uin /
// SubAccountUin are injected by the API gateway rather than carried in the
// request body, so they are validated on the Terraform side but not passed
// here. The response carries no business fields (only RequestId); dispatch
// completion is verified separately through DescribeStrategyDispatchStatus.
func (me *FwmService) DispatchFwmStrategy(ctx context.Context, product string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := fwmv20250611.NewDispatchStrategyRequest()
	request.Status = helper.Int64(fwmDispatchStrategyStatusBegin)
	request.Product = helper.String(product)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseFwmV20250611Client().DispatchStrategyWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		return nil
	})

	if err != nil {
		errRet = err
	}

	return
}

// DescribeFwmStrategyDispatchStatus queries the latest dispatch task status for
// the given product (app id / uin are injected by the gateway). The response is
// the latest task of that account x product, so a nil status means no task has
// been started yet.
func (me *FwmService) DescribeFwmStrategyDispatchStatus(ctx context.Context, product string) (status *fwmv20250611.DescribeStrategyDispatchStatusResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := fwmv20250611.NewDescribeStrategyDispatchStatusRequest()
	request.Product = helper.String(product)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseFwmV20250611Client().DescribeStrategyDispatchStatusWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("DescribeStrategyDispatchStatus returned an empty response."))
		}

		status = result.Response
		return nil
	})

	if err != nil {
		errRet = err
	}

	return
}

// WaitFwmStrategyDispatchComplete polls DescribeStrategyDispatchStatus until the
// dispatch task reaches a terminal state (succeeded=2 / failed=3), or the
// timeout budget is exhausted. It returns the last observed status snapshot so
// callers can attach it to a timeout/failure diagnostic.
func (me *FwmService) WaitFwmStrategyDispatchComplete(ctx context.Context, product string, timeout time.Duration, pollInterval time.Duration) (last *fwmv20250611.DescribeStrategyDispatchStatusResponseParams, err error) {
	deadline := time.Now().Add(timeout)
	for {
		latest, e := me.DescribeFwmStrategyDispatchStatus(ctx, product)
		if e != nil {
			err = e
			return
		}
		last = latest
		if latest != nil && latest.Status != nil {
			switch *latest.Status {
			case fwmDispatchStatusSucceeded:
				return
			case fwmDispatchStatusFailed:
				errMsg := "unknown"
				if latest.ErrorMsg != nil {
					errMsg = *latest.ErrorMsg
				}
				err = fmt.Errorf("dispatch strategy failed, ErrorMsg: %s", errMsg)
				return
			}
		}

		if time.Now().After(deadline) {
			err = fmt.Errorf("dispatch strategy did not converge within %s", timeout)
			return
		}
		time.Sleep(pollInterval)
	}
}
