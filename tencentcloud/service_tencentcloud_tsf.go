package tencentcloud

import (
	"context"
	"fmt"
	"log"

	tsf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tsf/v20180326"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type TsfService struct {
	client *connectivity.TencentCloudClient
}

func (me *TsfService) DescribeTsfApplicationConfigById(ctx context.Context, configId, configName string) (applicationConfig *tsf.Config, errRet error) {
	logId := getLogId(ctx)

	request := tsf.NewDescribeConfigsRequest()
	request.ConfigId = &configId
	if configId != "" {
		request.ConfigId = &configId
	}
	if configName != "" {
		request.ConfigName = &configName
	}
	if configId == "" && configName == "" {
		errRet = fmt.Errorf("`configId` and `configName` cannot both be empty")
		return
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTsfClient().DescribeConfigs(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response.Result == nil {
		return
	}

	if len(response.Response.Result.Content) < 1 {
		return
	}

	applicationConfig = response.Response.Result.Content[0]
	return
}

func (me *TsfService) DeleteTsfApplicationConfigById(ctx context.Context, configId string) (errRet error) {
	logId := getLogId(ctx)

	request := tsf.NewDeleteConfigRequest()
	request.ConfigId = &configId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTsfClient().DeleteConfig(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TsfService) DescribeTsfMicroserviceById(ctx context.Context, namespaceId, microserviceId, microserviceName string) (microservice *tsf.Microservice, errRet error) {
	logId := getLogId(ctx)

	request := tsf.NewDescribeMicroservicesRequest()
	request.NamespaceId = &namespaceId
	if microserviceId != "" {
		request.MicroserviceIdList = []*string{&microserviceId}
	}
	if microserviceName != "" {
		request.MicroserviceNameList = []*string{&microserviceName}
	}
	if microserviceId == "" && microserviceName == "" {
		errRet = fmt.Errorf("`microserviceId` and `microserviceName` cannot both be empty")
		return
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTsfClient().DescribeMicroservices(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.Result.Content) < 1 {
		return
	}

	microservice = response.Response.Result.Content[0]
	return
}

func (me *TsfService) DeleteTsfMicroserviceById(ctx context.Context, microserviceId string) (errRet error) {
	logId := getLogId(ctx)

	request := tsf.NewDeleteMicroserviceRequest()
	request.MicroserviceId = &microserviceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTsfClient().DeleteMicroservice(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}
