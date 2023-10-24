package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	cdwpg "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdwpg/v20201230"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type CdwpgService struct {
	client *connectivity.TencentCloudClient
}

func (me *CdwpgService) DescribeCdwpgInstanceById(ctx context.Context, instanceId string) (instance *cdwpg.SimpleInstanceInfo, errRet error) {
	logId := getLogId(ctx)

	request := cdwpg.NewDescribeInstanceInfoRequest()
	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCdwpgClient().DescribeInstanceInfo(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil {
		return
	}

	instance = response.Response.SimpleInstanceInfo
	return
}

func (me *CdwpgService) DeleteCdwpgInstanceById(ctx context.Context, instanceId string) (errRet error) {
	logId := getLogId(ctx)

	request := cdwpg.NewDestroyInstanceByApiRequest()
	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCdwpgClient().DestroyInstanceByApi(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *CdwpgService) InstanceStateRefreshFunc(instanceId string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		request := cdwpg.NewDescribeInstanceStateRequest()
		request.InstanceId = &instanceId
		ratelimit.Check(request.GetAction())
		object, err := me.client.UseCdwpgClient().DescribeInstanceState(request)

		if err != nil {
			return nil, "", err
		}
		if object == nil || object.Response == nil || object.Response.InstanceState == nil {
			return nil, "", nil
		}

		return object, *object.Response.InstanceState, nil
	}
}
