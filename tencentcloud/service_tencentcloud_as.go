package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	as "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/as/v20180419"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/connectivity"
)

type AsService struct {
	client *connectivity.TencentCloudClient
}

func (me *AsService) DescribeLaunchConfigurationById(ctx context.Context, configurationId string) (config *as.LaunchConfiguration, errRet error) {
	logId := getLogId(ctx)
	request := as.NewDescribeLaunchConfigurationsRequest()
	request.LaunchConfigurationIds = []*string{&configurationId}
	response, err := me.client.UseAsClient().DescribeLaunchConfigurations(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.LaunchConfigurationSet) < 1 {
		errRet = fmt.Errorf("configuration id is not found")
		return
	}
	config = response.Response.LaunchConfigurationSet[0]
	return
}

func (me *AsService) DescribeLaunchConfigurationByFilter(ctx context.Context, configurationId, configurationName string) (configs []*as.LaunchConfiguration, errRet error) {
	logId := getLogId(ctx)
	request := as.NewDescribeLaunchConfigurationsRequest()
	request.Filters = make([]*as.Filter, 0)
	if configurationId != "" {
		filter := &as.Filter{
			Name:   stringToPointer("launch-configuration-id"),
			Values: []*string{&configurationId},
		}
		request.Filters = append(request.Filters, filter)
	}
	if configurationName != "" {
		filter := &as.Filter{
			Name:   stringToPointer("launch-configuration-name"),
			Values: []*string{&configurationName},
		}
		request.Filters = append(request.Filters, filter)
	}

	offset := 0
	pageSize := 100
	configs = make([]*as.LaunchConfiguration, 0)
	for {
		request.Offset = intToPointer(offset)
		request.Limit = intToPointer(pageSize)
		response, err := me.client.UseAsClient().DescribeLaunchConfigurations(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.LaunchConfigurationSet) < 1 {
			break
		}

		configs = append(configs, response.Response.LaunchConfigurationSet...)

		if len(response.Response.LaunchConfigurationSet) < pageSize {
			break
		}

		offset += pageSize
	}
	return
}

