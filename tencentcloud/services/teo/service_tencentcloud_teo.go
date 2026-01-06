package teo

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func NewTeoService(client *connectivity.TencentCloudClient) TeoService {
	return TeoService{client: client}
}

type TeoService struct {
	client *connectivity.TencentCloudClient
}

func (me *TeoService) DescribeTeoZone(ctx context.Context, zoneId string) (zone *teo.Zone, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = teo.NewDescribeZonesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	if zoneId != "" {
		request.Filters = append(
			request.Filters,
			&teo.AdvancedFilter{
				Name:   helper.String("zone-id"),
				Values: []*string{&zoneId},
			},
		)
	}

	ratelimit.Check(request.GetAction())
	var iacExtInfo connectivity.IacExtInfo
	iacExtInfo.InstanceId = zoneId

	var offset int64 = 0
	var pageSize int64 = 100
	instances := make([]*teo.Zone, 0)

	for {
		request.Offset = &offset
		request.Limit = &pageSize
		ratelimit.Check(request.GetAction())
		response := teo.NewDescribeZonesResponse()
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			result, e := me.client.UseTeoClient(iacExtInfo).DescribeZones(request)
			if e != nil {
				return tccommon.RetryError(e)
			}
			response = result
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || response.Response == nil || len(response.Response.Zones) < 1 {
			break
		}
		instances = append(instances, response.Response.Zones...)
		if len(response.Response.Zones) < int(pageSize) {
			break
		}
		offset += pageSize
	}

	if len(instances) < 1 {
		return
	}
	zone = instances[0]

	return
}

