package tencentcloud

import (
	"context"
	"log"
	"strconv"

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

func (me *WafService) DescribeWafCiphersByFilter(ctx context.Context) (ciphers []*waf.TLSCiphers, errRet error) {
	var (
		logId   = getLogId(ctx)
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
		logId   = getLogId(ctx)
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
	logId := getLogId(ctx)

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
	logId := getLogId(ctx)

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
	logId := getLogId(ctx)

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
	logId := getLogId(ctx)

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
	logId := getLogId(ctx)

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
		logId   = getLogId(ctx)
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
		logId   = getLogId(ctx)
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

func (me *WafService) DescribeWafWafInfosByFilter(ctx context.Context, param map[string]interface{}) (wafInfos []*waf.ClbHostResult, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = waf.NewDescribeWafInfoRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "Params" {
			request.Params = v.([]*waf.ClbHostsParams)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWafClient().DescribeWafInfo(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || len(response.Response.HostList) < 1 {
		return
	}

	wafInfos = response.Response.HostList
	return
}

func (me *WafService) DescribeWafPortsByFilter(ctx context.Context, param map[string]interface{}) (ports *waf.DescribePortsResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
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
		logId   = getLogId(ctx)
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
	logId := getLogId(ctx)

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

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWafClient().DescribeInstances(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.Instances) < 1 {
		return
	}

	instance = response.Response.Instances[0]
	return
}

func (me *WafService) DescribeWafAttackLogHistogramByFilter(ctx context.Context, param map[string]interface{}) (AttackLogHistogram *waf.GetAttackHistogramResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
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
		logId   = getLogId(ctx)
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
		logId   = getLogId(ctx)
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
		logId   = getLogId(ctx)
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
		logId   = getLogId(ctx)
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
	logId := getLogId(ctx)

	request := waf.NewDescribeAntiFakeRulesRequest()
	request.Domain = &domain
	request.Offset = common.Uint64Ptr(0)
	request.Limit = common.Uint64Ptr(10)
	request.Filters = []*waf.FiltersItemNew{
		{
			Name:       common.StringPtr("RuleID"),
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
	logId := getLogId(ctx)

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
	logId := getLogId(ctx)

	request := waf.NewDescribeAntiInfoLeakageRulesRequest()
	request.Domain = &domain

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

	ruleIdInt, _ := strconv.ParseUint(ruleId, 10, 64)
	for _, item := range response.Response.RuleList {
		if *item.RuleId == ruleIdInt {
			antiInfoLeak = item
			break
		}
	}

	return
}

func (me *WafService) DeleteWafAntiInfoLeakById(ctx context.Context, ruleId, domain string) (errRet error) {
	logId := getLogId(ctx)

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
		logId   = getLogId(ctx)
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
	logId := getLogId(ctx)

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
	logId := getLogId(ctx)

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
	logId := getLogId(ctx)

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
	logId := getLogId(ctx)

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
		logId   = getLogId(ctx)
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
	logId := getLogId(ctx)

	request := waf.NewDescribeCCRuleListRequest()
	request.Domain = &domain
	request.Filters = []*waf.FiltersItemNew{
		{
			Name:       common.StringPtr("RuleID"),
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

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWafClient().DescribeCCRuleList(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || len(response.Response.Data.Res) != 1 {
		return
	}

	cc = response.Response.Data.Res[0]
	return
}

func (me *WafService) DeleteWafCcById(ctx context.Context, domain, ruleId, name string) (errRet error) {
	logId := getLogId(ctx)

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

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWafClient().DeleteCCRule(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *WafService) DescribeWafCcAutoStatusById(ctx context.Context, domain string) (CcAutoStatus *waf.DescribeCCAutoStatusResponseParams, errRet error) {
	logId := getLogId(ctx)

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
	logId := getLogId(ctx)

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
	logId := getLogId(ctx)

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
	logId := getLogId(ctx)

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
	logId := getLogId(ctx)

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
	logId := getLogId(ctx)

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
	logId := getLogId(ctx)

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
