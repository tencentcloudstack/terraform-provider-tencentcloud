package tencentcloud

import (
	"context"
	"log"
	"strconv"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220106"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type TeoService struct {
	client *connectivity.TencentCloudClient
}

func (me *TeoService) DescribeTeoZone(ctx context.Context, zoneId string) (zone *teo.DescribeZoneDetailsResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = teo.NewDescribeZoneDetailsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()
	request.Id = &zoneId

	response, err := me.client.UseTeoClient().DescribeZoneDetails(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	zone = response.Response
	return
}

func (me *TeoService) DeleteTeoZoneById(ctx context.Context, zoneId string) (errRet error) {
	logId := getLogId(ctx)

	request := teo.NewDeleteZoneRequest()
	request.Id = &zoneId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTeoClient().DeleteZone(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TeoService) DescribeTeoDnsRecord(ctx context.Context, zoneId, name string) (dnsRecord *teo.DnsRecord, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = teo.NewDescribeDnsRecordsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	request.ZoneId = &zoneId
	request.Filters = append(
		request.Filters,
		&teo.DnsRecordFilter{
			Name:   helper.String("name"),
			Values: []*string{&name},
		},
	)
	ratelimit.Check(request.GetAction())

	var offset int64 = 0
	var pageSize int64 = 100
	instances := make([]*teo.DnsRecord, 0)

	for {
		request.Offset = &offset
		request.Limit = &pageSize
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseTeoClient().DescribeDnsRecords(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Records) < 1 {
			break
		}
		instances = append(instances, response.Response.Records...)
		if len(response.Response.Records) < int(pageSize) {
			break
		}
		offset += pageSize
	}

	if len(instances) < 1 {
		return
	}
	dnsRecord = instances[0]

	return

}

func (me *TeoService) DeleteTeoDnsRecordById(ctx context.Context, zoneId, dnsRecordId string) (errRet error) {
	logId := getLogId(ctx)

	request := teo.NewDeleteDnsRecordsRequest()
	request.Ids = []*string{&dnsRecordId}
	request.ZoneId = &zoneId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTeoClient().DeleteDnsRecords(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TeoService) DescribeTeoLoadBalancing(ctx context.Context, zoneId string, loadBalancingId string) (loadBalancing *teo.DescribeLoadBalancingDetailResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = teo.NewDescribeLoadBalancingDetailRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()
	request.ZoneId = &zoneId
	request.LoadBalancingId = &loadBalancingId

	response, err := me.client.UseTeoClient().DescribeLoadBalancingDetail(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	loadBalancing = response.Response
	return
}

func (me *TeoService) DeleteTeoLoadBalancingById(ctx context.Context, zoneId string, loadBalancingId string) (errRet error) {
	logId := getLogId(ctx)

	request := teo.NewDeleteLoadBalancingRequest()
	request.ZoneId = &zoneId
	request.LoadBalancingId = &loadBalancingId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTeoClient().DeleteLoadBalancing(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TeoService) DescribeTeoOriginGroup(ctx context.Context, zoneId string, originGroupId string) (originGroup *teo.DescribeOriginGroupDetailResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = teo.NewDescribeOriginGroupDetailRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()
	request.ZoneId = &zoneId
	request.OriginId = &originGroupId

	response, err := me.client.UseTeoClient().DescribeOriginGroupDetail(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	originGroup = response.Response
	return
}

func (me *TeoService) DeleteTeoOriginGroupById(ctx context.Context, zoneId string, originGroupId string) (errRet error) {
	logId := getLogId(ctx)

	request := teo.NewDeleteOriginGroupRequest()
	request.ZoneId = &zoneId
	request.OriginId = &originGroupId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTeoClient().DeleteOriginGroup(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TeoService) DescribeTeoRuleEngine(ctx context.Context, zoneId, ruleId string) (ruleEngine *teo.RuleSettingDetail, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = teo.NewDescribeRulesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	request.ZoneId = &zoneId
	request.Filters = append(
		request.Filters,
		&teo.RuleFilter{
			Name:   helper.String("RULE_ID"),
			Values: []*string{&ruleId},
		},
	)
	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTeoClient().DescribeRules(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	instances := response.Response.RuleList

	if len(instances) < 1 {
		return
	}
	ruleEngine = instances[0]

	return

}

func (me *TeoService) DeleteTeoRuleEngineById(ctx context.Context, zoneId, ruleId string) (errRet error) {
	logId := getLogId(ctx)

	request := teo.NewDeleteRulesRequest()

	request.ZoneId = &zoneId
	request.RuleIds = []*string{&ruleId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTeoClient().DeleteRules(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TeoService) DescribeTeoApplicationProxy(ctx context.Context, zoneId, proxyId string) (applicationProxy *teo.DescribeApplicationProxyDetailResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = teo.NewDescribeApplicationProxyDetailRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()
	request.ZoneId = &zoneId
	request.ProxyId = &proxyId

	response, err := me.client.UseTeoClient().DescribeApplicationProxyDetail(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	applicationProxy = response.Response
	return
}

func (me *TeoService) DeleteTeoApplicationProxyById(ctx context.Context, zoneId, proxyId string) (errRet error) {
	logId := getLogId(ctx)

	request := teo.NewDeleteApplicationProxyRequest()

	request.ZoneId = &zoneId
	request.ProxyId = &proxyId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTeoClient().DeleteApplicationProxy(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TeoService) DescribeTeoApplicationProxyRule(ctx context.Context, zoneId, proxyId, ruleId string) (applicationProxyRule *teo.ApplicationProxyRule, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = teo.NewDescribeApplicationProxyDetailRequest()
	)

	rules := make([]*teo.ApplicationProxyRule, 0)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()
	request.ZoneId = &zoneId
	request.ProxyId = &proxyId

	response, err := me.client.UseTeoClient().DescribeApplicationProxyDetail(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	rules = response.Response.Rule
	for _, rule := range rules {
		if *rule.RuleId == ruleId {
			applicationProxyRule = rule
			return
		}
	}
	return
}

func (me *TeoService) DeleteTeoApplicationProxyRuleById(ctx context.Context, zoneId, proxyId, ruleId string) (errRet error) {
	logId := getLogId(ctx)

	request := teo.NewDeleteApplicationProxyRuleRequest()

	request.ZoneId = &zoneId
	request.ProxyId = &proxyId
	request.RuleId = &ruleId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTeoClient().DeleteApplicationProxyRule(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TeoService) DescribeTeoZoneSetting(ctx context.Context, zoneId string) (zoneSetting *teo.DescribeZoneSettingResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = teo.NewDescribeZoneSettingRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()
	request.ZoneId = &zoneId

	response, err := me.client.UseTeoClient().DescribeZoneSetting(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	zoneSetting = response.Response
	return
}

func (me *TeoService) DescribeTeoSecurityPolicy(ctx context.Context, zoneId, entity string) (securityPolicy *teo.DescribeSecurityPolicyResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = teo.NewDescribeSecurityPolicyRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()
	request.ZoneId = &zoneId
	request.Entity = &entity

	response, err := me.client.UseTeoClient().DescribeSecurityPolicy(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	securityPolicy = response.Response
	return
}

func (me *TeoService) DescribeTeoHostCertificate(ctx context.Context, zoneId, host string) (hostCertificate *teo.HostCertSetting, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = teo.NewDescribeHostsCertificateRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	request.ZoneId = &zoneId

	request.Filters = append(
		request.Filters,
		&teo.CertFilter{
			Name:   helper.String("host"),
			Values: []*string{&host},
		},
	)
	ratelimit.Check(request.GetAction())

	var offset int64 = 0
	var pageSize int64 = 100
	instances := make([]*teo.HostCertSetting, 0)

	for {
		request.Offset = &offset
		request.Limit = &pageSize
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseTeoClient().DescribeHostsCertificate(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Hosts) < 1 {
			break
		}
		instances = append(instances, response.Response.Hosts...)
		if len(response.Response.Hosts) < int(pageSize) {
			break
		}
		offset += pageSize
	}

	if len(instances) < 1 {
		return
	}
	hostCertificate = instances[0]

	return
}

func (me *TeoService) DescribeTeoDnsSec(ctx context.Context, zoneId string) (dnsSec *teo.DescribeDnssecResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = teo.NewDescribeDnssecRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()
	request.Id = &zoneId

	response, err := me.client.UseTeoClient().DescribeDnssec(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	dnsSec = response.Response
	return
}

func (me *TeoService) DescribeTeoDefaultCertificate(ctx context.Context, zoneId string) (defaultCertificate *teo.DefaultServerCertInfo, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = teo.NewDescribeDefaultCertificatesRequest()
	)

	defaultCertificates := make([]*teo.DefaultServerCertInfo, 0)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()
	request.ZoneId = &zoneId

	response, err := me.client.UseTeoClient().DescribeDefaultCertificates(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	defaultCertificates = response.Response.CertInfo

	for _, cert := range defaultCertificates {
		if *cert.CertId != "" {
			defaultCertificate = cert
			return
		}
	}
	return
}

func (me *TeoService) DescribeTeoDdosPolicy(ctx context.Context, zoneId, policyId string) (ddosPolicy *teo.DescribeDDoSPolicyResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = teo.NewDescribeDDoSPolicyRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	policyId64, errRet := strconv.ParseInt(policyId, 10, 64)
	if errRet != nil {
		log.Printf("[DEBUG]%s api[%s] error, Type conversion failed, [%s] conversion int64 failed\n",
			logId, request.GetAction(), policyId)
		return nil, errRet
	}

	request.ZoneId = &zoneId
	request.PolicyId = &policyId64

	response, err := me.client.UseTeoClient().DescribeDDoSPolicy(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	ddosPolicy = response.Response
	return
}

func (me *TeoService) DescribeZoneDDoSPolicy(ctx context.Context, zoneId string) (ddosPolicy *teo.DescribeZoneDDoSPolicyResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = teo.NewDescribeZoneDDoSPolicyRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	if zoneId != "" {
		request.ZoneId = &zoneId
	}

	response, err := me.client.UseTeoClient().DescribeZoneDDoSPolicy(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	ddosPolicy = response.Response
	return
}

func (me *TeoService) DeleteTeoDdosPolicyById(ctx context.Context, zoneId, policyId string) (errRet error) {
	logId := getLogId(ctx)

	request := teo.NewModifyDDoSPolicyRequest()

	policyId64, errRet := strconv.ParseInt(policyId, 10, 64)
	if errRet != nil {
		log.Printf("[DEBUG]%s api[%s] error, Type conversion failed, [%s] conversion int64 failed\n",
			logId, request.GetAction(), policyId)
		return errRet
	}

	request.ZoneId = &zoneId
	request.PolicyId = &policyId64

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTeoClient().ModifyDDoSPolicy(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TeoService) DescribeAvailablePlans(ctx context.Context) (availablePlans *teo.DescribeAvailablePlansResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = teo.NewDescribeAvailablePlansRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	response, err := me.client.UseTeoClient().DescribeAvailablePlans(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	availablePlans = response.Response
	return
}