func (me *TeoService) DeleteTeoZoneById(ctx context.Context, zoneId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := teo.NewDeleteZoneRequest()
	request.ZoneId = &zoneId

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

func (me *TeoService) DescribeTeoOriginGroup(ctx context.Context,
	zoneId, originGroupId string) (originGroup *teo.OriginGroup, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = teo.NewDescribeOriginGroupRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	request.Filters = append(
		request.Filters,
		&teo.AdvancedFilter{
			Name:   helper.String("zone-id"),
			Values: []*string{&zoneId},
		},
	)
	request.Filters = append(
		request.Filters,
		&teo.AdvancedFilter{
			Name:   helper.String("origin-group-id"),
			Values: []*string{&originGroupId},
		},
	)

	var offset uint64 = 0
	var pageSize uint64 = 100
	originGroups := make([]*teo.OriginGroup, 0)

	for {
		request.Offset = &offset
		request.Limit = &pageSize
		ratelimit.Check(request.GetAction())
		response := teo.NewDescribeOriginGroupResponse()
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			result, e := me.client.UseTeoClient().DescribeOriginGroup(request)
			if e != nil {
				return tccommon.RetryError(e)
			}
			response = result
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || response.Response == nil || len(response.Response.OriginGroups) < 1 {
			break
		}
		originGroups = append(originGroups, response.Response.OriginGroups...)
		if len(response.Response.OriginGroups) < int(pageSize) {
			break
		}
		offset += pageSize
	}

	if len(originGroups) < 1 {
		return
	}
	originGroup = originGroups[0]

	return
}

func (me *TeoService) DeleteTeoOriginGroupById(ctx context.Context, zoneId, originGroupId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := teo.NewDeleteOriginGroupRequest()
	request.ZoneId = &zoneId
	request.GroupId = &originGroupId

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

func (me *TeoService) DescribeTeoRuleEngine(ctx context.Context, zoneId, ruleId string) (ruleEngine *teo.RuleItem,
	errRet error) {
	var (
		logId    = tccommon.GetLogId(ctx)
		request  = teo.NewDescribeRulesRequest()
		response = teo.NewDescribeRulesResponse()
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
		&teo.Filter{
			Name:   helper.String("rule-id"),
			Values: []*string{&ruleId},
		},
	)
	ratelimit.Check(request.GetAction())
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := me.client.UseTeoClient().DescribeRules(request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response != nil && response.Response != nil && response.Response.RuleItems != nil {
		for _, v := range response.Response.RuleItems {
			if *v.RuleId == ruleId {
				ruleEngine = v
				return
			}
		}
	}

	return

}

func (me *TeoService) DescribeTeoRuleEngines(ctx context.Context, zoneId string) (ruleEngines []*teo.RuleItem,
	errRet error) {
	var (
		logId    = tccommon.GetLogId(ctx)
		request  = teo.NewDescribeRulesRequest()
		response = teo.NewDescribeRulesResponse()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	request.ZoneId = &zoneId
	ratelimit.Check(request.GetAction())
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := me.client.UseTeoClient().DescribeRules(request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response != nil && response.Response != nil && response.Response.RuleItems != nil {
		ruleEngines = response.Response.RuleItems
	}

	return

}

func (me *TeoService) DeleteTeoRuleEngineById(ctx context.Context, zoneId, ruleId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

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

func (me *TeoService) DescribeTeoApplicationProxy(ctx context.Context,
	zoneId, proxyId string) (applicationProxy *teo.ApplicationProxy, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = teo.NewDescribeApplicationProxiesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	if zoneId != "" {
		request.Filters = append(
			request.Filters,
			&teo.Filter{
				Name:   helper.String("zone-id"),
				Values: []*string{&zoneId},
			},
		)
	}

	if proxyId != "" {
		request.Filters = append(
			request.Filters,
			&teo.Filter{
				Name:   helper.String("proxy-id"),
				Values: []*string{&proxyId},
			},
		)
	}

	ratelimit.Check(request.GetAction())

	var offset int64 = 0
	var pageSize int64 = 100
	instances := make([]*teo.ApplicationProxy, 0)

	for {
		request.Offset = &offset
		request.Limit = &pageSize
		ratelimit.Check(request.GetAction())
		response := teo.NewDescribeApplicationProxiesResponse()
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			result, e := me.client.UseTeoClient().DescribeApplicationProxies(request)
			if e != nil {
				return tccommon.RetryError(e)
			}
			response = result
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || response.Response == nil || len(response.Response.ApplicationProxies) < 1 {
			break
		}
		instances = append(instances, response.Response.ApplicationProxies...)
		if len(response.Response.ApplicationProxies) < int(pageSize) {
			break
		}
		offset += pageSize
	}

	if len(instances) < 1 {
		return
	}
	applicationProxy = instances[0]

	return
}

func (me *TeoService) DeleteTeoApplicationProxyById(ctx context.Context, zoneId, proxyId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

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

func (me *TeoService) DescribeTeoApplicationProxyRule(ctx context.Context,
	zoneId, proxyId, ruleId string) (applicationProxyRule *teo.ApplicationProxyRule, errRet error) {
	var (
		logId    = tccommon.GetLogId(ctx)
		request  = teo.NewDescribeApplicationProxiesRequest()
		response = teo.NewDescribeApplicationProxiesResponse()
	)

	request.Filters = append(
		request.Filters,
		&teo.Filter{
			Name:   helper.String("zone-id"),
			Values: []*string{&zoneId},
		},
	)
	request.Filters = append(
		request.Filters,
		&teo.Filter{
			Name:   helper.String("proxy-id"),
			Values: []*string{&proxyId},
		},
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := me.client.UseTeoClient().DescribeApplicationProxies(request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response == nil || len(response.Response.ApplicationProxies) < 1 {
		return
	}
	for _, v := range response.Response.ApplicationProxies[0].ApplicationProxyRules {
		if *v.RuleId == ruleId {
			applicationProxyRule = v
			return
		}
	}

	return
}

func (me *TeoService) DeleteTeoApplicationProxyRuleById(ctx context.Context,
	zoneId, proxyId, ruleId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

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

func (me *TeoService) DescribeTeoZoneSetting(ctx context.Context, zoneId string) (zoneSetting *teo.ZoneSetting,
	errRet error) {
	var (
		logId    = tccommon.GetLogId(ctx)
		request  = teo.NewDescribeZoneSettingRequest()
		response = teo.NewDescribeZoneSettingResponse()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()
	request.ZoneId = &zoneId
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := me.client.UseTeoClient().DescribeZoneSetting(request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response == nil {
		return
	}
	zoneSetting = response.Response.ZoneSetting
	return
}

func (me *TeoService) DescribeTeoDefaultCertificate(ctx context.Context,
	zoneId, certId string) (defaultCertificate *teo.DefaultServerCertInfo, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = teo.NewDescribeDefaultCertificatesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	request.Filters = append(
		request.Filters,
		&teo.Filter{
			Name:   helper.String("zone-id"),
			Values: []*string{&zoneId},
		},
	)

	var offset int64 = 0
	var pageSize int64 = 100
	certificates := make([]*teo.DefaultServerCertInfo, 0)

	for {
		request.Offset = &offset
		request.Limit = &pageSize
		ratelimit.Check(request.GetAction())
		response := teo.NewDescribeDefaultCertificatesResponse()
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			result, e := me.client.UseTeoClient().DescribeDefaultCertificates(request)
			if e != nil {
				return tccommon.RetryError(e)
			}
			response = result
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || response.Response == nil || len(response.Response.DefaultServerCertInfo) < 1 {
			break
		}
		certificates = append(certificates, response.Response.DefaultServerCertInfo...)
		if len(response.Response.DefaultServerCertInfo) < int(pageSize) {
			break
		}
		offset += pageSize
	}

	if len(certificates) < 1 {
		return
	}
	for _, v := range certificates {
		if *v.CertId == certId {
			defaultCertificate = v
			return
		}
	}

	return
}

func (me *TeoService) DescribeTeoZoneAvailablePlansByFilter(ctx context.Context, param map[string]interface{}) (ret []*teo.PlanInfo, errRet error) {
	var (
		logId    = tccommon.GetLogId(ctx)
		request  = teo.NewDescribeAvailablePlansRequest()
		response = teo.NewDescribeAvailablePlansResponse()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := me.client.UseTeoClient().DescribeAvailablePlans(request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		response = result
		return nil
	})
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response == nil || len(response.Response.PlanInfo) < 1 {
		return
	}

	ret = response.Response.PlanInfo
	return
}

func (me *TeoService) DescribeTeoRuleEnginePriority(ctx context.Context,
	zoneId string) (ruleEnginePriority []*teo.RuleItem, errRet error) {
	var (
		logId    = tccommon.GetLogId(ctx)
		request  = teo.NewDescribeRulesRequest()
		response = teo.NewDescribeRulesResponse()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()
	request.ZoneId = &zoneId

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := me.client.UseTeoClient().DescribeRules(request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response == nil {
		return
	}
	ruleEnginePriority = response.Response.RuleItems
	return
}

func (me *TeoService) DescribeTeoRuleEngineSettingsByFilter(ctx context.Context, param map[string]interface{}) (ret []*teo.RulesSettingAction, errRet error) {
	var (
		logId    = tccommon.GetLogId(ctx)
		request  = teo.NewDescribeRulesSettingRequest()
		response = teo.NewDescribeRulesSettingResponse()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := me.client.UseTeoClient().DescribeRulesSetting(request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		response = result
		return nil
	})
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response == nil || len(response.Response.Actions) < 1 {
		return
	}

	ret = response.Response.Actions
	return
}

func (me *TeoService) CheckZoneComplete(ctx context.Context, zoneId string) error {
	zone, err := me.DescribeTeoZone(ctx, zoneId)
	if err != nil {
		return err
	}
	if zone == nil || zone.Type == nil || zone.Status == nil || zone.CnameStatus == nil {
		return fmt.Errorf("get zone[%s] info failed", zoneId)
	}
	if *zone.Type == "full" && *zone.Status != "active" {
		return fmt.Errorf("`zone.Status` is not `active`, please modify NS records from the domain name provider first")
	}
	if *zone.Type == "partial" && *zone.CnameStatus != "finished" {
		return fmt.Errorf("`zone.CnameStatus` is not `finished`, please verify ownership of the site first")
	}
	return nil
}

func (me *TeoService) DescribeTeoAccelerationDomainById(ctx context.Context, zoneId string, domainName string) (ret *teo.AccelerationDomain, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := teo.NewDescribeAccelerationDomainsRequest()
	request.ZoneId = helper.String(zoneId)
	advancedFilter := &teo.AdvancedFilter{
		Name:   helper.String("domain-name"),
		Values: []*string{helper.String(domainName)},
	}
	request.Filters = append(request.Filters, advancedFilter)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	var instances []*teo.AccelerationDomain
	for {
		request.Offset = &offset
		request.Limit = &limit
		response := teo.NewDescribeAccelerationDomainsResponse()
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			result, e := me.client.UseTeoClient().DescribeAccelerationDomains(request)
			if e != nil {
				return tccommon.RetryError(e)
			}
			response = result
			return nil
		})
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.AccelerationDomains) < 1 {
			break
		}
		instances = append(instances, response.Response.AccelerationDomains...)
		if len(response.Response.AccelerationDomains) < int(limit) {
			break
		}

		offset += limit
	}

	if len(instances) < 1 {
		return
	}

	ret = instances[0]
	return
}

func (me *TeoService) DescribeIdentifications(ctx context.Context, domain string) (identifications []*teo.Identification, errRet error) {
	var (
		logId    = tccommon.GetLogId(ctx)
		request  = teo.NewDescribeIdentificationsRequest()
		response = teo.NewDescribeIdentificationsResponse()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()
	request.Filters = append(
		request.Filters,
		&teo.Filter{
			Name:   helper.String("zone-name"),
			Values: []*string{&domain},
		},
	)

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := me.client.UseTeoClient().DescribeIdentifications(request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response == nil {
		return nil, nil
	}

	identifications = response.Response.Identifications
	return
}

func (me *TeoService) ModifyZoneStatus(ctx context.Context, zoneId string, paused bool, operate string) error {
	logId := tccommon.GetLogId(ctx)

	req := teo.NewModifyZoneStatusRequest()
	req.ZoneId, req.Paused = &zoneId, helper.Bool(paused)
	_, e := me.client.UseTeoClient().ModifyZoneStatus(req)
	if e != nil {
		log.Printf("[CRITAL]%s modify zone status failed, reason:%+v", logId, e)
		return e
	}

	err := resource.Retry(6*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		instance, errRet := me.DescribeTeoZone(ctx, zoneId)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}
		if operate == "delete" {
			if *instance.ActiveStatus == "paused" {
				return nil
			}
		} else {
			if *instance.ActiveStatus == "inactive" || *instance.ActiveStatus == "paused" {
				return nil
			}
		}
		return resource.RetryableError(fmt.Errorf("zone status is %v, retry...", *instance.ActiveStatus))
	})
	if err != nil {
		return err
	}

	return nil
}

func (me *TeoService) CheckAccelerationDomainStatus(ctx context.Context, zoneId, domainName, operate string) error {
	err := resource.Retry(6*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		instance, errRet := me.DescribeTeoAccelerationDomainById(ctx, zoneId, domainName)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}
		if operate == "delete" {
			if *instance.DomainStatus == "offline" {
				return nil
			}
		} else {
			if *instance.DomainStatus == "online" {
				return nil
			}
		}

		return resource.RetryableError(fmt.Errorf("AccelerationDomain status is %v, retry...", *instance.DomainStatus))
	})
	if err != nil {
		return err
	}

	return nil
}

func (me *TeoService) DescribeTeoApplicationProxyRuleById(ctx context.Context, ruleId string) (ret *teo.ApplicationProxyRule, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := teo.NewDescribeApplicationProxiesRequest()
	response := teo.NewDescribeApplicationProxiesResponse()

	if err := resourceTencentCloudTeoApplicationProxyRuleReadPostFillRequest0(ctx, request); err != nil {
		return nil, err
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := me.client.UseTeoClient().DescribeApplicationProxies(request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		response = result
		return nil
	})
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	var tmpRet *teo.ApplicationProxy
	if response.Response == nil || len(response.Response.ApplicationProxies) < 1 {
		return
	}

	tmpRet = response.Response.ApplicationProxies[0]
	if len(tmpRet.ApplicationProxyRules) < 1 {
		return
	}

	for _, info := range tmpRet.ApplicationProxyRules {
		if info.RuleId != nil && *info.RuleId == ruleId {
			ret = info
			break
		}
	}
	return
}

func (me *TeoService) DescribeTeoZoneById(ctx context.Context, zoneId string) (ret *teo.Zone, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := teo.NewDescribeZonesRequest()
	advancedFilter := &teo.AdvancedFilter{
		Name:   helper.String("zone-id"),
		Values: []*string{&zoneId},
	}
	request.Filters = append(request.Filters, advancedFilter)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	var instances []*teo.Zone
	for {
		request.Offset = &offset
		request.Limit = &limit
		response := teo.NewDescribeZonesResponse()
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			result, e := me.client.UseTeoClient().DescribeZones(request)
			if e != nil {
				return tccommon.RetryError(e)
			}
			response = result
			return nil
		})
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Zones) < 1 {
			break
		}
		instances = append(instances, response.Response.Zones...)
		if len(response.Response.Zones) < int(limit) {
			break
		}

		offset += limit
	}

	if len(instances) < 1 {
		return
	}

	ret = instances[0]
	return
}

func (me *TeoService) DescribeTeoZoneSettingById(ctx context.Context, zoneId string) (ret *teo.ZoneSetting, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := teo.NewDescribeZoneSettingRequest()
	response := teo.NewDescribeZoneSettingResponse()
	request.ZoneId = &zoneId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := me.client.UseTeoClient().DescribeZoneSetting(request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		response = result
		return nil
	})
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response == nil {
		return
	}

	ret = response.Response.ZoneSetting
	return
}

func (me *TeoService) DescribeTeoRuleEngineById(ctx context.Context, zoneId string, ruleId string) (ret *teo.RuleItem, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := teo.NewDescribeRulesRequest()
	response := teo.NewDescribeRulesResponse()
	request.ZoneId = helper.String(zoneId)
	filter := &teo.Filter{
		Name:   helper.String("rule-id"),
		Values: []*string{helper.String(ruleId)},
	}
	request.Filters = append(request.Filters, filter)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := me.client.UseTeoClient().DescribeRules(request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		response = result
		return nil
	})
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response == nil || len(response.Response.RuleItems) < 1 {
		return
	}

	ret = response.Response.RuleItems[0]
	return
}

func (me *TeoService) DescribeTeoOriginGroupById(ctx context.Context, originGroupId string) (ret *teo.OriginGroup, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := teo.NewDescribeOriginGroupRequest()
	advancedFilter := &teo.AdvancedFilter{
		Name:   helper.String("origin-group-id"),
		Values: []*string{helper.String(originGroupId)},
	}
	request.Filters = append(request.Filters, advancedFilter)

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
	var instances []*teo.OriginGroup
	for {
		request.Offset = &offset
		request.Limit = &limit
		response := teo.NewDescribeOriginGroupResponse()
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			result, e := me.client.UseTeoClient().DescribeOriginGroup(request)
			if e != nil {
				return tccommon.RetryError(e)
			}
			response = result
			return nil
		})
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.OriginGroups) < 1 {
			break
		}
		instances = append(instances, response.Response.OriginGroups...)
		if len(response.Response.OriginGroups) < int(limit) {
			break
		}

		offset += limit
	}

	if len(instances) < 1 {
		return
	}

	ret = instances[0]
	return
}

func (me *TeoService) DescribeTeoCertificateConfigById(ctx context.Context, zoneId string, host string) (ret *teo.AccelerationDomain, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := teo.NewDescribeAccelerationDomainsRequest()
	response := teo.NewDescribeAccelerationDomainsResponse()
	request.ZoneId = &zoneId
	advancedFilter := &teo.AdvancedFilter{
		Name:   helper.String("domain-name"),
		Values: []*string{&host},
	}
	request.Filters = append(request.Filters, advancedFilter)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := me.client.UseTeoClient().DescribeAccelerationDomains(request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		response = result
		return nil
	})
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response == nil || len(response.Response.AccelerationDomains) < 1 {
		return
	}

	ret = response.Response.AccelerationDomains[0]
	return
}

