package tencentcloud

import (
	apigateway "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apigateway/v20180808"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
	"log"
)

type ApigatewayService struct {
	client *connectivity.TencentCloudClient
}

func (me *ApigatewayService) DescribeApigatewayAPIById(ctx context.Context, apiId string) (API *apigateway.ApiInfo, errRet error) {
	logId := getLogId(ctx)

	request := apigateway.NewDescribeApiRequest()
	request.ApiId = &apiId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseApigatewayClient().DescribeApi(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.ApiInfo) < 1 {
		return
	}

	API = response.Response.ApiInfo[0]
	return
}

func (me *ApigatewayService) DeleteApigatewayAPIById(ctx context.Context, apiId string) (errRet error) {
	logId := getLogId(ctx)

	request := apigateway.NewDeleteApiRequest()
	request.ApiId = &apiId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseApigatewayClient().DeleteApi(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *ApigatewayService) DescribeApigatewayApiAppById(ctx context.Context, apiAppId string) (apiApp *apigateway.ApiAppApiInfos, errRet error) {
	logId := getLogId(ctx)

	request := apigateway.NewDescribeApiBindApiAppsStatusRequest()
	request.ApiAppId = &apiAppId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseApigatewayClient().DescribeApiBindApiAppsStatus(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.ApiAppApiInfos) < 1 {
		return
	}

	apiApp = response.Response.ApiAppApiInfos[0]
	return
}

func (me *ApigatewayService) DeleteApigatewayApiAppById(ctx context.Context, apiAppId string) (errRet error) {
	logId := getLogId(ctx)

	request := apigateway.NewUnbindApiAppRequest()
	request.ApiAppId = &apiAppId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseApigatewayClient().UnbindApiApp(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *ApigatewayService) DescribeApigatewayPluginById(ctx context.Context, pluginId string) (plugin *apigateway.AttachedApiSummary, errRet error) {
	logId := getLogId(ctx)

	request := apigateway.NewDescribePluginApisRequest()
	request.PluginId = &pluginId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseApigatewayClient().DescribePluginApis(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.AttachedApiSummary) < 1 {
		return
	}

	plugin = response.Response.AttachedApiSummary[0]
	return
}

func (me *ApigatewayService) DeleteApigatewayPluginById(ctx context.Context, pluginId string) (errRet error) {
	logId := getLogId(ctx)

	request := apigateway.NewDetachPluginRequest()
	request.PluginId = &pluginId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseApigatewayClient().DetachPlugin(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *ApigatewayService) DescribeApigatewayServiceById(ctx context.Context, serviceId string) (service *apigateway.DomainSets, errRet error) {
	logId := getLogId(ctx)

	request := apigateway.NewDescribeServiceSubDomainsRequest()
	request.ServiceId = &serviceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseApigatewayClient().DescribeServiceSubDomains(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.DomainSets) < 1 {
		return
	}

	service = response.Response.DomainSets[0]
	return
}

func (me *ApigatewayService) DeleteApigatewayServiceById(ctx context.Context, serviceId string) (errRet error) {
	logId := getLogId(ctx)

	request := apigateway.NewUnBindSubDomainRequest()
	request.ServiceId = &serviceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseApigatewayClient().UnBindSubDomain(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *ApigatewayService) DescribeApigatewayUsagePlanById(ctx context.Context, usagePlanId string) (usagePlan *apigateway.UsagePlanEnvironmentStatus, errRet error) {
	logId := getLogId(ctx)

	request := apigateway.NewDescribeUsagePlanEnvironmentsRequest()
	request.UsagePlanId = &usagePlanId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseApigatewayClient().DescribeUsagePlanEnvironments(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.UsagePlanEnvironmentStatus) < 1 {
		return
	}

	usagePlan = response.Response.UsagePlanEnvironmentStatus[0]
	return
}

func (me *ApigatewayService) DeleteApigatewayUsagePlanById(ctx context.Context, usagePlanId string) (errRet error) {
	logId := getLogId(ctx)

	request := apigateway.NewUnBindSecretIdsRequest()
	request.UsagePlanId = &usagePlanId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseApigatewayClient().UnBindSecretIds(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *ApigatewayService) DescribeApigatewayApiKeyById(ctx context.Context, accessKeyId string) (apiKey *apigateway.ApiKeysStatus, errRet error) {
	logId := getLogId(ctx)

	request := apigateway.NewDescribeApiKeysStatusRequest()
	request.AccessKeyId = &accessKeyId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseApigatewayClient().DescribeApiKeysStatus(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.ApiKeysStatus) < 1 {
		return
	}

	apiKey = response.Response.ApiKeysStatus[0]
	return
}

func (me *ApigatewayService) DescribeApigatewayUpstreamByFilter(ctx context.Context, param map[string]interface{}) (upstream []*apigateway.DescribeUpstreamBindApis, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = apigateway.NewDescribeUpstreamBindApisRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "UpstreamId" {
			request.UpstreamId = v.(*string)
		}
		if k == "Filters" {
			request.Filters = v.([]*apigateway.Filter)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseApigatewayClient().DescribeUpstreamBindApis(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Result) < 1 {
			break
		}
		upstream = append(upstream, response.Response.Result...)
		if len(response.Response.Result) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *ApigatewayService) DescribeApigatewayUpstreamById(ctx context.Context, upstreamId string) (upstream *apigateway.DescribeUpstreamInfo, errRet error) {
	logId := getLogId(ctx)

	request := apigateway.NewDescribeUpstreamsRequest()
	request.UpstreamId = &upstreamId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseApigatewayClient().DescribeUpstreams(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.DescribeUpstreamInfo) < 1 {
		return
	}

	upstream = response.Response.DescribeUpstreamInfo[0]
	return
}

func (me *ApigatewayService) DeleteApigatewayUpstreamById(ctx context.Context, upstreamId string) (errRet error) {
	logId := getLogId(ctx)

	request := apigateway.NewDeleteUpstreamRequest()
	request.UpstreamId = &upstreamId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseApigatewayClient().DeleteUpstream(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}
