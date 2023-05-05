package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	apigateway "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apigateway/v20180808"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type APIGatewayService struct {
	client *connectivity.TencentCloudClient
}

func (me *APIGatewayService) CreateApiKey(ctx context.Context, secretName string) (accessKeyId string, errRet error) {
	request := apigateway.NewCreateApiKeyRequest()
	request.SecretName = &secretName
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseAPIGatewayClient().CreateApiKey(request)
	if err != nil {
		errRet = err
		return
	}
	if response.Response.Result == nil || response.Response.Result.AccessKeyId == nil {
		errRet = fmt.Errorf("TencentCloud SDK %s return empty AccessKeyId", request.GetAction())
		return
	}
	accessKeyId = *response.Response.Result.AccessKeyId
	return
}

func (me *APIGatewayService) EnableApiKey(ctx context.Context, accessKeyId string) (errRet error) {
	request := apigateway.NewEnableApiKeyRequest()
	request.AccessKeyId = &accessKeyId
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseAPIGatewayClient().EnableApiKey(request)
	if err != nil {
		errRet = err
		return
	}
	if response.Response.Result == nil {
		errRet = fmt.Errorf("TencentCloud SDK %s return empty response", request.GetAction())
		return
	}
	if !*response.Response.Result {
		errRet = fmt.Errorf("enable API key fail")
		return
	}
	return
}

func (me *APIGatewayService) DisableApiKey(ctx context.Context, accessKeyId string) (errRet error) {
	request := apigateway.NewDisableApiKeyRequest()
	request.AccessKeyId = &accessKeyId
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseAPIGatewayClient().DisableApiKey(request)
	if err != nil {
		errRet = err
		return
	}
	if response.Response.Result == nil {
		errRet = fmt.Errorf("TencentCloud SDK %s return empty response", request.GetAction())
		return
	}
	if !*response.Response.Result {
		errRet = fmt.Errorf("disable API key fail")
		return
	}
	return
}

func (me *APIGatewayService) DescribeApiKey(ctx context.Context,
	accessKeyId string) (apiKey *apigateway.ApiKey, has bool, errRet error) {
	apiKeySet, err := me.DescribeApiKeysStatus(ctx, "", accessKeyId)
	if err != nil {
		errRet = err
		return
	}
	if len(apiKeySet) == 0 {
		return
	}
	has = true
	apiKey = apiKeySet[0]
	return
}

func (me *APIGatewayService) DescribeApiKeysStatus(ctx context.Context, secretName, accessKeyId string) (apiKeySet []*apigateway.ApiKey, errRet error) {
	request := apigateway.NewDescribeApiKeysStatusRequest()
	if secretName != "" || accessKeyId != "" {
		request.Filters = make([]*apigateway.Filter, 0, 2)
		if secretName != "" {
			request.Filters = append(request.Filters, &apigateway.Filter{Name: helper.String("SecretName"),
				Values: []*string{
					&secretName,
				}})
		}
		if accessKeyId != "" {
			request.Filters = append(request.Filters, &apigateway.Filter{Name: helper.String("AccessKeyId"),
				Values: []*string{
					&accessKeyId,
				}})
		}
	}

	var limit int64 = 100
	var offset int64 = 0

	request.Limit = &limit
	request.Offset = &offset

	for {
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseAPIGatewayClient().DescribeApiKeysStatus(request)
		if err != nil {
			errRet = err
			return
		}
		if response.Response.Result == nil {
			errRet = fmt.Errorf("TencentCloud SDK %s return empty response", request.GetAction())
			return
		}
		if len(response.Response.Result.ApiKeySet) > 0 {
			apiKeySet = append(apiKeySet, response.Response.Result.ApiKeySet...)
		}
		if len(response.Response.Result.ApiKeySet) < int(limit) {
			return
		}
		offset += limit
	}
}

func (me *APIGatewayService) DeleteApiKey(ctx context.Context, accessKeyId string) (errRet error) {
	request := apigateway.NewDeleteApiKeyRequest()
	request.AccessKeyId = &accessKeyId
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseAPIGatewayClient().DeleteApiKey(request)
	if err != nil {
		errRet = err
		return
	}
	if response.Response.Result == nil {
		errRet = fmt.Errorf("TencentCloud SDK %s return empty response", request.GetAction())
		return
	}
	if !*response.Response.Result {
		errRet = fmt.Errorf("delete API key fail")
		return
	}
	return
}

func (me *APIGatewayService) CreateUsagePlan(ctx context.Context, usagePlanName string, usagePlanDesc *string,
	maxRequestNum, maxRequestNumPreSec int64) (usagePlanId string, errRet error) {

	logId := getLogId(ctx)

	request := apigateway.NewCreateUsagePlanRequest()
	request.UsagePlanName = &usagePlanName
	request.MaxRequestNum = &maxRequestNum
	request.MaxRequestNumPreSec = &maxRequestNumPreSec
	if nil != usagePlanDesc {
		request.UsagePlanDesc = usagePlanDesc
	}

	errRet = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())

		response, err := me.client.UseAPIGatewayClient().CreateUsagePlan(request)
		if err != nil {
			log.Printf("[CRITAL]%s API[%s] fail, reason:%s", logId, request.GetAction(), err.Error())
			return retryError(err)
		}
		if response.Response.Result == nil {
			return resource.NonRetryableError(fmt.Errorf("TencentCloud SDK %s return empty response", request.GetAction()))
		}
		usagePlanId = *response.Response.Result.UsagePlanId
		return nil
	})

	return
}

func (me *APIGatewayService) DescribeUsagePlan(ctx context.Context, usagePlanId string) (info apigateway.UsagePlanInfo, has bool, errRet error) {
	request := apigateway.NewDescribeUsagePlanRequest()
	request.UsagePlanId = &usagePlanId

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseAPIGatewayClient().DescribeUsagePlan(request)
	if err != nil {
		if sdkErr, ok := err.(*errors.TencentCloudSDKError); ok && sdkErr.GetCode() == "ResourceNotFound.InvalidUsagePlan" {
			return
		}
		errRet = err
		return
	}
	if response.Response.Result == nil {
		errRet = fmt.Errorf("TencentCloud SDK %s return empty response", request.GetAction())
		return
	}
	has = true
	info = *response.Response.Result
	return
}