func (me *TeoService) DescribeTeoL4ProxyById(ctx context.Context, zoneId string, proxyId string) (ret *teo.L4Proxy, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := teo.NewDescribeL4ProxyRequest()
	response := teo.NewDescribeL4ProxyResponse()
	request.ZoneId = &zoneId
	filter := &teo.Filter{
		Name:   helper.String("proxy-id"),
		Values: []*string{&proxyId},
	}
	request.Filters = append(request.Filters, filter)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := me.client.UseTeoClient().DescribeL4Proxy(request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		response = result
		return nil
	})
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response == nil || len(response.Response.L4Proxies) < 1 {
		return
	}

	ret = response.Response.L4Proxies[0]
	return
}

func (me *TeoService) DescribeTeoRealtimeLogDeliveryById(ctx context.Context, zoneId string, taskId string) (ret *teo.RealtimeLogDeliveryTask, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := teo.NewDescribeRealtimeLogDeliveryTasksRequest()
	response := teo.NewDescribeRealtimeLogDeliveryTasksResponse()
	request.ZoneId = helper.String(zoneId)
	advancedFilter := &teo.AdvancedFilter{
		Name:   helper.String("task-id"),
		Values: []*string{helper.String(taskId)},
	}
	request.Filters = append(request.Filters, advancedFilter)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := me.client.UseTeoClient().DescribeRealtimeLogDeliveryTasks(request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		response = result
		return nil
	})
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response == nil || len(response.Response.RealtimeLogDeliveryTasks) < 1 {
		return
	}

	ret = response.Response.RealtimeLogDeliveryTasks[0]
	return
}

