package tencentcloud

import (
	"context"
	"log"
	"strconv"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	waf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/waf/v20180125"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type WafService struct {
	client *connectivity.TencentCloudClient
}

func (me *WafService) DescribeWafCustomRuleById(ctx context.Context, domain, ruleId string) (CustomRule *waf.DescribeCustomRulesRspRuleListItem, errRet error) {
	logId := getLogId(ctx)

	request := waf.NewDescribeCustomRuleListRequest()
	request.Domain = &domain
	request.Offset = common.Uint64Ptr(0)
	request.Limit = common.Uint64Ptr(20)
	request.Filters = []*waf.FiltersItemNew{
		{
			Name:       common.StringPtr("RuleID"),
			Values:     common.StringPtrs([]string{ruleId}),
			ExactMatch: common.BoolPtr(true),
		},
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWafClient().DescribeCustomRuleList(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.RuleList) < 1 {
		return
	}

	CustomRule = response.Response.RuleList[0]
	return
}

func (me *WafService) DeleteWafCustomRuleById(ctx context.Context, domain, ruleId string) (errRet error) {
	logId := getLogId(ctx)

	request := waf.NewDeleteCustomRuleRequest()
	request.Domain = &domain
	request.RuleId = &ruleId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWafClient().DeleteCustomRule(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *WafService) DescribeWafCustomWhiteRuleById(ctx context.Context, domain, ruleId string) (CustomWhiteRule *waf.DescribeCustomRulesRspRuleListItem, errRet error) {
	logId := getLogId(ctx)

	request := waf.NewDescribeCustomWhiteRuleRequest()
	request.Domain = &domain
	request.Offset = common.Uint64Ptr(0)
	request.Limit = common.Uint64Ptr(20)
	request.Filters = []*waf.FiltersItemNew{
		{
			Name:       common.StringPtr("RuleID"),
			Values:     common.StringPtrs([]string{ruleId}),
			ExactMatch: common.BoolPtr(true),
		},
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWafClient().DescribeCustomWhiteRule(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.RuleList) < 1 {
		return
	}

	CustomWhiteRule = response.Response.RuleList[0]
	return
}

func (me *WafService) DeleteWafCustomWhiteRuleById(ctx context.Context, domain, ruleId string) (errRet error) {
	logId := getLogId(ctx)

	request := waf.NewDeleteCustomWhiteRuleRequest()
	request.Domain = &domain
	tmpRuleId, _ := strconv.ParseUint(ruleId, 10, 64)
	request.RuleId = &tmpRuleId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWafClient().DeleteCustomWhiteRule(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}
