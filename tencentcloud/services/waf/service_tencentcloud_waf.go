package waf

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	waf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/waf/v20180125"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type WafService struct {
	client *connectivity.TencentCloudClient
}

func (me *WafService) DescribeWafCustomRuleById(ctx context.Context, domain, ruleId string) (CustomRule *waf.DescribeCustomRulesRspRuleListItem, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := waf.NewDescribeCustomRuleListRequest()
	response := waf.NewDescribeCustomRuleListResponse()
	request.Domain = &domain
	request.Offset = common.Uint64Ptr(0)
	request.Limit = common.Uint64Ptr(20)
	request.Filters = []*waf.FiltersItemNew{
		{
			Name:       common.StringPtr("RuleId"),
			Values:     common.StringPtrs([]string{ruleId}),
			ExactMatch: common.BoolPtr(true),
		},
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseWafV20180125Client().DescribeCustomRuleList(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	if len(response.Response.RuleList) < 1 {
		return
	}

	CustomRule = response.Response.RuleList[0]
	return
}

func (me *WafService) DeleteWafCustomRuleById(ctx context.Context, domain, ruleId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := waf.NewDeleteCustomRuleRequest()
	request.Domain = &domain
	request.RuleId = &ruleId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseWafV20180125Client().DeleteCustomRule(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	return
}

func (me *WafService) DescribeWafCustomWhiteRuleById(ctx context.Context, domain, ruleId string) (CustomWhiteRule *waf.DescribeCustomRulesRspRuleListItem, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := waf.NewDescribeCustomWhiteRuleRequest()
	response := waf.NewDescribeCustomWhiteRuleResponse()
	request.Domain = &domain
	request.Offset = common.Uint64Ptr(0)
	request.Limit = common.Uint64Ptr(20)
	request.Filters = []*waf.FiltersItemNew{
		{
			Name:       common.StringPtr("RuleId"),
			Values:     common.StringPtrs([]string{ruleId}),
			ExactMatch: common.BoolPtr(true),
		},
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseWafClient().DescribeCustomWhiteRule(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.RuleList == nil {
			return resource.NonRetryableError(fmt.Errorf("Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	if err != nil {
		errRet = err
		return
	}

	if len(response.Response.RuleList) < 1 {
		return
	}

	CustomWhiteRule = response.Response.RuleList[0]
	return
}

func (me *WafService) DeleteWafCustomWhiteRuleById(ctx context.Context, domain, ruleId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := waf.NewDeleteCustomWhiteRuleRequest()
	request.Domain = &domain
	tmpRuleId, _ := strconv.ParseUint(ruleId, 10, 64)
	request.RuleId = &tmpRuleId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseWafClient().DeleteCustomWhiteRule(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	return
}

func (me *WafService) DescribeWafCiphersByFilter(ctx context.Context) (ciphers []*waf.TLSCiphers, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = waf.NewDescribeCiphersDetailRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWafClient().DescribeCiphersDetail(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || len(response.Response.Ciphers) < 1 {
		return
	}

	ciphers = response.Response.Ciphers
	return
}

func (me *WafService) DescribeWafTlsVersionsByFilter(ctx context.Context) (tlsVersions []*waf.TLSVersion, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = waf.NewDescribeTlsVersionRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWafClient().DescribeTlsVersion(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || len(response.Response.TLS) < 1 {
		return
	}

	tlsVersions = response.Response.TLS
	return
}

func (me *WafService) DescribeDomainsById(ctx context.Context, instanceID, domain string) (domainInfo *waf.DomainInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := waf.NewDescribeDomainsRequest()
	request.Offset = common.Uint64Ptr(0)
	request.Limit = common.Uint64Ptr(20)
	request.Filters = []*waf.FiltersItemNew{
		{
			Name:       common.StringPtr("InstanceId"),
			Values:     common.StringPtrs([]string{instanceID}),
			ExactMatch: common.BoolPtr(true),
		},
		{
			Name:       common.StringPtr("Domain"),
			Values:     common.StringPtrs([]string{domain}),
			ExactMatch: common.BoolPtr(true),
		},
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWafClient().DescribeDomains(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.Domains) < 1 {
		return
	}

	domainInfo = response.Response.Domains[0]
	return
}

func (me *WafService) DescribeWafClbDomainById(ctx context.Context, instanceID, domain, domainId string) (clbDomainInfo *waf.ClbDomainsInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := waf.NewDescribeDomainDetailsClbRequest()
	request.InstanceId = &instanceID
	request.Domain = &domain
	request.DomainId = &domainId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWafClient().DescribeDomainDetailsClb(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response.DomainsClbPartInfo == nil {
		return
	}

	clbDomainInfo = response.Response.DomainsClbPartInfo
	return
}

func (me *WafService) DeleteWafClbDomainById(ctx context.Context, instanceID, domain, domainId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := waf.NewDeleteHostRequest()
	request.HostsDel = []*waf.HostDel{
		{
			Domain:     common.StringPtr(domain),
			InstanceID: common.StringPtr(instanceID),
			DomainId:   common.StringPtr(domainId),
		},
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWafClient().DeleteHost(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *WafService) DescribeWafSaasDomainById(ctx context.Context, instanceID, domain, domainId string) (saasDomain *waf.DomainsPartInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := waf.NewDescribeDomainDetailsSaasRequest()
	request.InstanceId = &instanceID
	request.Domain = &domain
	request.DomainId = &domainId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWafClient().DescribeDomainDetailsSaas(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response.DomainsPartInfo == nil {
		return
	}

	saasDomain = response.Response.DomainsPartInfo
	return
}

func (me *WafService) DeleteWafSaasDomainById(ctx context.Context, instanceID, domain string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := waf.NewDeleteSpartaProtectionRequest()
	request.InstanceID = common.StringPtr(instanceID)
	request.Domains = common.StringPtrs([]string{domain})

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWafClient().DeleteSpartaProtection(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *WafService) DescribeWafDomainsByFilter(ctx context.Context, instanceID, domain string) (domains []*waf.DomainInfo, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = waf.NewDescribeDomainsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	tmpFilter := []*waf.FiltersItemNew{}
	if instanceID != "" {
		tmpFilter = append(tmpFilter, &waf.FiltersItemNew{
			Name:       common.StringPtr("InstanceId"),
			Values:     common.StringPtrs([]string{instanceID}),
			ExactMatch: common.BoolPtr(true),
		})
	}

	if domain != "" {
		tmpFilter = append(tmpFilter, &waf.FiltersItemNew{
			Name:       common.StringPtr("Domain"),
			Values:     common.StringPtrs([]string{domain}),
			ExactMatch: common.BoolPtr(true),
		})
	}

	request.Filters = tmpFilter

	ratelimit.Check(request.GetAction())

	var (
		offset uint64 = 0
		limit  uint64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseWafClient().DescribeDomains(request)
		if err != nil {
			errRet = err
			return
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Domains) < 1 {
			break
		}

		domains = append(domains, response.Response.Domains...)
		if len(response.Response.Domains) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *WafService) DescribeWafFindDomainsByFilter(ctx context.Context, param map[string]interface{}) (findDomains []*waf.FindAllDomainDetail, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = waf.NewDescribeFindDomainListRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "Key" {
			request.Key = v.(*string)
		}

		if k == "IsWafDomain" {
			request.IsWafDomain = v.(*string)
		}

		if k == "By" {
			request.By = v.(*string)
		}

		if k == "Order" {
			request.Order = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset uint64 = 1
		limit  uint64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseWafClient().DescribeFindDomainList(request)
		if err != nil {
			errRet = err
			return
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.List) < 1 {
			break
		}

		findDomains = append(findDomains, response.Response.List...)
		if len(response.Response.List) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *WafService) DescribeWafPortsByFilter(ctx context.Context, param map[string]interface{}) (ports *waf.DescribePortsResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = waf.NewDescribePortsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "Edition" {
			request.Edition = v.(*string)
		}

		if k == "InstanceID" {
			request.InstanceID = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWafClient().DescribePorts(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil {
		return
	}

	ports = response.Response
	return
}

func (me *WafService) DescribeWafUserDomainsByFilter(ctx context.Context) (userDomains []*waf.UserDomainInfo, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = waf.NewDescribeUserDomainInfoRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWafClient().DescribeUserDomainInfo(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || len(response.Response.UsersInfo) < 1 {
		return
	}

	userDomains = response.Response.UsersInfo
	return
}

func (me *WafService) DescribeWafInstanceById(ctx context.Context, instanceId string) (instance *waf.InstanceInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := waf.NewDescribeInstancesRequest()
	response := waf.NewDescribeInstancesResponse()
	request.Offset = common.Uint64Ptr(1)
	request.Limit = common.Uint64Ptr(20)
	request.Filters = []*waf.FiltersItemNew{
		{
			Name:       common.StringPtr("InstanceId"),
			Values:     common.StringPtrs([]string{instanceId}),
			ExactMatch: common.BoolPtr(true),
		},
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var iacExtInfo connectivity.IacExtInfo
	iacExtInfo.InstanceId = instanceId

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseWafClient(iacExtInfo).DescribeInstances(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	if len(response.Response.Instances) < 1 {
		return
	}

	instance = response.Response.Instances[0]
	return
}

func (me *WafService) DescribeWafInstanceWaitStatusById(ctx context.Context, instanceId string) error {
	logId := tccommon.GetLogId(ctx)

	request := waf.NewDescribeInstancesRequest()
	request.Offset = common.Uint64Ptr(1)
	request.Limit = common.Uint64Ptr(20)
	request.Filters = []*waf.FiltersItemNew{
		{
			Name:       common.StringPtr("InstanceId"),
			Values:     common.StringPtrs([]string{instanceId}),
			ExactMatch: common.BoolPtr(true),
		},
	}

	err := resource.Retry(tccommon.ReadRetryTimeout*10, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseWafClient().DescribeInstances(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || len(result.Response.Instances) < 1 {
			return resource.NonRetryableError(fmt.Errorf("DescribeInstances response is nil."))
		}

		instance := result.Response.Instances[0]
		if instance.Status != nil && *instance.Status == 0 {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("Waf instance still running, status is %d...", *instance.Status))
	})

	if err != nil {
		return err
	}

	return nil
}

func (me *WafService) DescribeWafAttackLogHistogramByFilter(ctx context.Context, param map[string]interface{}) (AttackLogHistogram *waf.GetAttackHistogramResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = waf.NewGetAttackHistogramRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "Domain" {
			request.Domain = v.(*string)
		}

		if k == "StartTime" {
			request.StartTime = v.(*string)
		}

		if k == "EndTime" {
			request.EndTime = v.(*string)
		}

		if k == "QueryString" {
			request.QueryString = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWafClient().GetAttackHistogram(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil {
		return
	}

	AttackLogHistogram = response.Response
	return
}

func (me *WafService) DescribeWafAttackLogListByFilter(ctx context.Context, param map[string]interface{}) (AttackLogList []*waf.AttackLogInfo, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = waf.NewSearchAttackLogRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "Domain" {
			request.Domain = v.(*string)
		}

		if k == "StartTime" {
			request.StartTime = v.(*string)
		}

		if k == "EndTime" {
			request.EndTime = v.(*string)
		}

		if k == "Count" {
			request.Count = v.(*int64)
		}

		if k == "QueryString" {
			request.QueryString = v.(*string)
		}

		if k == "Sort" {
			request.Sort = v.(*string)
		}

		if k == "Page" {
			request.Page = v.(*int64)
		}
	}

	request.Context = common.StringPtr("")

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWafClient().SearchAttackLog(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil {
		return
	}

	AttackLogList = response.Response.Data
	return
}

func (me *WafService) DescribeWafAttackOverviewByFilter(ctx context.Context, param map[string]interface{}) (AttackOverview *waf.DescribeAttackOverviewResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = waf.NewDescribeAttackOverviewRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "FromTime" {
			request.FromTime = v.(*string)
		}

		if k == "ToTime" {
			request.ToTime = v.(*string)
		}

		if k == "Appid" {
			request.Appid = v.(*uint64)
		}

		if k == "Domain" {
			request.Domain = v.(*string)
		}

		if k == "Edition" {
			request.Edition = v.(*string)
		}

		if k == "InstanceID" {
			request.InstanceID = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWafClient().DescribeAttackOverview(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil {
		return
	}

	AttackOverview = response.Response
	return
}

func (me *WafService) DescribeWafAttackTotalCountByFilter(ctx context.Context, param map[string]interface{}) (AttackTotalCount *waf.GetAttackTotalCountResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = waf.NewGetAttackTotalCountRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "StartTime" {
			request.StartTime = v.(*string)
		}

		if k == "EndTime" {
			request.EndTime = v.(*string)
		}

		if k == "Domain" {
			request.Domain = v.(*string)
		}

		if k == "QueryString" {
			request.QueryString = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWafClient().GetAttackTotalCount(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil {
		return
	}

	AttackTotalCount = response.Response
	return
}

func (me *WafService) DescribeWafPeakPointsByFilter(ctx context.Context, param map[string]interface{}) (PeakPoints []*waf.PeakPointsItem, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = waf.NewDescribePeakPointsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "FromTime" {
			request.FromTime = v.(*string)
		}

		if k == "ToTime" {
			request.ToTime = v.(*string)
		}

		if k == "Domain" {
			request.Domain = v.(*string)
		}

		if k == "Edition" {
			request.Edition = v.(*string)
		}

		if k == "InstanceID" {
			request.InstanceID = v.(*string)
		}

		if k == "MetricName" {
			request.MetricName = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWafClient().DescribePeakPoints(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil {
		return
	}

	PeakPoints = response.Response.Points
	return
}

func (me *WafService) DescribeWafAntiFakeById(ctx context.Context, id, domain string) (antiFake *waf.CacheUrlItems, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := waf.NewDescribeAntiFakeRulesRequest()
	request.Domain = &domain
	request.Offset = common.Uint64Ptr(0)
	request.Limit = common.Uint64Ptr(10)
	request.Filters = []*waf.FiltersItemNew{
		{
			Name:       common.StringPtr("RuleId"),
			Values:     common.StringPtrs([]string{id}),
			ExactMatch: common.BoolPtr(true),
		},
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWafClient().DescribeAntiFakeRules(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.Data) < 1 {
		return
	}

	antiFake = response.Response.Data[0]
	return
}

func (me *WafService) DeleteWafAntiFakeById(ctx context.Context, id, domain string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := waf.NewDeleteAntiFakeUrlRequest()
	idInt, _ := strconv.ParseUint(id, 10, 64)
	request.Id = common.Uint64Ptr(idInt)
	request.Domain = common.StringPtr(domain)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWafClient().DeleteAntiFakeUrl(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *WafService) DescribeWafAntiInfoLeakById(ctx context.Context, ruleId, domain string) (antiInfoLeak *waf.DescribeAntiLeakageItem, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := waf.NewDescribeAntiInfoLeakageRulesRequest()
	request.Domain = &domain
	request.Limit = common.Uint64Ptr(10)
	request.Offset = common.Uint64Ptr(0)
	request.Filters = []*waf.FiltersItemNew{
		{
			Name:       common.StringPtr("RuleId"),
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

	response, err := me.client.UseWafClient().DescribeAntiInfoLeakageRules(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.RuleList) < 1 {
		return
	}

	antiInfoLeak = response.Response.RuleList[0]
	return
}

func (me *WafService) DeleteWafAntiInfoLeakById(ctx context.Context, ruleId, domain string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := waf.NewDeleteAntiInfoLeakRuleRequest()
	ruleIdInt, _ := strconv.ParseUint(ruleId, 10, 64)
	request.Domain = &domain
	request.RuleId = &ruleIdInt

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWafClient().DeleteAntiInfoLeakRule(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *WafService) DescribeWafInstanceQpsLimitByFilter(ctx context.Context, param map[string]interface{}) (instanceQpsLimit *waf.QpsData, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = waf.NewGetInstanceQpsLimitRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "InstanceId" {
			request.InstanceId = v.(*string)
		}

		if k == "Type" {
			request.Type = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWafClient().GetInstanceQpsLimit(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil {
		return
	}

	instanceQpsLimit = response.Response.QpsData
	return
}

func (me *WafService) DescribeWafAutoDenyRulesById(ctx context.Context, domain string) (autoDenyRules *waf.DescribeWafAutoDenyRulesResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := waf.NewDescribeWafAutoDenyRulesRequest()
	request.Domain = &domain

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWafClient().DescribeWafAutoDenyRules(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	autoDenyRules = response.Response
	return
}

func (me *WafService) DescribeWafModuleStatusById(ctx context.Context, domain string) (moduleStatus *waf.DescribeModuleStatusResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := waf.NewDescribeModuleStatusRequest()
	request.Domain = &domain

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWafClient().DescribeModuleStatus(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	moduleStatus = response.Response
	return
}

func (me *WafService) DescribeSpartaProtectionInfoById(ctx context.Context, domain, edition string) (protectionInfo *waf.DescribeSpartaProtectionInfoResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := waf.NewDescribeSpartaProtectionInfoRequest()
	request.Domain = &domain
	request.Edition = &edition

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWafClient().DescribeSpartaProtectionInfo(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	protectionInfo = response.Response
	return
}

func (me *WafService) DescribeWafWebShellById(ctx context.Context, domain string) (webShell *waf.DescribeWebshellStatusResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := waf.NewDescribeWebshellStatusRequest()
	request.Domain = &domain

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWafClient().DescribeWebshellStatus(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	webShell = response.Response
	return
}

func (me *WafService) DescribeWafUserClbRegionsByFilter(ctx context.Context) (userClbRegions *waf.DescribeUserClbWafRegionsResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = waf.NewDescribeUserClbWafRegionsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWafClient().DescribeUserClbWafRegions(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil {
		return
	}

	userClbRegions = response.Response
	return
}

func (me *WafService) DescribeWafCcById(ctx context.Context, domain, ruleId string) (cc *waf.CCRuleItems, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := waf.NewDescribeCCRuleListRequest()
	response := waf.NewDescribeCCRuleListResponse()
	request.Domain = &domain
	request.Filters = []*waf.FiltersItemNew{
		{
			Name:       common.StringPtr("RuleId"),
			Values:     common.StringPtrs([]string{ruleId}),
			ExactMatch: common.BoolPtr(true),
		},
	}
	request.Offset = common.Uint64Ptr(0)
	request.Limit = common.Uint64Ptr(10)
	request.By = common.StringPtr("ts_version")

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseWafClient().DescribeCCRuleList(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.Data == nil {
			return resource.NonRetryableError(fmt.Errorf("Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	if len(response.Response.Data.Res) != 1 {
		return
	}

	cc = response.Response.Data.Res[0]
	return
}

func (me *WafService) DeleteWafCcById(ctx context.Context, domain, ruleId, name string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := waf.NewDeleteCCRuleRequest()
	request.Domain = common.StringPtr(domain)
	request.Name = common.StringPtr(name)
	ruleIdInt, _ := strconv.ParseInt(ruleId, 10, 64)
	request.RuleId = common.Int64Ptr(ruleIdInt)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseWafClient().DeleteCCRule(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	return
}

func (me *WafService) DescribeWafCcAutoStatusById(ctx context.Context, domain string) (CcAutoStatus *waf.DescribeCCAutoStatusResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := waf.NewDescribeCCAutoStatusRequest()
	request.Domain = &domain

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWafClient().DescribeCCAutoStatus(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	CcAutoStatus = response.Response
	return
}

func (me *WafService) DeleteWafCcAutoStatusById(ctx context.Context, domain, edition string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := waf.NewUpsertCCAutoStatusRequest()
	request.Domain = &domain
	request.Edition = &edition
	request.Value = helper.IntInt64(0)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWafClient().UpsertCCAutoStatus(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *WafService) DescribeWafCcSessionById(ctx context.Context, domain, edition, sessionID string) (ccSession *waf.SessionItem, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := waf.NewDescribeSessionRequest()
	request.Domain = &domain
	request.Edition = &edition
	sessionIDInt, _ := strconv.ParseInt(sessionID, 10, 64)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWafClient().DescribeSession(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || len(response.Response.Data.Res) < 1 {
		return
	}

	for _, item := range response.Response.Data.Res {
		if *item.SessionId == sessionIDInt {
			ccSession = item
			break
		}
	}

	return
}

func (me *WafService) DeleteWafCcSessionById(ctx context.Context, domain, edition, sessionID string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := waf.NewDeleteSessionRequest()
	request.Domain = &domain
	request.Edition = &edition
	sessionIDInt, _ := strconv.ParseInt(sessionID, 10, 64)
	request.SessionID = &sessionIDInt

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWafClient().DeleteSession(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *WafService) DescribeWafIpAccessControlById(ctx context.Context, domain string) (ipAccessControlList []*waf.IpAccessControlItem, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := waf.NewDescribeIpAccessControlRequest()
	request.Domain = &domain
	request.Count = common.Uint64Ptr(1)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	var (
		offset uint64 = 0
		limit  uint64 = 20
	)

	for {
		request.OffSet = &offset
		request.Limit = &limit
		response, err := me.client.UseWafClient().DescribeIpAccessControl(request)
		if err != nil {
			errRet = err
			return
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Data.Res) < 1 {
			break
		}

		ipAccessControlList = append(ipAccessControlList, response.Response.Data.Res...)
		if len(response.Response.Data.Res) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *WafService) DeleteWafIpAccessControlByDiff(ctx context.Context, domain string, ids []string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := waf.NewDeleteIpAccessControlRequest()
	request.Domain = &domain
	request.IsId = common.BoolPtr(true)
	request.Items = common.StringPtrs(ids)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWafClient().DeleteIpAccessControl(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *WafService) DeleteWafIpAccessControlById(ctx context.Context, domain string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := waf.NewDeleteIpAccessControlRequest()
	request.Domain = &domain
	request.Items = common.StringPtrs([]string{""})
	request.DeleteAll = common.BoolPtr(true)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWafClient().DeleteIpAccessControl(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *WafService) DescribeWafIpAccessControlV2ById(ctx context.Context, domain string, ruleId string) (ret *waf.DescribeIpAccessControlResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := waf.NewDescribeIpAccessControlRequest()
	request.Domain = helper.String(domain)
	request.Count = helper.Uint64(1)
	request.RuleId = helper.StrToUint64Point(ruleId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWafV20180125Client().DescribeIpAccessControl(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	ret = response.Response
	return
}

func (me *WafService) DescribeWafLogPostClsFlowById(ctx context.Context, logType int64) (ret *waf.DescribePostCLSFlowsResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := waf.NewDescribePostCLSFlowsRequest()
	response := waf.NewDescribePostCLSFlowsResponse()
	request.LogType = &logType

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseWafV20180125Client().DescribePostCLSFlows(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
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

func (me *WafService) DescribeWafLogPostCkafkaFlowById(ctx context.Context, logType int64) (ret *waf.DescribePostCKafkaFlowsResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := waf.NewDescribePostCKafkaFlowsRequest()
	response := waf.NewDescribePostCKafkaFlowsResponse()
	request.LogType = &logType

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseWafV20180125Client().DescribePostCKafkaFlows(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
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

func (me *WafService) DescribeWafDomainPostActionById(ctx context.Context, domain string) (domains []*waf.DomainInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := waf.NewDescribeDomainsRequest()
	response := waf.NewDescribeDomainsResponse()
	tmpFilter := []*waf.FiltersItemNew{}
	if domain != "" {
		tmpFilter = append(tmpFilter, &waf.FiltersItemNew{
			Name:       common.StringPtr("Domain"),
			Values:     common.StringPtrs([]string{domain}),
			ExactMatch: common.BoolPtr(true),
		})
	}

	request.Filters = tmpFilter

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var (
		offset uint64 = 0
		limit  uint64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			result, e := me.client.UseWafV20180125Client().DescribeDomains(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			response = result
			return nil
		})

		if err != nil {
			errRet = err
			return
		}

		if response == nil || len(response.Response.Domains) < 1 {
			break
		}

		domains = append(domains, response.Response.Domains...)
		if len(response.Response.Domains) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *WafService) DescribeWafBotSceneStatusConfigById(ctx context.Context, domain string, sceneId string) (ret *waf.BotSceneInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := waf.NewDescribeBotSceneListRequest()
	response := waf.NewDescribeBotSceneListResponse()
	request.Domain = &domain
	request.SceneId = &sceneId
	// wait waf sdk update
	// request.BusinessType = common.StringPtrs([]string{"all"})
	request.BusinessType = common.StringPtrs([]string{"login", "seckill", "crawl", "scan", "key-protect", "click-farming", "junk-mail", "social-media", "auto-download", "custom"})

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var (
		BotSceneList []*waf.BotSceneInfo
		offset       int64 = 0
		limit        int64 = 20
	)

	for {
		request.Offset = &offset
		request.Limit = &limit
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			result, e := me.client.UseWafV20180125Client().DescribeBotSceneList(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			response = result
			return nil
		})

		if err != nil {
			errRet = err
			return
		}

		if response == nil || len(response.Response.BotSceneList) < 1 {
			break
		}

		BotSceneList = append(BotSceneList, response.Response.BotSceneList...)
		if len(response.Response.BotSceneList) < int(limit) {
			break
		}

		offset += limit
	}

	for _, item := range BotSceneList {
		if *item.SceneId == sceneId {
			ret = item
			return
		}
	}

	return
}

func (me *WafService) DescribeWafBotStatusConfigById(ctx context.Context, domain string) (ret *waf.DescribeBotSceneOverviewResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := waf.NewDescribeBotSceneOverviewRequest()
	response := waf.NewDescribeBotSceneOverviewResponse()
	request.Domain = helper.String(domain)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseWafV20180125Client().DescribeBotSceneOverview(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
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

func (me *WafService) DescribeWafBotSceneUCBRuleById(ctx context.Context, domain, sceneId, ruleId string) (ret *waf.InOutputBotUCBRule, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := waf.NewDescribeBotSceneUCBRuleRequest()
	response := waf.NewDescribeBotSceneUCBRuleResponse()
	request.Domain = &domain
	request.SceneId = &sceneId
	request.RuleId = &ruleId
	request.Sort = helper.String("timestamp:-1")

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var (
		BotSceneUCBRuleList []*waf.InOutputBotUCBRule
		skip                uint64 = 0
		limit               uint64 = 20
	)

	for {
		request.Skip = &skip
		request.Limit = &limit
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			result, e := me.client.UseWafV20180125Client().DescribeBotSceneUCBRule(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			response = result
			return nil
		})

		if err != nil {
			errRet = err
			return
		}

		if response == nil || response.Response == nil || response.Response.Data == nil || len(response.Response.Data.Res) < 1 {
			break
		}

		BotSceneUCBRuleList = append(BotSceneUCBRuleList, response.Response.Data.Res...)
		if len(response.Response.Data.Res) < int(limit) {
			break
		}

		limit += skip
	}

	for _, item := range BotSceneUCBRuleList {
		if *item.SceneId == sceneId {
			ret = item
			return
		}
	}

	return
}

func (me *WafService) DeleteWafBotSceneUCBRuleById(ctx context.Context, domain, sceneId, ruleId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := waf.NewDeleteBotSceneUCBRuleRequest()
	request.Domain = helper.String(domain)
	request.SceneId = helper.String(sceneId)
	request.RuleId = helper.String(ruleId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseWafV20180125Client().DeleteBotSceneUCBRule(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	return
}

func (me *WafService) DescribeWafAttackWhiteRuleById(ctx context.Context, domain string, ruleId uint64) (ret *waf.UserWhiteRule, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := waf.NewDescribeAttackWhiteRuleRequest()
	response := waf.NewDescribeAttackWhiteRuleResponse()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	request.Domain = &domain

	var (
		offset uint64 = 0
		limit  uint64 = 20
		wrList []*waf.UserWhiteRule
	)

	for {
		request.Offset = &offset
		request.Limit = &limit
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			result, e := me.client.UseWafV20180125Client().DescribeAttackWhiteRule(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			response = result
			return nil
		})

		if err != nil {
			errRet = err
			return
		}

		if response == nil || len(response.Response.List) < 1 {
			break
		}

		wrList = append(wrList, response.Response.List...)
		if len(response.Response.List) < int(limit) {
			break
		}

		offset += limit
	}

	for _, item := range wrList {
		if item.WhiteRuleId != nil && *item.WhiteRuleId == ruleId {
			ret = item
			break
		}
	}

	return
}

func (me *WafService) DescribeWafOwaspRuleTypeConfigById(ctx context.Context, domain, typeId string) (ret *waf.OwaspRuleType, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := waf.NewDescribeOwaspRuleTypesRequest()
	response := waf.NewDescribeOwaspRuleTypesResponse()
	request.Domain = &domain

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var (
		offset uint64 = 0
		limit  uint64 = 100
		orList []*waf.OwaspRuleType
	)

	for {
		request.Offset = &offset
		request.Limit = &limit
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			result, e := me.client.UseWafV20180125Client().DescribeOwaspRuleTypes(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil || result.Response.List == nil {
				return resource.NonRetryableError(fmt.Errorf("Describe owasp rule types failed, Response is nil."))
			}

			response = result
			return nil
		})

		if err != nil {
			errRet = err
			return
		}

		if response == nil || len(response.Response.List) < 1 {
			break
		}

		orList = append(orList, response.Response.List...)
		if len(response.Response.List) < int(limit) {
			break
		}

		offset += limit
	}

	for _, item := range orList {
		if item != nil && item.TypeId != nil {
			respTypeId := helper.UInt64ToStr(*item.TypeId)
			if respTypeId == typeId {
				ret = item
				return
			}
		}
	}

	return
}

func (me *WafService) DescribeWafOwaspRuleTypesByFilter(ctx context.Context, param map[string]interface{}) (ret []*waf.OwaspRuleType, errRet error) {
	var (
		logId    = tccommon.GetLogId(ctx)
		request  = waf.NewDescribeOwaspRuleTypesRequest()
		response = waf.NewDescribeOwaspRuleTypesResponse()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "Domain" {
			request.Domain = v.(*string)
		}

		if k == "Filters" {
			request.Filters = v.([]*waf.FiltersItemNew)
		}
	}

	var (
		offset uint64 = 0
		limit  uint64 = 100
	)

	for {
		request.Offset = &offset
		request.Limit = &limit
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			result, e := me.client.UseWafV20180125Client().DescribeOwaspRuleTypes(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil || result.Response.List == nil {
				return resource.NonRetryableError(fmt.Errorf("Describe owasp rule types failed, Response is nil."))
			}

			response = result
			return nil
		})

		if err != nil {
			errRet = err
			return
		}

		if response == nil || len(response.Response.List) < 1 {
			break
		}

		ret = append(ret, response.Response.List...)
		if len(response.Response.List) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *WafService) DescribeWafOwaspRuleStatusConfigById(ctx context.Context, domain, ruleId string) (ret *waf.OwaspRule, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := waf.NewDescribeOwaspRulesRequest()
	response := waf.NewDescribeOwaspRulesResponse()
	request.Domain = &domain
	request.Filters = []*waf.FiltersItemNew{
		{
			Name:       common.StringPtr("RuleId"),
			Values:     common.StringPtrs([]string{ruleId}),
			ExactMatch: common.BoolPtr(true),
		},
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseWafV20180125Client().DescribeOwaspRules(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.List == nil || len(result.Response.List) < 1 {
			return resource.NonRetryableError(fmt.Errorf("Describe owasp rules failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	ret = response.Response.List[0]
	return
}

func (me *WafService) DescribeWafOwaspRulesByFilter(ctx context.Context, param map[string]interface{}) (ret []*waf.OwaspRule, errRet error) {
	var (
		logId    = tccommon.GetLogId(ctx)
		request  = waf.NewDescribeOwaspRulesRequest()
		response = waf.NewDescribeOwaspRulesResponse()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "Domain" {
			request.Domain = v.(*string)
		}

		if k == "By" {
			request.By = v.(*string)
		}

		if k == "Order" {
			request.Order = v.(*string)
		}

		if k == "Filters" {
			request.Filters = v.([]*waf.FiltersItemNew)
		}
	}

	var (
		offset uint64 = 0
		limit  uint64 = 100
	)

	for {
		request.Offset = &offset
		request.Limit = &limit
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			result, e := me.client.UseWafV20180125Client().DescribeOwaspRules(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil || result.Response.List == nil {
				return resource.NonRetryableError(fmt.Errorf("Describe owasp rules failed, Response is nil."))
			}

			response = result
			return nil
		})

		if err != nil {
			errRet = err
			return
		}

		if response == nil || len(response.Response.List) < 1 {
			break
		}

		ret = append(ret, response.Response.List...)
		if len(response.Response.List) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *WafService) DescribeWafOwaspWhiteRuleById(ctx context.Context, domain, ruleId string) (ret *waf.OwaspWhiteRule, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := waf.NewDescribeOwaspWhiteRulesRequest()
	response := waf.NewDescribeOwaspWhiteRulesResponse()
	request.Domain = &domain
	request.Filters = []*waf.FiltersItemNew{
		{
			Name:       common.StringPtr("RuleId"),
			Values:     common.StringPtrs([]string{ruleId}),
			ExactMatch: common.BoolPtr(true),
		},
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseWafV20180125Client().DescribeOwaspWhiteRules(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.List == nil || len(result.Response.List) < 1 {
			return resource.NonRetryableError(fmt.Errorf("Describe owasp white rules failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	ret = response.Response.List[0]
	return
}

func (me *WafService) DescribeWafObjectById(ctx context.Context, objectId string, role *string) (ret *waf.ClbObject, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := waf.NewDescribeObjectsRequest()
	response := waf.NewDescribeObjectsResponse()
	request.Filters = []*waf.FiltersItemNew{
		{
			Name:       common.StringPtr("ObjectId"),
			Values:     common.StringPtrs([]string{objectId}),
			ExactMatch: common.BoolPtr(true),
		},
	}

	if *role == "Admin" || *role == "DelegatedAdmin" {
		request.IsCrossAccount = common.Int64Ptr(1)
	} else if *role == "Member" || *role == "NoMember" {
		request.IsCrossAccount = common.Int64Ptr(0)
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseWafV20180125Client().DescribeObjects(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.ClbObjects == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe objects failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	if len(response.Response.ClbObjects) == 0 {
		return
	}

	ret = response.Response.ClbObjects[0]
	return
}

func (me *WafService) DescribeOrganizationRole(ctx context.Context) (ret *string, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := waf.NewGetOrganizationRoleRequest()
	response := waf.NewGetOrganizationRoleResponse()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseWafV20180125Client().GetOrganizationRole(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe organization role failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	ret = response.Response.Role
	return
}