func (me *TeoService) DescribeTeoSecurityIpGroupById(ctx context.Context, zoneId string, groupId string) (ret *teo.DescribeSecurityIPGroupResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := teo.NewDescribeSecurityIPGroupRequest()
	response := teo.NewDescribeSecurityIPGroupResponse()
	request.ZoneId = helper.String(zoneId)
	request.GroupIds = []*int64{helper.StrToInt64Point(groupId)}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := me.client.UseTeoV20220901Client().DescribeSecurityIPGroup(request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		response = result
		return nil
	})
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	ret = response.Response
	return
}

func (me *TeoService) DescribeTeoFunctionById(ctx context.Context, zoneId string, functionId string) (ret *teo.Function, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := teo.NewDescribeFunctionsRequest()
	response := teo.NewDescribeFunctionsResponse()
	request.ZoneId = helper.String(zoneId)
	request.FunctionIds = []*string{helper.String(functionId)}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := me.client.UseTeoV20220901Client().DescribeFunctions(request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		response = result
		return nil
	})
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response == nil || len(response.Response.Functions) < 1 {
		return
	}

	ret = response.Response.Functions[0]
	return
}

func (me *TeoService) DescribeTeoFunctionRuleById(ctx context.Context, zoneId string, functionId string, ruleId string) (ret *teo.FunctionRule, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := teo.NewDescribeFunctionRulesRequest()
	response := teo.NewDescribeFunctionRulesResponse()
	request.ZoneId = helper.String(zoneId)
	request.Filters = []*teo.Filter{{
		Name:   helper.String("function-id"),
		Values: []*string{helper.String(functionId)},
	}, {
		Name:   helper.String("rule-id"),
		Values: []*string{helper.String(ruleId)},
	}}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := me.client.UseTeoV20220901Client().DescribeFunctionRules(request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		response = result
		return nil
	})
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response == nil || len(response.Response.FunctionRules) < 1 {
		return
	}

	ret = response.Response.FunctionRules[0]
	return
}