func (me *APIGatewayService) DeleteUsagePlan(ctx context.Context, usagePlanId string) (errRet error) {
	request := apigateway.NewDeleteUsagePlanRequest()
	request.UsagePlanId = &usagePlanId

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseAPIGatewayClient().DeleteUsagePlan(request)

	if err != nil {
		return err
	}
	if response.Response.Result == nil {
		return fmt.Errorf("TencentCloud SDK %s return empty response", request.GetAction())
	}

	if !*response.Response.Result {
		return fmt.Errorf("delete usage plan fail")
	}

	return
}

func (me *APIGatewayService) ModifyUsagePlan(ctx context.Context,
	usagePlanId string,
	usagePlanName string,
	usagePlanDesc *string,
	maxRequestNum,
	maxRequestNumPreSec int64) (errRet error) {

	request := apigateway.NewModifyUsagePlanRequest()
	request.UsagePlanId = &usagePlanId

	ratelimit.Check(request.GetAction())
	request.UsagePlanName = &usagePlanName
	if usagePlanDesc != nil {
		request.UsagePlanDesc = usagePlanDesc
	}
	request.MaxRequestNum = &maxRequestNum
	request.MaxRequestNumPreSec = &maxRequestNumPreSec

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseAPIGatewayClient().ModifyUsagePlan(request)
	if err != nil {
		errRet = err
		return
	}
	if response.Response.Result == nil {
		errRet = fmt.Errorf("TencentCloud SDK %s return empty response", request.GetAction())
		return
	}

	return nil
}

func (me *APIGatewayService) DescribeUsagePlanEnvironments(ctx context.Context,
	usagePlanId string, bindType string) (list []*apigateway.UsagePlanEnvironment, errRet error) {

	request := apigateway.NewDescribeUsagePlanEnvironmentsRequest()
	request.UsagePlanId = &usagePlanId
	request.BindType = &bindType

	var limit int64 = 100
	var offset int64 = 0

	request.Limit = &limit
	request.Offset = &offset

	for {
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseAPIGatewayClient().DescribeUsagePlanEnvironments(request)
		if err != nil {
			errRet = err
			return
		}
		if response.Response.Result == nil {
			errRet = fmt.Errorf("TencentCloud SDK %s return empty response", request.GetAction())
			return
		}
		if len(response.Response.Result.EnvironmentList) > 0 {
			list = append(list, response.Response.Result.EnvironmentList...)
		}
		if len(response.Response.Result.EnvironmentList) < int(limit) {
			return
		}
		offset += limit
	}
}

func (me *APIGatewayService) DescribeUsagePlansStatus(ctx context.Context,
	usagePlanId string, usagePlanName string) (infos []*apigateway.UsagePlanStatusInfo, errRet error) {

	request := apigateway.NewDescribeUsagePlansStatusRequest()

	request.Filters = make([]*apigateway.Filter, 0, 2)
	if usagePlanId != "" {
		request.Filters = append(request.Filters, &apigateway.Filter{Name: helper.String("UsagePlanId"),
			Values: []*string{
				&usagePlanId,
			}})
	}
	if usagePlanName != "" {
		request.Filters = append(request.Filters, &apigateway.Filter{Name: helper.String("UsagePlanName"),
			Values: []*string{
				&usagePlanName,
			}})
	}

	var limit int64 = 100
	var offset int64 = 0

	request.Limit = &limit
	request.Offset = &offset

	for {
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseAPIGatewayClient().DescribeUsagePlansStatus(request)
		if err != nil {
			errRet = err
			return
		}
		if response.Response.Result == nil {
			errRet = fmt.Errorf("TencentCloud SDK %s return empty response", request.GetAction())
			return
		}
		if len(response.Response.Result.UsagePlanStatusSet) > 0 {
			infos = append(infos, response.Response.Result.UsagePlanStatusSet...)
		}
		if len(response.Response.Result.UsagePlanStatusSet) < int(limit) {
			return
		}
		offset += limit
	}
}

