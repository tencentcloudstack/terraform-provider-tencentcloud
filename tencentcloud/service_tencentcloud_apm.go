package tencentcloud

import (
	"context"
	"log"

	apm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apm/v20210622"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type ApmService struct {
	client *connectivity.TencentCloudClient
}

func (me *ApmService) DescribeApmInstanceById(ctx context.Context, instanceId string) (instance *apm.ApmInstanceDetail, errRet error) {
	logId := getLogId(ctx)

	request := apm.NewDescribeApmInstancesRequest()
	request.InstanceIds = []*string{&instanceId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseApmClient().DescribeApmInstances(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.Instances) < 1 {
		return
	}

	instance = response.Response.Instances[0]
	return
}

func (me *ApmService) DeleteApmInstanceById(ctx context.Context, instanceId string) (errRet error) {
	logId := getLogId(ctx)

	request := apm.NewTerminateApmInstanceRequest()
	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseApmClient().TerminateApmInstance(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}