func (me *TeoService) DescribeTeoFunctionRulePriorityById(ctx context.Context, zoneId string, functionId string) (ret *teo.DescribeFunctionRulesResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := teo.NewDescribeFunctionRulesRequest()
	response := teo.NewDescribeFunctionRulesResponse()
	request.ZoneId = helper.String(zoneId)
	filter := &teo.Filter{
		Name:   helper.String("function-id"),
		Values: []*string{helper.String(functionId)},
	}
	request.Filters = append(request.Filters, filter)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := me.client.UseTeoV20220901Client().DescribeFunctionRules(request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		response = result
		return nil
	})
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	ret = response.Response
	return
}

func (me *TeoService) DescribeTeoFunctionRuntimeEnvironmentById(ctx context.Context, zoneId string, functionId string) (ret *teo.DescribeFunctionRuntimeEnvironmentResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := teo.NewDescribeFunctionRuntimeEnvironmentRequest()
	response := teo.NewDescribeFunctionRuntimeEnvironmentResponse()
	request.ZoneId = helper.String(zoneId)
	request.FunctionId = helper.String(functionId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := me.client.UseTeoV20220901Client().DescribeFunctionRuntimeEnvironment(request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		response = result
		return nil
	})
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	ret = response.Response
	return
}

