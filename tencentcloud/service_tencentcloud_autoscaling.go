package tencentcloud

import (
	autoscaling "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/as/v20180419"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
	"log"
)

type AutoscalingService struct {
	client *connectivity.TencentCloudClient
}

func (me *AutoscalingService) DescribeAutoscalingAdvicesByFilter(ctx context.Context, param map[string]interface{}) (advices []*autoscaling.AutoScalingAdvice, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = autoscaling.NewDescribeAutoScalingAdvicesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "AutoScalingGroupIds" {
			request.AutoScalingGroupIds = v.([]*string)
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
		response, err := me.client.UseAutoscalingClient().DescribeAutoScalingAdvices(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.AutoScalingAdviceSet) < 1 {
			break
		}
		advices = append(advices, response.Response.AutoScalingAdviceSet...)
		if len(response.Response.AutoScalingAdviceSet) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *AutoscalingService) DescribeAutoscalingLimitsByFilter(ctx context.Context, param map[string]interface{}) (limits []*autoscaling.DescribeAccountLimitsResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = autoscaling.NewDescribeAccountLimitsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseAutoscalingClient().DescribeAccountLimits(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.MaxNumberOfLaunchConfigurations) < 1 {
			break
		}
		limits = append(limits, response.Response.MaxNumberOfLaunchConfigurations...)
		if len(response.Response.MaxNumberOfLaunchConfigurations) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *AutoscalingService) DescribeAutoscalingLastActivityByFilter(ctx context.Context, param map[string]interface{}) (lastActivity []*autoscaling.Activity, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = autoscaling.NewDescribeAutoScalingGroupLastActivitiesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "AutoScalingGroupIds" {
			request.AutoScalingGroupIds = v.([]*string)
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
		response, err := me.client.UseAutoscalingClient().DescribeAutoScalingGroupLastActivities(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.ActivitySet) < 1 {
			break
		}
		lastActivity = append(lastActivity, response.Response.ActivitySet...)
		if len(response.Response.ActivitySet) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *AutoscalingService) DescribeAutoscalingAutoScalingGroupById(ctx context.Context, autoScalingGroupId string) (autoScalingGroup *autoscaling.AutoScalingGroup, errRet error) {
	logId := getLogId(ctx)

	request := autoscaling.NewDescribeAutoScalingGroupsRequest()
	request.AutoScalingGroupId = &autoScalingGroupId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseAutoscalingClient().DescribeAutoScalingGroups(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.AutoScalingGroup) < 1 {
		return
	}

	autoScalingGroup = response.Response.AutoScalingGroup[0]
	return
}