func (me *APIGatewayService) DescribeIPStrategysStatus(ctx context.Context,
	serviceId, strategyName string) (infos []*apigateway.IPStrategy, errRet error) {

	request := apigateway.NewDescribeIPStrategysStatusRequest()
	request.ServiceId = &serviceId

	if strategyName != "" {
		request.Filters = make([]*apigateway.Filter, 0, 1)
		if strategyName != "" {
			request.Filters = append(request.Filters, &apigateway.Filter{Name: helper.String("StrategyName"),
				Values: []*string{
					&strategyName,
				}})
		}
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseAPIGatewayClient().DescribeIPStrategysStatus(request)
	if err != nil {
		errRet = err
		return
	}
	if response.Response.Result == nil {
		errRet = fmt.Errorf("TencentCloud SDK %s return empty response", request.GetAction())
		return
	}
	if len(response.Response.Result.StrategySet) > 0 {
		infos = append(infos, response.Response.Result.StrategySet...)
	}
	return
}

func (me *APIGatewayService) DescribeIPStrategies(ctx context.Context, serviceId, strategyId, environmentName string) (ipStrategies *apigateway.IPStrategy, errRet error) {
	request := apigateway.NewDescribeIPStrategyRequest()

	request.ServiceId = &serviceId
	request.StrategyId = &strategyId
	request.EnvironmentName = &environmentName

	var (
		limit   int64 = 100
		offset  int64 = 0
		apiList []*apigateway.DesApisStatus
	)

	request.Limit = &limit
	request.Offset = &offset

	for {
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseAPIGatewayClient().DescribeIPStrategy(request)
		if err != nil {
			errRet = err
			return
		}
		if response.Response.Result == nil {
			errRet = fmt.Errorf("TencentCloud SDK %s return empty response", request.GetAction())
			return
		}
		if len(response.Response.Result.BindApis) > 0 {
			apiList = append(apiList, response.Response.Result.BindApis...)
		}
		if len(response.Response.Result.BindApis) < int(limit) {
			ipStrategies = response.Response.Result
			ipStrategies.BindApis = apiList
			return
		}
		offset += limit
	}
}

func (me *APIGatewayService) DescribeServiceSubDomains(ctx context.Context, serviceId string) (domainList []*apigateway.DomainSetList, errRet error) {
	request := apigateway.NewDescribeServiceSubDomainsRequest()
	request.ServiceId = &serviceId

	var (
		limit  int64 = 100
		offset int64 = 0
	)

	request.Limit = &limit
	request.Offset = &offset

	for {
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseAPIGatewayClient().DescribeServiceSubDomains(request)
		if err != nil {
			errRet = err
			return
		}
		if response.Response.Result == nil {
			errRet = fmt.Errorf("TencentCloud SDK %s return empty response", request.GetAction())
			return
		}
		if len(response.Response.Result.DomainSet) > 0 {
			domainList = append(domainList, response.Response.Result.DomainSet...)
		}
		if len(response.Response.Result.DomainSet) < int(limit) {
			return
		}
		offset += limit
	}
}

func (me *APIGatewayService) BindSecretId(ctx context.Context,
	usagePlanId string, apiKeyId string) (errRet error) {

	request := apigateway.NewBindSecretIdsRequest()
	request.UsagePlanId = &usagePlanId
	request.AccessKeyIds = []*string{&apiKeyId}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseAPIGatewayClient().BindSecretIds(request)

	if err != nil {
		return err
	}
	if response.Response.Result == nil {
		return fmt.Errorf("TencentCloud SDK %s return empty response", request.GetAction())
	}

	if !*response.Response.Result {
		return fmt.Errorf("bind API key to usage plan fail")
	}

	return
}

func flattenOauthConfigMappings(v *apigateway.OauthConfig) map[string]interface{} {
	if v != nil {
		return map[string]interface{}{
			"login_redirect_url": *v.LoginRedirectUrl,
			"public_key":         *v.PublicKey,
			"token_location":     *v.TokenLocation,
		}
	}
	return nil
}

func (me *APIGatewayService) UnBindSecretId(ctx context.Context,
	usagePlanId string,
	apiKeyId string) (errRet error) {
	request := apigateway.NewUnBindSecretIdsRequest()
	request.UsagePlanId = &usagePlanId
	request.AccessKeyIds = []*string{&apiKeyId}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseAPIGatewayClient().UnBindSecretIds(request)

	if err != nil {
		return err
	}
	if response.Response.Result == nil {
		return fmt.Errorf("TencentCloud SDK %s return empty response", request.GetAction())
	}

	if !*response.Response.Result {
		return fmt.Errorf("unbind API key to usage plan fail")
	}

	return
}

func (me *APIGatewayService) CreateService(ctx context.Context,
	serviceName,
	protocol,
	serviceDesc,
	exclusiveSetName,
	ipVersion,
	setServerName,
	appidType string,
	netTypes []string) (serviceId string, errRet error) {

	request := apigateway.NewCreateServiceRequest()
	request.ServiceName = &serviceName
	request.Protocol = &protocol
	if serviceDesc != "" {
		request.ServiceDesc = &serviceDesc
	}
	if exclusiveSetName != "" {
		request.ExclusiveSetName = &exclusiveSetName
	}
	if ipVersion != "" {
		request.IpVersion = &ipVersion
	}
	if appidType != "" {
		request.AppIdType = &appidType
	}
	if setServerName != "" {
		request.SetServerName = &setServerName
	}
	request.NetTypes = helper.Strings(netTypes)

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseAPIGatewayClient().CreateService(request)

	if err != nil {
		errRet = err
		return
	}
	serviceId = *response.Response.ServiceId
	return
}

func (me *APIGatewayService) DescribeService(ctx context.Context, serviceId string) (info apigateway.DescribeServiceResponse, has bool, errRet error) {
	request := apigateway.NewDescribeServiceRequest()
	request.ServiceId = &serviceId

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseAPIGatewayClient().DescribeService(request)
	if err != nil {
		if sdkError, ok := err.(*errors.TencentCloudSDKError); ok && sdkError.Code == SERVICE_ERR_CODE {
			return
		}
		errRet = err
		return
	}
	info = *response
	has = true
	return
}

func (me *APIGatewayService) ModifyService(ctx context.Context,
	serviceId,
	serviceName,
	protocol,
	serviceDesc string,
	netTypes []string) (errRet error) {

	request := apigateway.NewModifyServiceRequest()
	request.ServiceId = &serviceId
	request.ServiceName = &serviceName
	request.Protocol = &protocol
	request.ServiceDesc = &serviceDesc
	request.NetTypes = helper.Strings(netTypes)

	ratelimit.Check(request.GetAction())
	_, err := me.client.UseAPIGatewayClient().ModifyService(request)
	if err != nil {
		errRet = err
		return
	}
	return
}

func (me *APIGatewayService) DeleteService(ctx context.Context,
	serviceId string) (errRet error) {

	request := apigateway.NewDeleteServiceRequest()
	request.ServiceId = &serviceId

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseAPIGatewayClient().DeleteService(request)
	if err != nil {
		errRet = err
		return
	}
	if response.Response.Result == nil {
		return fmt.Errorf("TencentCloud SDK %s return empty response", request.GetAction())
	}

	if !*response.Response.Result {
		return fmt.Errorf("delete service fail")
	}

	return
}

func (me *APIGatewayService) UnReleaseService(ctx context.Context, serviceId, environment string) (errRet error) {
	request := apigateway.NewUnReleaseServiceRequest()
	request.ServiceId = &serviceId
	request.EnvironmentName = &environment

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseAPIGatewayClient().UnReleaseService(request)
	if err != nil {
		errRet = err
		return
	}
	if response.Response.Result == nil {
		return fmt.Errorf("TencentCloud SDK %s return empty response", request.GetAction())
	}

	if !*response.Response.Result {
		return fmt.Errorf("unrelease service %s.%s fail", serviceId, environment)
	}
	return
}

func (me *APIGatewayService) DescribeServiceUsagePlan(ctx context.Context,
	serviceId string) (list []*apigateway.ApiUsagePlan, errRet error) {

	request := apigateway.NewDescribeServiceUsagePlanRequest()
	request.ServiceId = &serviceId

	var limit int64 = 100
	var offset int64 = 0

	request.Limit = &limit
	request.Offset = &offset

	for {
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseAPIGatewayClient().DescribeServiceUsagePlan(request)
		if err != nil {
			errRet = err
			return
		}
		if response.Response.Result == nil {
			errRet = fmt.Errorf("TencentCloud SDK %s return empty response", request.GetAction())
			return
		}
		if len(response.Response.Result.ServiceUsagePlanList) > 0 {
			list = append(list, response.Response.Result.ServiceUsagePlanList...)
		}
		if len(response.Response.Result.ServiceUsagePlanList) < int(limit) {
			return
		}
		offset += limit
	}
}

func (me *APIGatewayService) DescribeApiUsagePlan(ctx context.Context,
	serviceId string) (list []*apigateway.ApiUsagePlan, errRet error) {

	request := apigateway.NewDescribeApiUsagePlanRequest()
	request.ServiceId = &serviceId

	var limit int64 = 100
	var offset int64 = 0

	request.Limit = &limit
	request.Offset = &offset

	for {
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseAPIGatewayClient().DescribeApiUsagePlan(request)
		if err != nil {
			errRet = err
			return
		}
		if response.Response.Result == nil {
			errRet = fmt.Errorf("TencentCloud SDK %s return empty response", request.GetAction())
			return
		}
		if len(response.Response.Result.ApiUsagePlanList) > 0 {
			list = append(list, response.Response.Result.ApiUsagePlanList...)
		}
		if len(response.Response.Result.ApiUsagePlanList) < int(limit) {
			return
		}
		offset += limit
	}
}

func (me *APIGatewayService) BindEnvironment(ctx context.Context,
	serviceId, environment, bindType, usagePlanId, apiId string) (errRet error) {

	logId := getLogId(ctx)
	request := apigateway.NewBindEnvironmentRequest()
	request.ServiceId = &serviceId
	request.UsagePlanIds = []*string{&usagePlanId}
	request.Environment = &environment
	request.BindType = &bindType

	if bindType == API_GATEWAY_TYPE_API {
		request.ApiIds = []*string{&apiId}
	}

	errRet = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())

		response, err := me.client.UseAPIGatewayClient().BindEnvironment(request)
		if err != nil {
			log.Printf("[CRITAL]%s API[%s] fail, reason:%s", logId, request.GetAction(), err.Error())
			return retryError(err)
		}

		if response.Response.Result == nil {
			return resource.NonRetryableError(fmt.Errorf("TencentCloud SDK %s return empty response", request.GetAction()))
		}

		if !*response.Response.Result {
			return resource.RetryableError(fmt.Errorf("%s attach to %s.%s fail", usagePlanId, serviceId, apiId))
		}

		return nil
	})

	return
}