func (me *TeoService) DescribeTeoL7AccSettingById(ctx context.Context, zoneId string) (ret *teo.ZoneConfigParameters, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := teo.NewDescribeL7AccSettingRequest()
	response := teo.NewDescribeL7AccSettingResponse()
	request.ZoneId = helper.String(zoneId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := me.client.UseTeoV20220901Client().DescribeL7AccSetting(request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		response = result
		return nil
	})
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil {
		return
	}

	ret = response.Response.ZoneSetting
	return
}

func (me *TeoService) DescribeTeoL4ProxyRuleById(ctx context.Context, zoneId string, proxyId string, ruleId string) (ret *teov20220901.L4ProxyRule, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := teov20220901.NewDescribeL4ProxyRulesRequest()
	request.ZoneId = helper.String(zoneId)
	request.ProxyId = helper.String(proxyId)
	filter := &teo.Filter{
		Name:   helper.String("rule-id"),
		Values: []*string{helper.String(ruleId)},
	}
	request.Filters = append(request.Filters, filter)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	var (
		offset uint64 = 0
		limit  int64  = 20
	)
	var instances []*teov20220901.L4ProxyRule
	for {
		request.Offset = &offset
		request.Limit = &limit
		response := teo.NewDescribeL4ProxyRulesResponse()
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			result, e := me.client.UseTeoClient().DescribeL4ProxyRules(request)
			if e != nil {
				return tccommon.RetryError(e)
			}
			response = result
			return nil
		})
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.L4ProxyRules) < 1 {
			break
		}
		instances = append(instances, response.Response.L4ProxyRules...)
		if len(response.Response.L4ProxyRules) < int(limit) {
			break
		}

		offset = offset + uint64(limit)
	}

	if len(instances) < 1 {
		return
	}

	ret = instances[0]
	return
}

func (me *TeoService) DescribeTeoL7AccRuleById(ctx context.Context, zoneId string, ruleId string) (ret *teov20220901.DescribeL7AccRulesResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := teov20220901.NewDescribeL7AccRulesRequest()
	response := teov20220901.NewDescribeL7AccRulesResponse()
	request.ZoneId = helper.String(zoneId)

	if ruleId != "" {
		request.Filters = []*teov20220901.Filter{
			{
				Name:   helper.String("rule-id"),
				Values: helper.Strings([]string{ruleId}),
			},
		}
	}
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := me.client.UseTeoV20220901Client().DescribeL7AccRules(request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		response = result
		return nil
	})
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil {
		return
	}

	ret = response.Response
	return
}

func (me *TeoService) DescribeTeoSecurityPolicyConfigById(ctx context.Context, zoneId, entity, host, templateId string) (ret *teo.SecurityPolicy, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := teo.NewDescribeSecurityPolicyRequest()
	response := teo.NewDescribeSecurityPolicyResponse()
	request.ZoneId = &zoneId
	request.Entity = &entity
	if host != "" {
		request.Host = &host
	}

	if templateId != "" {
		request.TemplateId = &templateId
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseTeoV20220901Client().DescribeSecurityPolicy(request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	if response.Response == nil {
		return
	}

	ret = response.Response.SecurityPolicy
	return
}

func (me *TeoService) DescribeTeoZonesByFilter(ctx context.Context, param map[string]interface{}) (ret []*teov20220901.Zone, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = teov20220901.NewDescribeZonesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "Filters" {
			request.Filters = v.([]*teov20220901.AdvancedFilter)
		}
		if k == "Order" {
			request.Order = v.(*string)
		}
		if k == "Direction" {
			request.Direction = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 100
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response := teo.NewDescribeZonesResponse()
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			result, e := me.client.UseTeoClient().DescribeZones(request)
			if e != nil {
				return tccommon.RetryError(e)
			}
			response = result
			return nil
		})
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || response.Response == nil || len(response.Response.Zones) < 1 {
			break
		}
		ret = append(ret, response.Response.Zones...)
		if len(response.Response.Zones) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *TeoService) TeoL7AccRuleStateRefreshFunc(zoneId, taskId string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		request := teov20220901.NewDescribeZoneConfigImportResultRequest()
		request.ZoneId = helper.String(zoneId)
		request.TaskId = helper.String(taskId)
		ratelimit.Check(request.GetAction())
		object, err := me.client.UseTeoV20220901Client().DescribeZoneConfigImportResult(request)

		if err != nil {
			return nil, "", err
		}
		if object == nil || object.Response == nil || object.Response.Status == nil {
			return nil, "", nil
		}
		status := helper.PString(object.Response.Status)
		if len(failStates) > 0 {
			for _, state := range failStates {
				if strings.Contains(status, state) {
					return object, status, fmt.Errorf("teo[%s] sync check task[%s] failed, status is on [%s], return...", zoneId, taskId, status)
				}
			}
		}

		return object, status, nil
	}
}