func (me *AsService) DeleteLaunchConfiguration(ctx context.Context, configurationId string) error {
	logId := getLogId(ctx)
	request := as.NewDeleteLaunchConfigurationRequest()
	request.LaunchConfigurationId = &configurationId
	response, err := me.client.UseAsClient().DeleteLaunchConfiguration(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	return nil
}

func (me *AsService) DescribeAutoScalingGroupById(ctx context.Context, scalingGroupId string) (scalingGroup *as.AutoScalingGroup, errRet error) {
	logId := getLogId(ctx)
	request := as.NewDescribeAutoScalingGroupsRequest()
	request.AutoScalingGroupIds = []*string{&scalingGroupId}
	response, err := me.client.UseAsClient().DescribeAutoScalingGroups(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.AutoScalingGroupSet) < 1 {
		errRet = fmt.Errorf("configuration id is not found")
		return
	}
	scalingGroup = response.Response.AutoScalingGroupSet[0]
	return
}

func (me *AsService) DescribeAutoScalingGroupByFilter(ctx context.Context, scalingGroupId, configurationId, scalingGroupName string) (scalingGroups []*as.AutoScalingGroup, errRet error) {
	logId := getLogId(ctx)
	request := as.NewDescribeAutoScalingGroupsRequest()
	request.Filters = make([]*as.Filter, 0)
	if scalingGroupId != "" {
		filter := &as.Filter{
			Name:   stringToPointer("auto-scaling-group-id"),
			Values: []*string{&scalingGroupId},
		}
		request.Filters = append(request.Filters, filter)
	}
	if configurationId != "" {
		filter := &as.Filter{
			Name:   stringToPointer("launch-configuration-id"),
			Values: []*string{&configurationId},
		}
		request.Filters = append(request.Filters, filter)
	}
	if scalingGroupName != "" {
		filter := &as.Filter{
			Name:   stringToPointer("auto-scaling-group-name"),
			Values: []*string{&scalingGroupName},
		}
		request.Filters = append(request.Filters, filter)
	}

	offset := 0
	pageSize := 100
	scalingGroups = make([]*as.AutoScalingGroup, 0)
	for {
		request.Offset = intToPointer(offset)
		request.Limit = intToPointer(pageSize)
		response, err := me.client.UseAsClient().DescribeAutoScalingGroups(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.AutoScalingGroupSet) < 1 {
			break
		}

		scalingGroups = append(scalingGroups, response.Response.AutoScalingGroupSet...)

		if len(response.Response.AutoScalingGroupSet) < pageSize {
			break
		}

		offset += pageSize
	}
	return
}

// set the scaling group desired capacity to 0
func (me *AsService) ClearScalingGroupInstance(ctx context.Context, scalingGroupId string) error {
	logId := getLogId(ctx)
	request := as.NewModifyAutoScalingGroupRequest()
	request.AutoScalingGroupId = &scalingGroupId
	request.MinSize = intToPointer(0)
	request.MaxSize = intToPointer(0)
	request.DesiredCapacity = intToPointer(0)
	response, err := me.client.UseAsClient().ModifyAutoScalingGroup(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	return nil
}

func (me *AsService) DeleteScalingGroup(ctx context.Context, scalingGroupId string) error {
	logId := getLogId(ctx)
	request := as.NewDeleteAutoScalingGroupRequest()
	request.AutoScalingGroupId = &scalingGroupId
	response, err := me.client.UseAsClient().DeleteAutoScalingGroup(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	return nil
}

func (me *AsService) AttachInstances(ctx context.Context, scalingGroupId string, instanceIds []string) error {
	logId := getLogId(ctx)
	request := as.NewAttachInstancesRequest()
	request.AutoScalingGroupId = &scalingGroupId
	request.InstanceIds = make([]*string, 0, len(instanceIds))
	for i := range instanceIds {
		request.InstanceIds = append(request.InstanceIds, &instanceIds[i])
	}
	response, err := me.client.UseAsClient().AttachInstances(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	activityId := *response.Response.ActivityId

	err = resource.Retry(10*time.Minute, func() *resource.RetryError {
		status, err := me.DescribeActivityById(ctx, activityId)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if status == SCALING_GROUP_ACTIVITY_STATUS_INIT || status == SCALING_GROUP_ACTIVITY_STATUS_RUNNING {
			return resource.RetryableError(fmt.Errorf("attach status is running(%s)", status))
		}
		if status == SCALING_GROUP_ACTIVITY_STATUS_SUCCESSFUL {
			return nil
		}
		return resource.NonRetryableError(fmt.Errorf("attach status is failed(%s)", status))
	})
	if err != nil {
		return err
	}
	return nil
}

func (me *AsService) DescribeActivityById(ctx context.Context, activityId string) (status string, errRet error) {
	logId := getLogId(ctx)
	request := as.NewDescribeAutoScalingActivitiesRequest()
	request.ActivityIds = []*string{&activityId}
	response, err := me.client.UseAsClient().DescribeAutoScalingActivities(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.ActivitySet) < 1 {
		errRet = fmt.Errorf("activity id set is nil")
	}
	status = *response.Response.ActivitySet[0].StatusCode
	return
}

func (me *AsService) DetachInstances(ctx context.Context, scalingGroupId string, instanceIds []string) error {
	logId := getLogId(ctx)
	request := as.NewDetachInstancesRequest()
	request.AutoScalingGroupId = &scalingGroupId
	request.InstanceIds = make([]*string, 0, len(instanceIds))
	for i := range instanceIds {
		request.InstanceIds = append(request.InstanceIds, &instanceIds[i])
	}
	response, err := me.client.UseAsClient().DetachInstances(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	activityId := *response.Response.ActivityId

	err = resource.Retry(10*time.Minute, func() *resource.RetryError {
		status, err := me.DescribeActivityById(ctx, activityId)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if status == SCALING_GROUP_ACTIVITY_STATUS_INIT || status == SCALING_GROUP_ACTIVITY_STATUS_RUNNING {
			return resource.RetryableError(fmt.Errorf("detach status is running(%s)", status))
		}
		if status == SCALING_GROUP_ACTIVITY_STATUS_SUCCESSFUL {
			return nil
		}
		return resource.NonRetryableError(fmt.Errorf("detach status is failed(%s)", status))
	})
	if err != nil {
		return err
	}
	return nil
}

func (me *AsService) DescribeAutoScalingAttachment(ctx context.Context, scalingGroupId string) (instanceIds []string, errRet error) {
	logId := getLogId(ctx)
	request := as.NewDescribeAutoScalingInstancesRequest()
	request.Filters = []*as.Filter{
		{
			Name:   stringToPointer("auto-scaling-group-id"),
			Values: []*string{&scalingGroupId},
		},
	}
	response, err := me.client.UseAsClient().DescribeAutoScalingInstances(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	instanceIds = make([]string, 0)
	for _, instance := range response.Response.AutoScalingInstanceSet {
		if *instance.CreationType == "MANUAL_ATTACHING" {
			instanceIds = append(instanceIds, *instance.InstanceId)
		}
	}
	return
}

func (me *AsService) DescribeScalingPolicyById(ctx context.Context, scalingPolicyId string) (scalingPolicy *as.ScalingPolicy, errRet error) {
	logId := getLogId(ctx)
	request := as.NewDescribeScalingPoliciesRequest()
	request.AutoScalingPolicyIds = []*string{&scalingPolicyId}
	response, err := me.client.UseAsClient().DescribeScalingPolicies(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.ScalingPolicySet) < 1 {
		errRet = fmt.Errorf("scaling policy id is not found")
		return
	}
	scalingPolicy = response.Response.ScalingPolicySet[0]
	return
}

func (me *AsService) DescribeScalingPolicyByFilter(ctx context.Context, policyId, policyName, scalingGroupId string) (scalingPolicies []*as.ScalingPolicy, errRet error) {
	logId := getLogId(ctx)
	request := as.NewDescribeScalingPoliciesRequest()
	request.Filters = make([]*as.Filter, 0)
	if policyId != "" {
		filter := &as.Filter{
			Name:   stringToPointer("auto-scaling-policy-id"),
			Values: []*string{&policyId},
		}
		request.Filters = append(request.Filters, filter)
	}
	if policyName != "" {
		filter := &as.Filter{
			Name:   stringToPointer("scaling-policy-name"),
			Values: []*string{&policyName},
		}
		request.Filters = append(request.Filters, filter)
	}
	if scalingGroupId != "" {
		filter := &as.Filter{
			Name:   stringToPointer("auto-scaling-group-id"),
			Values: []*string{&scalingGroupId},
		}
		request.Filters = append(request.Filters, filter)
	}

	offset := 0
	pageSize := 100
	scalingPolicies = make([]*as.ScalingPolicy, 0)
	for {
		request.Offset = intToPointer(offset)
		request.Limit = intToPointer(pageSize)
		response, err := me.client.UseAsClient().DescribeScalingPolicies(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.ScalingPolicySet) < 1 {
			break
		}

		scalingPolicies = append(scalingPolicies, response.Response.ScalingPolicySet...)

		if len(response.Response.ScalingPolicySet) < pageSize {
			break
		}

		offset += pageSize
	}
	return
}

func (me *AsService) DeleteScalingPolicy(ctx context.Context, scalingPolicyId string) error {
	logId := getLogId(ctx)
	request := as.NewDeleteScalingPolicyRequest()
	request.AutoScalingPolicyId = &scalingPolicyId
	response, err := me.client.UseAsClient().DeleteScalingPolicy(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	return nil
}

func (me *AsService) DescribeScheduledActionById(ctx context.Context, scheduledActionId string) (scheduledAction *as.ScheduledAction, errRet error) {
	logId := getLogId(ctx)
	request := as.NewDescribeScheduledActionsRequest()
	request.ScheduledActionIds = []*string{&scheduledActionId}
	response, err := me.client.UseAsClient().DescribeScheduledActions(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.ScheduledActionSet) < 1 {
		errRet = fmt.Errorf("scheduled action id is not found")
		return
	}
	scheduledAction = response.Response.ScheduledActionSet[0]
	return
}

func (me *AsService) DeleteScheduledAction(ctx context.Context, scheduledActonId string) error {
	logId := getLogId(ctx)
	request := as.NewDeleteScheduledActionRequest()
	request.ScheduledActionId = &scheduledActonId
	response, err := me.client.UseAsClient().DeleteScheduledAction(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	return nil
}

func (me *AsService) DescribeLifecycleHookById(ctx context.Context, lifecycleHookId string) (lifecycleHook *as.LifecycleHook, errRet error) {
	logId := getLogId(ctx)
	request := as.NewDescribeLifecycleHooksRequest()
	request.LifecycleHookIds = []*string{&lifecycleHookId}
	response, err := me.client.UseAsClient().DescribeLifecycleHooks(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.LifecycleHookSet) < 1 {
		errRet = fmt.Errorf("lifecycle hook id is not found")
		return
	}
	lifecycleHook = response.Response.LifecycleHookSet[0]
	return
}

func (me *AsService) DeleteLifecycleHook(ctx context.Context, lifecycleHookId string) error {
	logId := getLogId(ctx)
	request := as.NewDeleteLifecycleHookRequest()
	request.LifecycleHookId = &lifecycleHookId
	response, err := me.client.UseAsClient().DeleteLifecycleHook(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	return nil
}

func (me *AsService) DescribeNotificationById(ctx context.Context, notificationId string) (notification *as.AutoScalingNotification, errRet error) {
	logId := getLogId(ctx)
	request := as.NewDescribeNotificationConfigurationsRequest()
	request.AutoScalingNotificationIds = []*string{&notificationId}
	response, err := me.client.UseAsClient().DescribeNotificationConfigurations(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.AutoScalingNotificationSet) < 1 {
		errRet = fmt.Errorf("notification id is not found")
		return
	}
	notification = response.Response.AutoScalingNotificationSet[0]
	return
}

func (me *AsService) DeleteNotification(ctx context.Context, notificationId string) error {
	logId := getLogId(ctx)
	request := as.NewDeleteNotificationConfigurationRequest()
	request.AutoScalingNotificationId = &notificationId
	response, err := me.client.UseAsClient().DeleteNotificationConfiguration(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	return nil
}

func flattenDataDiskMappings(list []*as.DataDisk) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(list))
	for _, v := range list {
		disk := map[string]interface{}{
			"disk_type": *v.DiskType,
			"disk_size": *v.DiskSize,
		}
		result = append(result, disk)
	}
	return result
}

func flattenInstanceTagsMapping(list []*as.InstanceTag) map[string]interface{} {
	result := make(map[string]interface{}, len(list))
	for _, v := range list {
		result[*v.Key] = *v.Value
	}
	return result
}
