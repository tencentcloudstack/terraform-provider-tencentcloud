package tencentcloud

import (
	apiGateway "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apigateway/v20180808"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
	"log"
)

type ApiGatewayService struct {
	client *connectivity.TencentCloudClient
}

func (me *ApiGatewayService) DescribeApiGatewayPluginById(ctx context.Context, pluginId string) (plugin *apiGateway.PluginSummary, errRet error) {
	logId := getLogId(ctx)

	request := apiGateway.NewDescribePluginsRequest()
	request.PluginId = &pluginId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseApiGatewayClient().DescribePlugins(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.PluginSummary) < 1 {
		return
	}

	plugin = response.Response.PluginSummary[0]
	return
}

func (me *ApiGatewayService) DeleteApiGatewayPluginById(ctx context.Context, pluginId string) (errRet error) {
	logId := getLogId(ctx)

	request := apiGateway.NewDeletePluginRequest()
	request.PluginId = &pluginId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseApiGatewayClient().DeletePlugin(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *ApiGatewayService) DescribeApiGatewayPluginAttachmentById(ctx context.Context, pluginId string, serviceId string, environmentName string, apiId string) (pluginAttachment *apiGateway.AttachedApiSummary, errRet error) {
	logId := getLogId(ctx)

	request := apiGateway.NewDescribePluginApisRequest()
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

	response, err := me.client.UseApiGatewayClient().DescribePluginApis(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.AttachedApiSummary) < 1 {
		return
	}

	pluginAttachment = response.Response.AttachedApiSummary[0]
	return
}

func (me *ApiGatewayService) DeleteApiGatewayPluginAttachmentById(ctx context.Context, pluginId string, serviceId string, environmentName string, apiId string) (errRet error) {
	logId := getLogId(ctx)

	request := apiGateway.NewDetachPluginRequest()
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

	response, err := me.client.UseApiGatewayClient().DetachPlugin(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}