func (me *APIGatewayService) UnBindEnvironment(ctx context.Context,
	serviceId, environment, bindType, usagePlanId, apiId string) (errRet error) {

	request := apigateway.NewUnBindEnvironmentRequest()
	request.ServiceId = &serviceId
	request.UsagePlanIds = []*string{&usagePlanId}
	request.Environment = &environment
	request.BindType = &bindType

	if bindType == API_GATEWAY_TYPE_API {
		request.ApiIds = []*string{&apiId}
	}

	errRet = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())

		response, errRet := me.client.UseAPIGatewayClient().UnBindEnvironment(request)
		if errRet != nil {
			return retryError(errRet)
		}

		if response.Response.Result == nil {
			return resource.NonRetryableError(fmt.Errorf("TencentCloud SDK %s return empty response", request.GetAction()))
		}

		if !*response.Response.Result {
			return resource.RetryableError(fmt.Errorf("%s unattach to %s.%s fail", usagePlanId, serviceId, apiId))
		}

		return nil
	})

	return
}

func (me *APIGatewayService) DescribeApi(ctx context.Context,
	serviceId,
	apiId string) (info apigateway.ApiInfo, has bool, errRet error) {

	request := apigateway.NewDescribeApiRequest()
	request.ServiceId = &serviceId
	request.ApiId = &apiId

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseAPIGatewayClient().DescribeApi(request)
	if err != nil {
		if sdkError, ok := err.(*errors.TencentCloudSDKError); ok && sdkError.Code == SERVICE_ERR_CODE || sdkError.Code == API_ERR_CODE {

			return
		}
		errRet = err
		return
	}

	if response.Response.Result == nil {
		errRet = fmt.Errorf("TencentCloud SDK %s return empty response", request.GetAction())
		return
	}

	has = true
	info = *response.Response.Result
	return
}

func (me *APIGatewayService) DeleteApi(ctx context.Context, serviceId,
	apiId string) (errRet error) {
	request := apigateway.NewDeleteApiRequest()
	request.ServiceId = &serviceId
	request.ApiId = &apiId
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseAPIGatewayClient().DeleteApi(request)
	if err != nil {
		errRet = err
		return
	}
	if response.Response.Result == nil {
		errRet = fmt.Errorf("TencentCloud SDK %s return empty response", request.GetAction())
		return
	}
	if !*response.Response.Result {
		errRet = fmt.Errorf("delete API fail")
		return
	}
	return
}

