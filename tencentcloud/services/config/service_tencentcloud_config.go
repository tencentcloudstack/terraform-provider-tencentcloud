package config

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	configv20220802 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/config/v20220802"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func NewConfigService(client *connectivity.TencentCloudClient) ConfigService {
	return ConfigService{client: client}
}

type ConfigService struct {
	client *connectivity.TencentCloudClient
}

func (me *ConfigService) DescribeConfigCompliancePacksByFilter(ctx context.Context, paramMap map[string]interface{}) (items []*configv20220802.ConfigCompliancePack, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := configv20220802.NewListCompliancePacksRequest()
	request.Limit = helper.Uint64(100)
	request.Offset = helper.Uint64(0)

	if v, ok := paramMap["CompliancePackName"]; ok {
		request.CompliancePackName = helper.String(v.(string))
	}

	if v, ok := paramMap["RiskLevel"]; ok {
		request.RiskLevel = v.([]*uint64)
	}

	if v, ok := paramMap["Status"]; ok {
		request.Status = helper.String(v.(string))
	}

	if v, ok := paramMap["ComplianceResult"]; ok {
		request.ComplianceResult = v.([]*string)
	}

	if v, ok := paramMap["OrderType"]; ok {
		request.OrderType = helper.String(v.(string))
	}

	items = make([]*configv20220802.ConfigCompliancePack, 0)

	for {
		var (
			total     uint64
			pageItems []*configv20220802.ConfigCompliancePack
		)

		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			result, e := me.client.UseConfigV20220802Client().ListCompliancePacksWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("list config compliance packs failed, Response is nil"))
			}

			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

			if result.Response.Total != nil {
				total = *result.Response.Total
			}

			pageItems = result.Response.Items
			return nil
		})

		if err != nil {
			errRet = err
			return
		}

		items = append(items, pageItems...)

		if uint64(len(items)) >= total {
			break
		}

		*request.Offset += *request.Limit
	}

	return
}

func (me *ConfigService) DescribeSystemConfigCompliancePacks(ctx context.Context) (items []*configv20220802.SystemCompliancePack, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := configv20220802.NewListSystemCompliancePacksRequest()
	request.Limit = helper.Int64(100)
	request.Offset = helper.Int64(0)

	items = make([]*configv20220802.SystemCompliancePack, 0)

	for {
		var (
			total     uint64
			pageItems []*configv20220802.SystemCompliancePack
		)

		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			result, e := me.client.UseConfigV20220802Client().ListSystemCompliancePacksWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("list system config compliance packs failed, Response is nil"))
			}

			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

			if result.Response.Total != nil {
				total = *result.Response.Total
			}

			pageItems = result.Response.Items
			return nil
		})

		if err != nil {
			errRet = err
			return
		}

		items = append(items, pageItems...)

		if uint64(len(items)) >= total {
			break
		}

		*request.Offset += *request.Limit
	}

	return
}

func (me *ConfigService) DescribeConfigSystemRulesByFilter(ctx context.Context, paramMap map[string]interface{}) (items []*configv20220802.SystemConfigRule, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := configv20220802.NewListSystemRulesRequest()
	request.Limit = helper.Uint64(100)
	request.Offset = helper.Uint64(0)

	if v, ok := paramMap["Keyword"]; ok {
		request.Keyword = helper.String(v.(string))
	}

	if v, ok := paramMap["RiskLevel"]; ok {
		request.RiskLevel = helper.Uint64(uint64(v.(int)))
	}

	items = make([]*configv20220802.SystemConfigRule, 0)

	for {
		var (
			total     uint64
			pageItems []*configv20220802.SystemConfigRule
		)

		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			result, e := me.client.UseConfigV20220802Client().ListSystemRulesWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("list config system rules failed, Response is nil"))
			}

			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

			if result.Response.Total != nil {
				total = *result.Response.Total
			}

			pageItems = result.Response.Items
			return nil
		})

		if err != nil {
			errRet = err
			return
		}

		items = append(items, pageItems...)

		if uint64(len(items)) >= total {
			break
		}

		*request.Offset += *request.Limit
	}

	return
}

