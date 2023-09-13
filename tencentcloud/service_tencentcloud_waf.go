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