func (me *APIGatewayService) DescribeServicesStatus(ctx context.Context,
	serviceId,
	serviceName string) (infos []*apigateway.Service, errRet error) {

	request := apigateway.NewDescribeServicesStatusRequest()

	if serviceId != "" || serviceName != "" {
		request.Filters = make([]*apigateway.Filter, 0, 2)
		if serviceId != "" {
			request.Filters = append(request.Filters, &apigateway.Filter{Name: helper.String("ServiceId"),
				Values: []*string{
					&serviceId,
				}})
		}
		if serviceName != "" {
			request.Filters = append(request.Filters, &apigateway.Filter{Name: helper.String("ServiceName"),
				Values: []*string{
					&serviceName,
				}})
		}
	}

	var limit int64 = 100
	var offset int64 = 0

	request.Limit = &limit
	request.Offset = &offset

	for {
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseAPIGatewayClient().DescribeServicesStatus(request)
		if err != nil {
			errRet = err
			return
		}
		if response.Response.Result == nil {
			errRet = fmt.Errorf("TencentCloud SDK %s return empty response", request.GetAction())
			return
		}
		if len(response.Response.Result.ServiceSet) > 0 {
			infos = append(infos, response.Response.Result.ServiceSet...)
		}
		if len(response.Response.Result.ServiceSet) < int(limit) {
			return
		}
		offset += limit
	}
}

func (me *APIGatewayService) DescribeApisStatus(ctx context.Context,
	serviceId, apiName, apiId string) (infos []*apigateway.DesApisStatus, errRet error) {

	request := apigateway.NewDescribeApisStatusRequest()
	request.ServiceId = &serviceId

	if apiId != "" || apiName != "" {
		request.Filters = make([]*apigateway.Filter, 0, 2)
		if apiId != "" {
			request.Filters = append(request.Filters, &apigateway.Filter{Name: helper.String("ApiId"),
				Values: []*string{
					&apiId,
				}})
		}
		if apiName != "" {
			request.Filters = append(request.Filters, &apigateway.Filter{Name: helper.String("ApiName"),
				Values: []*string{
					&apiName,
				}})
		}
	}

	var limit int64 = 100
	var offset int64 = 0

	request.Limit = &limit
	request.Offset = &offset

	for {
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseAPIGatewayClient().DescribeApisStatus(request)
		if err != nil {
			errRet = err
			return
		}
		if response.Response.Result == nil {
			errRet = fmt.Errorf("TencentCloud SDK %s return empty response", request.GetAction())
			return
		}
		if len(response.Response.Result.ApiIdStatusSet) > 0 {
			infos = append(infos, response.Response.Result.ApiIdStatusSet...)
		}
		if len(response.Response.Result.ApiIdStatusSet) < int(limit) {
			return
		}
		offset += limit
	}
}

//limit & domain
func (me *APIGatewayService) DescribeServiceEnvironmentStrategyList(ctx context.Context,
	serviceId string) (environmentList []*apigateway.ServiceEnvironmentStrategy, errRet error) {
	var (
		request  = apigateway.NewDescribeServiceEnvironmentStrategyRequest()
		response *apigateway.DescribeServiceEnvironmentStrategyResponse
		err      error
		limit    int64 = 100
		offset   int64 = 0
	)

	if serviceId == "" {
		errRet = fmt.Errorf("serviceId is must not empty")
		return
	}

	request.ServiceId = &serviceId
	request.Limit = &limit
	request.Offset = &offset

	for {
		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			response, err = me.client.UseAPIGatewayClient().DescribeServiceEnvironmentStrategy(request)
			if err != nil {
				return retryError(err, InternalError)
			}
			return nil
		})
		if err != nil {
			if sdkError, ok := err.(*errors.TencentCloudSDKError); ok && sdkError.Code == SERVICE_ERR_CODE {
				return
			}
			errRet = err
			return
		}

		if response.Response == nil {
			return nil, fmt.Errorf("Response is nil, serviceId: %s ", serviceId)
		}

		if response.Response.Result == nil {
			return
		}

		environmentList = append(environmentList, response.Response.Result.EnvironmentList...)
		if len(response.Response.Result.EnvironmentList) < int(limit) {
			break
		}
		offset += limit
	}
	return
}

func (me *APIGatewayService) DescribeApiEnvironmentStrategyList(ctx context.Context,
	serviceId string, environmentNames []string, apiId string) (environmentApiList []*apigateway.ApiEnvironmentStrategy, errRet error) {
	var (
		request  = apigateway.NewDescribeApiEnvironmentStrategyRequest()
		err      error
		response *apigateway.DescribeApiEnvironmentStrategyResponse

		limit  int64 = 100
		offset int64 = 0
	)

	if serviceId == "" {
		errRet = fmt.Errorf("serviceId is must not empty")
		return
	}

	if apiId != "" {
		request.ApiId = &apiId
	}

	request.ServiceId = &serviceId
	if len(environmentNames) > 0 {
		request.EnvironmentNames = append(request.EnvironmentNames, helper.Strings(environmentNames)...)
	}

	request.Limit = &limit
	request.Offset = &offset

	for {
		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			response, err = me.client.UseAPIGatewayClient().DescribeApiEnvironmentStrategy(request)
			if err != nil {
				return retryError(err, InternalError)
			}
			return nil
		})
		if err != nil {
			if sdkErr, ok := err.(*errors.TencentCloudSDKError); ok && sdkErr.Code == SERVICE_ERR_CODE {
				return
			}
			errRet = err
			return
		}

		if response.Response == nil {
			return nil, fmt.Errorf("Response is nil, serviceId: %s ", serviceId)
		}

		if response.Response.Result == nil || response.Response.Result.ApiEnvironmentStrategySet == nil {
			return
		}

		environmentApiList = append(environmentApiList, response.Response.Result.ApiEnvironmentStrategySet...)
		if len(response.Response.Result.ApiEnvironmentStrategySet) < int(limit) {
			break
		}
		offset += limit
	}
	return
}