func (me *ConfigService) DescribeConfigRuleEvaluationResultsByFilter(ctx context.Context, paramMap map[string]interface{}) (items []*configv20220802.EvaluationResult, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := configv20220802.NewListConfigRuleEvaluationResultsRequest()
	request.Limit = helper.Uint64(100)
	request.Offset = helper.Uint64(0)

	if v, ok := paramMap["ConfigRuleId"]; ok {
		request.ConfigRuleId = helper.String(v.(string))
	}

	if v, ok := paramMap["ResourceType"]; ok {
		request.ResourceType = v.([]*string)
	}

	if v, ok := paramMap["ComplianceType"]; ok {
		request.ComplianceType = v.([]*string)
	}

	items = make([]*configv20220802.EvaluationResult, 0)

	for {
		var (
			total     uint64
			pageItems []*configv20220802.EvaluationResult
		)

		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			result, e := me.client.UseConfigV20220802Client().ListConfigRuleEvaluationResultsWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("list config rule evaluation results failed, Response is nil"))
			}

			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

			if result.Response.Total != nil {
				total = *result.Response.Total
			}

			pageItems = result.Response.Items
			return nil
		})

		if err != nil {
			errRet = err
			return
		}

		items = append(items, pageItems...)

		if uint64(len(items)) >= total {
			break
		}

		*request.Offset += *request.Limit
	}

	return
}

func (me *ConfigService) DescribeConfigRulesByFilter(ctx context.Context, paramMap map[string]interface{}) (items []*configv20220802.ConfigRule, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := configv20220802.NewListConfigRulesRequest()
	request.Limit = helper.Int64(200)
	request.Offset = helper.Int64(0)

	if v, ok := paramMap["RuleName"]; ok {
		request.RuleName = helper.String(v.(string))
	}

	if v, ok := paramMap["RiskLevel"]; ok {
		request.RiskLevel = v.([]*uint64)
	}

	if v, ok := paramMap["State"]; ok {
		request.State = helper.String(v.(string))
	}

	if v, ok := paramMap["ComplianceResult"]; ok {
		request.ComplianceResult = v.([]*string)
	}

	if v, ok := paramMap["OrderType"]; ok {
		request.OrderType = helper.String(v.(string))
	}

	items = make([]*configv20220802.ConfigRule, 0)

	for {
		var (
			total     uint64
			pageItems []*configv20220802.ConfigRule
		)

		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			result, e := me.client.UseConfigV20220802Client().ListConfigRulesWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("list config rules failed, Response is nil"))
			}

			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

			if result.Response.Total != nil {
				total = *result.Response.Total
			}

			pageItems = result.Response.Items
			return nil
		})

		if err != nil {
			errRet = err
			return
		}

		items = append(items, pageItems...)

		if uint64(len(items)) >= total {
			break
		}

		*request.Offset += *request.Limit
	}

	return
}

func (me *ConfigService) DescribeConfigRuleById(ctx context.Context, ruleId string) (ret *configv20220802.ConfigRule, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := configv20220802.NewDescribeConfigRuleRequest()
	response := configv20220802.NewDescribeConfigRuleResponse()
	request.RuleId = &ruleId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseConfigV20220802Client().DescribeConfigRule(request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("describe config rule failed, Response is nil"))
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	ret = response.Response.ConfigRule
	return
}