func (me *TeoService) DescribeTeoDnsRecordById(ctx context.Context, zoneId, recordId string) (ret *teov20220901.DnsRecord, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := teov20220901.NewDescribeDnsRecordsRequest()
	request.ZoneId = helper.String(zoneId)
	request.Filters = []*teov20220901.AdvancedFilter{
		{
			Name:   helper.String("id"),
			Values: helper.Strings([]string{recordId}),
		},
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, e := me.client.UseTeoClient().DescribeDnsRecords(request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
		if len(response.Response.DnsRecords) > 0 {
			ret = response.Response.DnsRecords[0]
		}
		return nil
	})
	if err != nil {
		errRet = err
		return
	}
	return
}

func (me *TeoService) DescribeTeoBindSecurityTemplateById(ctx context.Context, zoneId string, templateId string, entity string) (ret *teov20220901.EntityStatus, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := teov20220901.NewDescribeSecurityTemplateBindingsRequest()
	request.ZoneId = helper.String(zoneId)
	request.TemplateId = []*string{helper.String(templateId)}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTeoV20220901Client().DescribeSecurityTemplateBindings(request)

	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	if response != nil && response.Response != nil {
		if response.Response.SecurityTemplate != nil && len(response.Response.SecurityTemplate) > 0 {
			if response.Response.SecurityTemplate[0] != nil && response.Response.SecurityTemplate[0].TemplateScope != nil && len(response.Response.SecurityTemplate[0].TemplateScope) > 0 {
				if response.Response.SecurityTemplate[0].TemplateScope[0] != nil && len(response.Response.SecurityTemplate[0].TemplateScope[0].EntityStatus) > 0 {
					for _, v := range response.Response.SecurityTemplate[0].TemplateScope[0].EntityStatus {
						if v != nil && *v.Entity == entity {
							ret = v
							return
						}
					}
				}
			}
		}
	}
	return
}