func (me *APIGatewayService) ModifyApiEnvironmentStrategy(ctx context.Context,
	serviceId string, strategy int64, environmentName string, apiIDs []string) (result bool, errRet error) {
	var (
		request  = apigateway.NewModifyApiEnvironmentStrategyRequest()
		err      error
		response *apigateway.ModifyApiEnvironmentStrategyResponse
	)

	request.ServiceId = &serviceId
	request.Strategy = &strategy
	request.EnvironmentName = &environmentName
	request.ApiIds = append(request.ApiIds, helper.Strings(apiIDs)...)

	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, err = me.client.UseAPIGatewayClient().ModifyApiEnvironmentStrategy(request)
		if err != nil {
			return retryError(err)
		}
		return nil
	})
	if err != nil {
		log.Printf("ModifyApiEnvironmentStrategy error: %v", err)
		errRet = err
		return
	}

	if response.Response == nil {
		return false, fmt.Errorf("Response is nil, serviceId: %s ", serviceId)
	}

	if response.Response.Result == nil {
		return
	}

	result = *response.Response.Result
	return
}

func (me *APIGatewayService) ModifyServiceEnvironmentStrategy(ctx context.Context,
	serviceId string, strategy int64, environmentName []string) (result bool, errRet error) {
	var (
		request  = apigateway.NewModifyServiceEnvironmentStrategyRequest()
		err      error
		response *apigateway.ModifyServiceEnvironmentStrategyResponse
	)

	request.ServiceId = &serviceId
	request.Strategy = &strategy
	request.EnvironmentNames = append(request.EnvironmentNames, helper.Strings(environmentName)...)

	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, err = me.client.UseAPIGatewayClient().ModifyServiceEnvironmentStrategy(request)
		if err != nil {
			return retryError(err)
		}
		return nil
	})
	if err != nil {
		errRet = err
		return
	}

	if response.Response == nil {
		return false, fmt.Errorf("Response is nil, serviceId: %s ", serviceId)
	}

	if response.Response.Result == nil {
		return
	}

	result = *response.Response.Result
	return
}

func (me *APIGatewayService) BindSubDomainService(ctx context.Context,
	serviceId, subDomain, protocol, netType, defaultDomain string, isDefaultMapping bool, certificateId string, pathMappings []string) (errRet error) {
	var (
		request = apigateway.NewBindSubDomainRequest()
		err     error
	)

	request.ServiceId = &serviceId
	request.SubDomain = &subDomain
	request.Protocol = &protocol
	request.NetType = &netType
	request.NetSubDomain = &defaultDomain
	request.IsDefaultMapping = &isDefaultMapping
	if certificateId != "" {
		request.CertificateId = &certificateId
	}
	for _, v := range pathMappings {
		results := strings.Split(v, "#")
		pathTmp := &apigateway.PathMapping{
			Path:        &results[0],
			Environment: &results[1],
		}
		request.PathMappingSet = append(request.PathMappingSet, pathTmp)
	}

	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		_, err = me.client.UseAPIGatewayClient().BindSubDomain(request)
		if err != nil {
			if ee, ok := err.(*errors.TencentCloudSDKError); ok {
				if ee.Code == CERTIFI_CATE_ID_EXPIRED || ee.Code == CERTIFICATE_ID_UNDER_VERIFY ||
					ee.Code == DOMAIN_NEED_BEIAN || ee.Code == EXCEEDED_DEFINE_MAPPING_LIMIT ||
					ee.Code == DOMAIN_RESOLVE_ERROR || ee.Code == DOMAIN_BIND_SERVICE {
					return resource.NonRetryableError(ee)
				}
			}
			return retryError(err)
		}
		return nil
	})
	if err != nil {
		errRet = err
		return
	}
	return
}

func (me *APIGatewayService) DescribeServiceSubDomainsService(ctx context.Context, serviceId, subDomain string) (resultList []*apigateway.DomainSetList, errRet error) {
	var (
		request  = apigateway.NewDescribeServiceSubDomainsRequest()
		err      error
		response *apigateway.DescribeServiceSubDomainsResponse

		limit  int64 = 100
		offset int64 = 0
	)
	request.ServiceId = &serviceId
	request.Limit = &limit
	request.Offset = &offset
	for {
		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			response, err = me.client.UseAPIGatewayClient().DescribeServiceSubDomains(request)
			if err != nil {
				return retryError(err, InternalError)
			}
			return nil
		})
		if err != nil {
			errRet = err
			return
		}
		if response.Response == nil || response.Response.Result == nil {
			errRet = fmt.Errorf("TencentCloud SDK %s return empty response", request.GetAction())
			return
		}

		resultList = append(resultList, response.Response.Result.DomainSet...)
		if len(response.Response.Result.DomainSet) < int(limit) {
			break
		}
		offset += limit
	}
	return
}

func (me *APIGatewayService) DescribeServiceSubDomainMappings(ctx context.Context, serviceId, subDomain string) (info *apigateway.ServiceSubDomainMappings, errRet error) {
	var (
		request  = apigateway.NewDescribeServiceSubDomainMappingsRequest()
		response *apigateway.DescribeServiceSubDomainMappingsResponse
		err      error
	)

	request.ServiceId = &serviceId
	request.SubDomain = &subDomain

	if err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, err = me.client.UseAPIGatewayClient().DescribeServiceSubDomainMappings(request)
		if err != nil {
			return retryError(err, InternalError)
		}
		return nil
	}); err != nil {
		errRet = err
		return
	}

	if response.Response == nil || response.Response.Result == nil {
		errRet = fmt.Errorf("TencentCloud SDK %s return empty response", request.GetAction())
		return
	}

	info = response.Response.Result
	return
}