func (me *ConfigService) DescribeConfigAlarmPolicyById(ctx context.Context, alarmPolicyId uint64) (ret *configv20220802.AlarmPolicyRsp, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := configv20220802.NewListAlarmPolicyRequest()
	request.Offset = helper.Uint64(0)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for {
		var result *configv20220802.ListAlarmPolicyResponse

		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			resp, e := me.client.UseConfigV20220802Client().ListAlarmPolicy(request)
			if e != nil {
				return tccommon.RetryError(e)
			}

			if resp == nil || resp.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("list config alarm policy failed, Response is nil"))
			}

			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), resp.ToJsonString())

			result = resp
			return nil
		})

		if err != nil {
			errRet = err
			return
		}

		for _, policy := range result.Response.AlarmPolicyList {
			if policy.AlarmPolicyId != nil && *policy.AlarmPolicyId == alarmPolicyId {
				ret = policy
				return
			}
		}

		if len(result.Response.AlarmPolicyList) == 0 {
			break
		}

		*request.Offset++
	}

	return
}

func (me *ConfigService) DescribeConfigRecorder(ctx context.Context) (ret *configv20220802.DescribeConfigRecorderResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := configv20220802.NewDescribeConfigRecorderRequest()
	response := configv20220802.NewDescribeConfigRecorderResponse()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseConfigV20220802Client().DescribeConfigRecorder(request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("describe config recorder failed, Response is nil"))
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

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

func (me *ConfigService) DescribeConfigDiscoveredResourcesByFilter(ctx context.Context, paramMap map[string]interface{}) (items []*configv20220802.ResourceListInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := configv20220802.NewListDiscoveredResourcesRequest()
	request.MaxResults = helper.Uint64(200)

	if v, ok := paramMap["Filters"]; ok {
		request.Filters = v.([]*configv20220802.Filter)
	}

	if v, ok := paramMap["Tags"]; ok {
		request.Tags = v.([]*configv20220802.Tag)
	}

	if v, ok := paramMap["OrderType"]; ok {
		request.OrderType = helper.String(v.(string))
	}

	items = make([]*configv20220802.ResourceListInfo, 0)

	for {
		var (
			nextToken *string
			pageItems []*configv20220802.ResourceListInfo
		)

		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			result, e := me.client.UseConfigV20220802Client().ListDiscoveredResourcesWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("list config discovered resources failed, Response is nil"))
			}

			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

			nextToken = result.Response.NextToken
			pageItems = result.Response.Items
			return nil
		})

		if err != nil {
			errRet = err
			return
		}

		items = append(items, pageItems...)

		if nextToken == nil || *nextToken == "" {
			break
		}

		request.NextToken = nextToken
	}

	return
}

func (me *ConfigService) DescribeConfigResourceTypes(ctx context.Context) (items []*configv20220802.ConfigResource, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := configv20220802.NewListResourceTypesRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseConfigV20220802Client().ListResourceTypesWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("list config resource types failed, Response is nil"))
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

		items = result.Response.Items
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	return
}

func (me *ConfigService) DescribeConfigDeliver(ctx context.Context) (ret *configv20220802.DescribeConfigDeliverResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := configv20220802.NewDescribeConfigDeliverRequest()
	response := configv20220802.NewDescribeConfigDeliverResponse()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseConfigV20220802Client().DescribeConfigDeliver(request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("describe config deliver failed, Response is nil"))
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

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

func (me *ConfigService) DescribeConfigRemediationById(ctx context.Context, ruleId, remediationId string) (ret *configv20220802.Remediation, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := configv20220802.NewListRemediationsRequest()
	request.Limit = helper.Uint64(200)
	if ruleId != "" {
		request.RuleIds = []*string{&ruleId}
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseConfigV20220802Client().ListRemediations(request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("list config remediations failed, Response is nil"))
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

		for _, item := range result.Response.Remediations {
			if item.RemediationId != nil && *item.RemediationId == remediationId {
				ret = item
				return nil
			}
		}

		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	return
}

func (me *ConfigService) DescribeConfigCompliancePackById(ctx context.Context, compliancePackId string) (ret *configv20220802.DescribeCompliancePackResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := configv20220802.NewDescribeCompliancePackRequest()
	response := configv20220802.NewDescribeCompliancePackResponse()
	request.CompliancePackId = &compliancePackId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseConfigV20220802Client().DescribeCompliancePack(request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("describe config compliance pack failed, Response is nil"))
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
