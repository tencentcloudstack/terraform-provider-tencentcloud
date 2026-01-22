package gwlb

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	gwlbv20240906 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/gwlb/v20240906"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func NewGwlbService(client *connectivity.TencentCloudClient) GwlbService {
	return GwlbService{client: client}
}

type GwlbService struct {
	client *connectivity.TencentCloudClient
}

func (me *GwlbService) DescribeGwlbInstanceById(ctx context.Context, instanceId string) (gatewayLoadBalancer *gwlbv20240906.GatewayLoadBalancer, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := gwlbv20240906.NewDescribeGatewayLoadBalancersRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	request.LoadBalancerIds = helper.Strings([]string{instanceId})

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, e := me.client.UseGwlbV20240906Client().DescribeGatewayLoadBalancers(request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
		if len(response.Response.LoadBalancerSet) > 0 {
			gatewayLoadBalancer = response.Response.LoadBalancerSet[0]
		}
		return nil
	})
	if err != nil {
		errRet = err
		return
	}
	return
}

func (me *GwlbService) DescribeGwlbTargetGroupById(ctx context.Context, instanceId string) (ret *gwlbv20240906.TargetGroupInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := gwlbv20240906.NewDescribeTargetGroupListRequest()

	request.TargetGroupIds = helper.Strings([]string{instanceId})

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, e := me.client.UseGwlbV20240906Client().DescribeTargetGroupList(request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
		if len(response.Response.TargetGroupSet) > 0 {
			ret = response.Response.TargetGroupSet[0]
		}
		return nil
	})
	if err != nil {
		errRet = err
		return
	}

	return
}

func (me *GwlbService) DescribeTargetGroupInstancesById(ctx context.Context, instanceId string) (targetGroupBackends []*gwlbv20240906.TargetGroupBackend, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := gwlbv20240906.NewDescribeTargetGroupInstancesRequest()

	request.Filters = []*gwlbv20240906.Filter{
		{
			Name:   helper.String("TargetGroupId"),
			Values: helper.Strings([]string{instanceId}),
		},
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, e := me.client.UseGwlbV20240906Client().DescribeTargetGroupInstances(request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
		targetGroupBackends = response.Response.TargetGroupInstanceSet
		return nil
	})
	if err != nil {
		errRet = err
		return
	}

	return
}

func (me *GwlbService) DescribeGwlbTargetGroupRegisterInstancesById(ctx context.Context, targetGroupId string) (ret *gwlbv20240906.DescribeTargetGroupInstancesResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := gwlbv20240906.NewDescribeTargetGroupInstancesRequest()
	request.Filters = []*gwlbv20240906.Filter{
		{
			Name:   helper.String("TargetGroupId"),
			Values: helper.Strings([]string{targetGroupId}),
		},
	}
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, e := me.client.UseGwlbV20240906Client().DescribeTargetGroupInstances(request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
		ret = response.Response
		return nil
	})
	if err != nil {
		errRet = err
		return
	}

	return
}

func (me *GwlbService) RegisterTargetGroupInstances(ctx context.Context, targetGroupId string, targetGroupInstances []*gwlbv20240906.TargetGroupInstance) error {
	logId := tccommon.GetLogId(ctx)
	request := gwlbv20240906.NewRegisterTargetGroupInstancesRequest()
	request.TargetGroupInstances = targetGroupInstances
	request.TargetGroupId = helper.String(targetGroupId)
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := me.client.UseGwlbV20240906Client().RegisterTargetGroupInstancesWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (me *GwlbService) DeregisterTargetGroupInstances(ctx context.Context, targetGroupId string, targetGroupInstances []*gwlbv20240906.TargetGroupInstance) error {
	logId := tccommon.GetLogId(ctx)

	request := gwlbv20240906.NewDeregisterTargetGroupInstancesRequest()
	request.TargetGroupId = helper.String(targetGroupId)
	request.TargetGroupInstances = targetGroupInstances
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := me.client.UseGwlbV20240906Client().DeregisterTargetGroupInstancesWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (me *GwlbService) DescribeTaskStatus(ctx context.Context, taskId string) (status *int64, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := gwlbv20240906.NewDescribeTaskStatusRequest()

	request.TaskId = helper.String(taskId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, e := me.client.UseGwlbV20240906Client().DescribeTaskStatus(request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
		status = response.Response.Status
		return nil
	})
	if err != nil {
		errRet = err
		return
	}

	return
}

func (me *GwlbService) TaskStatusRefreshFunc(ctx context.Context, taskId string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		status, err := me.DescribeTaskStatus(ctx, taskId)

		if err != nil {
			return nil, "", err
		}

		if status == nil {
			return nil, "", fmt.Errorf("task status is nil")
		}
		return status, helper.Int64ToStr(*status), nil
	}
}