func (me *APIGatewayService) ModifySubDomainService(ctx context.Context,
	serviceId, subDomain string, isDefaultMapping bool, certificateId, protocol, netType string, pathMappings []string) (errRet error) {
	var (
		request  = apigateway.NewModifySubDomainRequest()
		response *apigateway.ModifySubDomainResponse
		err      error
	)

	request.ServiceId = &serviceId
	request.SubDomain = &subDomain
	request.IsDefaultMapping = &isDefaultMapping
	if certificateId != "" {
		request.CertificateId = &certificateId
	}
	if protocol != "" {
		request.Protocol = &protocol
	}
	if netType != "" {
		request.NetType = &netType
	}
	for _, v := range pathMappings {
		results := strings.Split(v, "#")
		pathTmp := &apigateway.PathMapping{
			Path:        &results[0],
			Environment: &results[1],
		}
		request.PathMappingSet = append(request.PathMappingSet, pathTmp)
	}

	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, err = me.client.UseAPIGatewayClient().ModifySubDomain(request)
		if err != nil {
			return retryError(err)
		}
		return nil
	})
	if err != nil {
		errRet = err
		return
	}
	if response.Response == nil || response.Response.Result == nil {
		errRet = fmt.Errorf("TencentCloud SDK %s return empty response", request.GetAction())
		return
	}

	if !(*response.Response.Result) {
		errRet = fmt.Errorf("%s failed", request.GetAction())
		return
	}
	return
}

func (me *APIGatewayService) UnBindSubDomainService(ctx context.Context,
	serviceId, subDomain string) (errRet error) {
	var (
		request  = apigateway.NewUnBindSubDomainRequest()
		response *apigateway.UnBindSubDomainResponse
		err      error
	)

	request.ServiceId = &serviceId
	request.SubDomain = &subDomain

	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, err = me.client.UseAPIGatewayClient().UnBindSubDomain(request)
		if err != nil {
			return retryError(err)
		}
		return nil
	})
	if err != nil {
		errRet = err
		return
	}
	if response.Response == nil || response.Response.Result == nil {
		errRet = fmt.Errorf("TencentCloud SDK %s return empty response", request.GetAction())
		return
	}

	if !(*response.Response.Result) {
		errRet = fmt.Errorf("%s failed", request.GetAction())
		return
	}
	return
}

func (me *APIGatewayService) CreateIPStrategy(ctx context.Context,
	serviceId, strategyName, strategyType, strategyData string) (strategyId string, errRet error) {
	request := apigateway.NewCreateIPStrategyRequest()
	request.ServiceId = &serviceId
	request.StrategyName = &strategyName
	request.StrategyType = &strategyType
	request.StrategyData = &strategyData

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseAPIGatewayClient().CreateIPStrategy(request)
	if err != nil {
		errRet = err
		return
	}
	if response.Response.Result == nil || response.Response.Result.StrategyId == nil {
		errRet = fmt.Errorf("TencentCloud SDK %s return empty ", request.GetAction())
		return
	}
	strategyId = *response.Response.Result.StrategyId
	return
}

func (me *APIGatewayService) DescribeIPStrategyHas(ctx context.Context,
	serviceId, strategyId string) (has bool, errRet error) {

	request := apigateway.NewDescribeIPStrategysStatusRequest()
	request.ServiceId = &serviceId

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseAPIGatewayClient().DescribeIPStrategysStatus(request)
	if err != nil {
		if sdkErr, ok := err.(*errors.TencentCloudSDKError); ok && sdkErr.Code == SERVICE_ERR_CODE {
			return false, nil
		}
		errRet = err
		return
	}
	if response.Response.Result == nil {
		errRet = fmt.Errorf("TencentCloud SDK %s return empty response", request.GetAction())
		return
	}
	for _, status := range response.Response.Result.StrategySet {
		if *status.StrategyId == strategyId {
			has = true
			return
		}
	}

	return
}

func (me *APIGatewayService) DescribeIPStrategyStatus(ctx context.Context, serviceId,
	strategyId string) (ipStrategies *apigateway.IPStrategy, has bool, errRet error) {

	var apiList []*apigateway.DesApisStatus
	for _, env := range API_GATEWAY_SERVICE_ENVS {
		request := apigateway.NewDescribeIPStrategyRequest()

		request.ServiceId = &serviceId
		request.StrategyId = &strategyId
		request.EnvironmentName = &env

		var (
			limit  int64 = 100
			offset int64 = 0
		)

		request.Limit = &limit
		request.Offset = &offset

		for {
			ratelimit.Check(request.GetAction())
			response, err := me.client.UseAPIGatewayClient().DescribeIPStrategy(request)
			if err != nil {
				errRet = err
				return
			}
			if response.Response.Result == nil {
				errRet = fmt.Errorf("TencentCloud SDK %s return empty response", request.GetAction())
				return
			}
			if len(response.Response.Result.BindApis) > 0 {
				apiList = append(apiList, response.Response.Result.BindApis...)
			}
			if len(response.Response.Result.BindApis) < int(limit) {
				has = true
				ipStrategies = response.Response.Result
				ipStrategies.BindApis = apiList
				return
			}
			offset += limit
		}
	}
	return
}

func (me *APIGatewayService) UpdateIPStrategy(ctx context.Context, serviceId, strategyId, strategyData string) (errRet error) {
	request := apigateway.NewModifyIPStrategyRequest()
	request.StrategyId = &strategyId
	request.ServiceId = &serviceId
	request.StrategyData = &strategyData
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseAPIGatewayClient().ModifyIPStrategy(request)
	if err != nil {
		errRet = err
		return
	}
	if response.Response.Result == nil {
		errRet = fmt.Errorf("TencentCloud SDK %s return empty response", request.GetAction())
		return
	}
	if !*response.Response.Result {
		errRet = fmt.Errorf("update IP strategy fail")
		return
	}
	return
}

func (me *APIGatewayService) DeleteIPStrategy(ctx context.Context, serviceId, strategyId string) (errRet error) {
	request := apigateway.NewDeleteIPStrategyRequest()
	request.StrategyId = &strategyId
	request.ServiceId = &serviceId

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseAPIGatewayClient().DeleteIPStrategy(request)
	if err != nil {
		errRet = err
		return
	}
	if response.Response.Result == nil {
		errRet = fmt.Errorf("TencentCloud SDK %s return empty response", request.GetAction())
		return
	}
	if !*response.Response.Result {
		errRet = fmt.Errorf("delete IP strategy fail")
		return
	}
	return
}