func (me *TeoService) DescribeTeoPlansById(ctx context.Context, planId string) (ret *teo.Plan, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := teo.NewDescribePlansRequest()
	response := teo.NewDescribePlansResponse()
	request.Filters = []*teo.Filter{
		{
			Name:   helper.String("plan-id"),
			Values: helper.Strings([]string{planId}),
		},
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseTeoV20220901Client().DescribePlans(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.Plans == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe plans failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	if len(response.Response.Plans) > 0 {
		ret = response.Response.Plans[0]
	}

	return
}

func (me *TeoService) DescribeTeoPlansByFilters(ctx context.Context, paramMap map[string]interface{}) (ret []*teo.Plan, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := teo.NewDescribePlansRequest()
	response := teo.NewDescribePlansResponse()

	for k, v := range paramMap {
		if k == "Filters" {
			request.Filters = v.([]*teov20220901.Filter)
		}

		if k == "Order" {
			request.Order = v.(*string)
		}

		if k == "Direction" {
			request.Direction = v.(*string)
		}
	}

	var (
		offset int64 = 0
		limit  int64 = 200
	)

	for {
		request.Offset = &offset
		request.Limit = &limit
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			result, e := me.client.UseTeoClient().DescribePlans(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Describe plans failed, Response is nil."))
			}

			response = result
			return nil
		})

		if err != nil {
			errRet = err
			return
		}

		if len(response.Response.Plans) < 1 {
			break
		}

		ret = append(ret, response.Response.Plans...)
		if len(response.Response.Plans) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *TeoService) DescribeTeoContentIdentifierById(ctx context.Context, contentId string) (ret *teo.ContentIdentifier, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := teo.NewDescribeContentIdentifiersRequest()
	response := teo.NewDescribeContentIdentifiersResponse()
	request.Filters = []*teo.AdvancedFilter{
		{
			Name:   helper.String("content-id"),
			Values: helper.Strings([]string{contentId}),
			Fuzzy:  helper.Bool(false),
		},
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseTeoV20220901Client().DescribeContentIdentifiers(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || len(result.Response.ContentIdentifiers) == 0 {
			return resource.NonRetryableError(fmt.Errorf("Describe teo content identifier failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	if len(response.Response.ContentIdentifiers) != 1 {
		errRet = fmt.Errorf("`ContentIdentifiers` returning multiple values, Should be one.")
		return
	}

	ret = response.Response.ContentIdentifiers[0]
	return
}

func (me *TeoService) DescribeTeoCustomizeErrorPageById(ctx context.Context, zoneId, pageId string) (ret *teo.CustomErrorPage, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := teo.NewDescribeCustomErrorPagesRequest()
	response := teo.NewDescribeCustomErrorPagesResponse()
	request.ZoneId = &zoneId
	request.Filters = []*teo.AdvancedFilter{
		{
			Name:   helper.String("page-id"),
			Values: helper.Strings([]string{pageId}),
			Fuzzy:  helper.Bool(false),
		},
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseTeoV20220901Client().DescribeCustomErrorPages(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || len(result.Response.ErrorPages) == 0 {
			return resource.NonRetryableError(fmt.Errorf("Describe teo custom error pages failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	if len(response.Response.ErrorPages) != 1 {
		errRet = fmt.Errorf("`ErrorPages` returning multiple values, Should be one.")
		return
	}

	ret = response.Response.ErrorPages[0]
	return
}

func (me *TeoService) WaitTeoOriginACLById(ctx context.Context, timeout time.Duration, zoneId, status string) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := teo.NewDescribeOriginACLRequest()
	request.ZoneId = &zoneId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	errRet = resource.Retry(timeout, func() *resource.RetryError {
		result, e := me.client.UseTeoClient().DescribeOriginACLWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.OriginACLInfo == nil || result.Response.OriginACLInfo.Status == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe teo origin acl failed, Response is nil."))
		}

		if *result.Response.OriginACLInfo.Status == status {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("TEO zone %s origin acl is still %s. Please contact TEO for assistance.", zoneId, *result.Response.OriginACLInfo.Status))
	})

	return
}

func (me *TeoService) DescribeTeoOriginACLById(ctx context.Context, zoneId string) (originACLInfo *teo.OriginACLInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := teo.NewDescribeOriginACLRequest()
	response := teo.NewDescribeOriginACLResponse()
	request.ZoneId = &zoneId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	errRet = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := me.client.UseTeoClient().DescribeOriginACLWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe teo origin acl failed, Response is nil."))
		}

		response = result
		return nil
	})

	if errRet != nil {
		log.Printf("[CRITAL]%s describe teo origin acl failed, reason:%+v", logId, errRet)
		return
	}

	originACLInfo = response.Response.OriginACLInfo
	return
}

func (me *TeoService) DescribeTeoDdosProtectionConfigById(ctx context.Context, zoneId string) (ret *teo.DDoSProtection, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := teo.NewDescribeDDoSProtectionRequest()
	response := teo.NewDescribeDDoSProtectionResponse()
	request.ZoneId = &zoneId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseTeoV20220901Client().DescribeDDoSProtection(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe teo ddos protection failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	ret = response.Response.DDoSProtection
	return
}

func (me *TeoService) DescribeTeoOriginAclByFilter(ctx context.Context, param map[string]interface{}) (ret *teo.DescribeOriginACLResponseParams, errRet error) {
	var (
		logId    = tccommon.GetLogId(ctx)
		request  = teo.NewDescribeOriginACLRequest()
		response = teo.NewDescribeOriginACLResponse()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "ZoneId" {
			request.ZoneId = v.(*string)
		}
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseTeoV20220901Client().DescribeOriginACL(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe teo origin acl failed, Response is nil."))
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

func (me *TeoService) DescribeTeoWebSecurityTemplateById(ctx context.Context, zoneId, templateId string) (ret *teov20220901.SecurityPolicy, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := teov20220901.NewDescribeWebSecurityTemplateRequest()
	request.ZoneId = &zoneId
	request.TemplateId = &templateId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTeoV20220901Client().DescribeWebSecurityTemplate(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response == nil {
		return
	}

	ret = response.Response.SecurityPolicy
	return
}

func (me *TeoService) DescribeTeoWebSecurityTemplateNameById(ctx context.Context, zoneId string, templateId string) (templateName string, errRet error) {
	var (
		logId    = tccommon.GetLogId(ctx)
		request  = teo.NewDescribeWebSecurityTemplatesRequest()
		response = teo.NewDescribeWebSecurityTemplatesResponse()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	request.ZoneIds = []*string{&zoneId}

	ratelimit.Check(request.GetAction())

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := me.client.UseTeoV20220901Client().DescribeWebSecurityTemplates(request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		response = result
		return nil
	})
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response == nil || response.Response.SecurityPolicyTemplates == nil {
		return
	}

	for _, template := range response.Response.SecurityPolicyTemplates {
		if template.TemplateId != nil && *template.TemplateId == templateId {
			if template.TemplateName != nil {
				templateName = *template.TemplateName
			}
			break
		}
	}

	return
}
