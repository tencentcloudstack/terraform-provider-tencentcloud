package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	tem "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tem/v20210701"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type TemService struct {
	client *connectivity.TencentCloudClient
}

func (me *TemService) DescribeTemEnvironmentStatus(ctx context.Context, environmentId string) (environment *tem.NamespaceStatusInfo, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = tem.NewDescribeEnvironmentStatusRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()
	request.EnvironmentIds = []*string{&environmentId}

	response, err := me.client.UseTemClient().DescribeEnvironmentStatus(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	if len(response.Response.Result) < 1 {
		return
	}
	environment = response.Response.Result[0]
	return
}

func (me *TemService) DescribeTemEnvironment(ctx context.Context, environmentId string) (environment *tem.DescribeEnvironmentResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = tem.NewDescribeEnvironmentRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()
	request.EnvironmentId = &environmentId

	response, err := me.client.UseTemClient().DescribeEnvironment(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	environment = response.Response
	return
}

func (me *TemService) DeleteTemEnvironmentById(ctx context.Context, environmentId string) (errRet error) {
	logId := getLogId(ctx)

	request := tem.NewDestroyEnvironmentRequest()
	request.EnvironmentId = &environmentId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTemClient().DestroyEnvironment(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TemService) DescribeTemApplication(ctx context.Context, applicationId string) (application *tem.DescribeApplicationsResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = tem.NewDescribeApplicationsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()
	request.ApplicationId = &applicationId

	response, err := me.client.UseTemClient().DescribeApplications(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	application = response.Response
	return
}

func (me *TemService) DeleteTemApplicationById(ctx context.Context, applicationId string) (errRet error) {
	logId := getLogId(ctx)

	request := tem.NewDeleteApplicationRequest()
	request.ApplicationId = &applicationId
	request.EnvironmentId = helper.String("")

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTemClient().DeleteApplication(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TemService) DescribeTemWorkload(ctx context.Context, environmentId string, applicationId string) (workload *tem.DescribeApplicationInfoResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = tem.NewDescribeApplicationInfoRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()
	request.EnvironmentId = &environmentId
	request.ApplicationId = &applicationId

	response, err := me.client.UseTemClient().DescribeApplicationInfo(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	workload = response.Response
	return
}

func (me *TemService) DeleteTemWorkloadById(ctx context.Context, environmentId string, applicationId string) (errRet error) {
	logId := getLogId(ctx)

	request := tem.NewDeleteApplicationRequest()
	request.EnvironmentId = &environmentId
	request.ApplicationId = &applicationId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTemClient().DeleteApplication(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TemService) DescribeTemAppConfig(ctx context.Context, environmentId string, name string) (appConfig *tem.ConfigData, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = tem.NewDescribeConfigDataRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()
	request.EnvironmentId = &environmentId
	request.Name = &name

	response, err := me.client.UseTemClient().DescribeConfigData(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	appConfig = response.Response.Result
	return
}

func (me *TemService) DeleteTemAppConfigById(ctx context.Context, environmentId string, name string) (errRet error) {
	logId := getLogId(ctx)

	request := tem.NewDestroyConfigDataRequest()
	request.EnvironmentId = &environmentId
	request.Name = &name

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTemClient().DestroyConfigData(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TemService) DescribeTemLogConfig(ctx context.Context, environmentId string, applicationId string, name string) (logConfig *tem.LogConfig, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = tem.NewDescribeLogConfigRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()
	request.EnvironmentId = &environmentId
	request.ApplicationId = &applicationId
	request.Name = &name

	response, err := me.client.UseTemClient().DescribeLogConfig(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	logConfig = response.Response.Result
	return
}

func (me *TemService) DeleteTemLogConfigById(ctx context.Context, environmentId string, applicationId string, name string) (errRet error) {
	logId := getLogId(ctx)

	request := tem.NewDestroyLogConfigRequest()
	request.EnvironmentId = &environmentId
	request.ApplicationId = &applicationId
	request.Name = &name

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTemClient().DestroyLogConfig(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TemService) DescribeTemScaleRule(ctx context.Context, environmentId string, applicationId string, scaleRuleId string) (scaleRule *tem.Autoscaler, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = tem.NewDescribeApplicationAutoscalerListRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()
	request.EnvironmentId = &environmentId
	request.ApplicationId = &applicationId

	response, err := me.client.UseTemClient().DescribeApplicationAutoscalerList(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	scaleRuleList := response.Response.Result
	for _, rule := range scaleRuleList {
		if *rule.AutoscalerId == scaleRuleId {
			scaleRule = rule
			return
		}
	}
	return
}

func (me *TemService) DisableTemScaleRuleById(ctx context.Context, environmentId string, applicationId string, scaleRuleId string) (errRet error) {
	logId := getLogId(ctx)

	request := tem.NewDisableApplicationAutoscalerRequest()
	request.EnvironmentId = &environmentId
	request.ApplicationId = &applicationId
	request.AutoscalerId = &scaleRuleId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTemClient().DisableApplicationAutoscaler(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TemService) DeleteTemScaleRuleById(ctx context.Context, environmentId string, applicationId string, scaleRuleId string) (errRet error) {
	logId := getLogId(ctx)

	request := tem.NewDeleteApplicationAutoscalerRequest()
	request.EnvironmentId = &environmentId
	request.ApplicationId = &applicationId
	request.AutoscalerId = &scaleRuleId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTemClient().DeleteApplicationAutoscaler(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TemService) DescribeTemGateway(ctx context.Context, environmentId string, ingressName string) (gateway *tem.IngressInfo, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = tem.NewDescribeIngressRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()
	request.EnvironmentId = &environmentId
	request.IngressName = &ingressName
	request.ClusterNamespace = helper.String("default")

	err := resource.Retry(2*readRetryTimeout, func() *resource.RetryError {
		response, errRet := me.client.UseTemClient().DescribeIngress(request)
		if errRet != nil {
			return retryError(errRet, InternalError)
		}
		gateway = response.Response.Result
		if *gateway.ClbId != "" && *gateway.Vip != "" {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
			return nil
		}
		return resource.RetryableError(fmt.Errorf("gateway clb is not ready..."))
	})

	if err != nil {
		errRet = err
		return
	}
	return
}

func (me *TemService) DeleteTemGatewayById(ctx context.Context, environmentId string, ingressName string) (errRet error) {
	logId := getLogId(ctx)

	request := tem.NewDeleteIngressRequest()
	request.EnvironmentId = &environmentId
	request.IngressName = &ingressName
	request.ClusterNamespace = helper.String("default")

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTemClient().DeleteIngress(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TemService) DescribeTemApplicationServiceById(ctx context.Context, environmentId string, applicationId string) (applicationService *tem.DescribeApplicationServiceListResponseParams, errRet error) {
	logId := getLogId(ctx)

	request := tem.NewDescribeApplicationServiceListRequest()
	request.EnvironmentId = &environmentId
	request.ApplicationId = &applicationId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTemClient().DescribeApplicationServiceList(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	applicationService = response.Response
	return
}

func (me *TemService) DeleteTemApplicationServiceById(ctx context.Context, environmentId string, applicationId string, serviceName string) (errRet error) {
	logId := getLogId(ctx)

	request := tem.NewDeleteApplicationServiceRequest()
	request.EnvironmentId = &environmentId
	request.ApplicationId = &applicationId
	request.ServiceName = &serviceName

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTemClient().DeleteApplicationService(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}