func (me *APIGatewayService) CreateStrategyAttachment(ctx context.Context,
	serviceId, strategyId, envName, bindApiId string) (has bool, errRet error) {
	request := apigateway.NewBindIPStrategyRequest()
	var bindarr = []*string{&bindApiId}
	request.ServiceId = &serviceId
	request.StrategyId = &strategyId
	request.EnvironmentName = &envName
	request.BindApiIds = bindarr

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseAPIGatewayClient().BindIPStrategy(request)
	if err != nil {
		errRet = err
		return
	}
	if response.Response.Result == nil {
		errRet = fmt.Errorf("TencentCloud SDK %s return empty ", request.GetAction())
		return
	}
	has = *response.Response.Result
	return
}

func (me *APIGatewayService) DeleteStrategyAttachment(ctx context.Context,
	serviceId, strategyId, envName, bindApiId string) (has bool, errRet error) {
	request := apigateway.NewUnBindIPStrategyRequest()
	var unBindarr = []*string{&bindApiId}
	request.ServiceId = &serviceId
	request.StrategyId = &strategyId
	request.EnvironmentName = &envName
	request.UnBindApiIds = unBindarr

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseAPIGatewayClient().UnBindIPStrategy(request)
	if err != nil {
		errRet = err
		return
	}
	if response.Response.Result == nil {
		errRet = fmt.Errorf("TencentCloud SDK %s return empty ", request.GetAction())
		return
	}
	has = *response.Response.Result
	return
}

func (me *APIGatewayService) DescribeStrategyAttachment(ctx context.Context, serviceId, strategyId, bindApiId string) (has bool, errRet error) {
	ipStatus, _, err := me.DescribeIPStrategyStatus(ctx, serviceId, strategyId)
	if err != nil {
		if sdkError, ok := err.(*errors.TencentCloudSDKError); ok && sdkError.Code == SERVICE_ERR_CODE {
			return
		}
		errRet = err
		return
	}
	if ipStatus.BindApis == nil {
		has = true
		return
	}
	for _, bindApi := range ipStatus.BindApis {
		if *bindApi.ApiId == bindApiId {
			has = false
			return
		}
	}
	return
}

func (me *APIGatewayService) ReleaseService(ctx context.Context,
	serviceId, environmentName, releaseDesc string) (response *apigateway.ReleaseServiceResponse, err error) {

	request := apigateway.NewReleaseServiceRequest()
	request.ServiceId = &serviceId
	request.EnvironmentName = &environmentName
	request.ReleaseDesc = &releaseDesc

	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, err = me.client.UseAPIGatewayClient().ReleaseService(request)
		if err != nil {
			return retryError(err)
		}
		if response == nil || response.Response.Result == nil || response.Response.Result.ReleaseVersion == nil {
			return resource.NonRetryableError(fmt.Errorf("ReleaseService response is nil"))
		}
		return nil
	})
	return
}

func (me *APIGatewayService) DescribeServiceEnvironmentReleaseHistory(ctx context.Context,
	serviceId, envName string) (versionList []*apigateway.ServiceReleaseHistoryInfo, has bool, errRet error) {
	var (
		limit  int64 = 100
		offset int64 = 0
	)

	request := apigateway.NewDescribeServiceEnvironmentReleaseHistoryRequest()
	request.ServiceId = &serviceId
	request.EnvironmentName = &envName
	request.Limit = &limit
	request.Offset = &offset

	for {
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseAPIGatewayClient().DescribeServiceEnvironmentReleaseHistory(request)
		if err != nil {
			if sdkError, ok := err.(*errors.TencentCloudSDKError); ok && sdkError.Code == SERVICE_ERR_CODE {
				return
			}
			errRet = err
			return
		}
		if response.Response.Result == nil {
			errRet = fmt.Errorf("TencentCloud SDK %s return empty response", request.GetAction())
			return
		}
		if len(response.Response.Result.VersionList) > 0 {
			versionList = append(versionList, response.Response.Result.VersionList...)
		}

		if len(response.Response.Result.VersionList) < int(limit) {
			has = true
			return
		}
		offset += limit
	}
}

func (me *APIGatewayService) DescribeApiGatewayPluginById(ctx context.Context, pluginId string) (plugin *apigateway.Plugin, errRet error) {
	logId := getLogId(ctx)

	request := apigateway.NewDescribePluginsRequest()
	request.PluginIds = []*string{&pluginId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseAPIGatewayClient().DescribePlugins(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.Result.PluginSet) < 1 {
		return
	}

	plugin = response.Response.Result.PluginSet[0]
	return
}

func (me *APIGatewayService) DeleteApiGatewayPluginById(ctx context.Context, pluginId string) (errRet error) {
	logId := getLogId(ctx)

	request := apigateway.NewDeletePluginRequest()
	request.PluginId = &pluginId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseAPIGatewayClient().DeletePlugin(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *APIGatewayService) DescribeApiGatewayPluginAttachmentById(ctx context.Context, pluginId string, serviceId string, environmentName string, apiId string) (pluginAttachment *apigateway.AttachedApiInfo, errRet error) {
	logId := getLogId(ctx)

	request := apigateway.NewDescribePluginApisRequest()
	request.PluginId = &pluginId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseAPIGatewayClient().DescribePluginApis(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.Result.AttachedApis) < 1 {
		return
	}

	for _, api := range response.Response.Result.AttachedApis {
		if *api.ServiceId == serviceId && *api.Environment == environmentName && *api.ApiId == apiId {
			pluginAttachment = api
			return
		}
	}
	return
}

func (me *APIGatewayService) DeleteApiGatewayPluginAttachmentById(ctx context.Context, pluginId string, serviceId string, environmentName string, apiId string) (errRet error) {
	logId := getLogId(ctx)

	request := apigateway.NewDetachPluginRequest()
	request.PluginId = &pluginId
	request.ServiceId = &serviceId
	request.EnvironmentName = &environmentName
	request.ApiId = &apiId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseAPIGatewayClient().DetachPlugin(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}
