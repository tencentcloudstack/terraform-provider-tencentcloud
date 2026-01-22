package apm

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	apm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apm/v20210622"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type ApmService struct {
	client *connectivity.TencentCloudClient
}

func (me *ApmService) DescribeApmInstanceById(ctx context.Context, instanceId string) (instance *apm.ApmInstanceDetail, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := apm.NewDescribeApmInstancesRequest()
	response := apm.NewDescribeApmInstancesResponse()
	request.InstanceIds = []*string{&instanceId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, err := me.client.UseApmClient().DescribeApmInstances(request)
		if err != nil {
			return tccommon.RetryError(err)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.Instances == nil || len(result.Response.Instances) < 1 {
			return resource.NonRetryableError(fmt.Errorf("Describe apm instance failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	instance = response.Response.Instances[0]
	return
}

func (me *ApmService) DescribeApmAgentById(ctx context.Context, instanceId string) (apmAgent *apm.ApmAgentInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := apm.NewDescribeApmAgentRequest()
	response := apm.NewDescribeApmAgentResponse()
	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, err := me.client.UseApmClient().DescribeApmAgent(request)
		if err != nil {
			return tccommon.RetryError(err)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe apm agent failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	apmAgent = response.Response.ApmAgent
	return
}

func (me *ApmService) DeleteApmInstanceById(ctx context.Context, instanceId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := apm.NewTerminateApmInstanceRequest()
	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, err := me.client.UseApmClient().TerminateApmInstance(request)
		if err != nil {
			return tccommon.RetryError(err)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		return err
	}

	return
}

func (me *ApmService) DescribeApmSampleConfigById(ctx context.Context, instanceId, sampleName string) (ret *apm.ApmSampleConfig, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := apm.NewDescribeApmSampleConfigRequest()
	response := apm.NewDescribeApmSampleConfigResponse()
	request.InstanceId = &instanceId
	request.SampleName = &sampleName

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, err := me.client.UseApmClient().DescribeApmSampleConfig(request)
		if err != nil {
			return tccommon.RetryError(err)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.ApmSampleConfigs == nil || len(result.Response.ApmSampleConfigs) == 0 {
			return resource.NonRetryableError(fmt.Errorf("Describe apm sample config failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	ret = response.Response.ApmSampleConfigs[0]
	return
}

func (me *ApmService) DescribeApmApplicationConfigById(ctx context.Context, instanceId, serviceName string) (ret *apm.ApmAppConfig, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := apm.NewDescribeApmApplicationConfigRequest()
	response := apm.NewDescribeApmApplicationConfigResponse()
	request.InstanceId = &instanceId
	request.ServiceName = &serviceName

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, err := me.client.UseApmClient().DescribeApmApplicationConfig(request)
		if err != nil {
			return tccommon.RetryError(err)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.ApmAppConfig == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe apm application config failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	ret = response.Response.ApmAppConfig
	return
}

func (me *ApmService) DescribeApmAssociationById(ctx context.Context, instanceId, productName string) (ret *apm.ApmAssociation, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := apm.NewDescribeApmAssociationRequest()
	response := apm.NewDescribeApmAssociationResponse()
	request.InstanceId = &instanceId
	request.ProductName = &productName

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, err := me.client.UseApmClient().DescribeApmAssociation(request)
		if err != nil {
			return tccommon.RetryError(err)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.ApmAssociation == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe apm association failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	ret = response.Response.ApmAssociation
	return
}

func (me *ApmService) DescribeApmPrometheusRuleById(ctx context.Context, instanceId, ruleId string) (ret *apm.ApmPrometheusRules, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := apm.NewDescribeApmPrometheusRuleRequest()
	response := apm.NewDescribeApmPrometheusRuleResponse()
	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, err := me.client.UseApmClient().DescribeApmPrometheusRule(request)
		if err != nil {
			return tccommon.RetryError(err)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.ApmPrometheusRules == nil || len(result.Response.ApmPrometheusRules) == 0 {
			return resource.NonRetryableError(fmt.Errorf("Describe apm prometheus rule failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	for _, item := range response.Response.ApmPrometheusRules {
		if item != nil && item.Id != nil && helper.Int64ToStr(*item.Id) == ruleId {
			ret = item
			return
		}
	}

	return
}

func (me *ApmService) DescribeApmInstances(ctx context.Context, params map[string]interface{}) (instances []*apm.ApmInstanceDetail, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := apm.NewDescribeApmInstancesRequest()
	response := apm.NewDescribeApmInstancesResponse()

	// Set filter parameters
	if v, ok := params["instance_ids"]; ok {
		request.InstanceIds = helper.InterfacesStringsPoint(v.([]interface{}))
	}
	if v, ok := params["instance_id"]; ok {
		request.InstanceId = helper.String(v.(string))
	}
	if v, ok := params["instance_name"]; ok {
		request.InstanceName = helper.String(v.(string))
	}
	if v, ok := params["tags"]; ok {
		tags := v.(map[string]interface{})
		for key, value := range tags {
			tag := apm.ApmTag{
				Key:   helper.String(key),
				Value: helper.String(value.(string)),
			}
			request.Tags = append(request.Tags, &tag)
		}
	}
	if v, ok := params["demo_instance_flag"]; ok {
		request.DemoInstanceFlag = helper.IntInt64(v.(int))
	}
	if v, ok := params["all_regions_flag"]; ok {
		request.AllRegionsFlag = helper.IntInt64(v.(int))
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, err := me.client.UseApmClient().DescribeApmInstances(request)
		if err != nil {
			return tccommon.RetryError(err)
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

	if response.Response != nil && response.Response.Instances != nil {
		instances = response.Response.Instances
	}

	return
}
