package as

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	as "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/as/v20180419"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func NewAsService(client *connectivity.TencentCloudClient) AsService {
	return AsService{client: client}
}

type AsService struct {
	client *connectivity.TencentCloudClient
}

func (me *AsService) DescribeLaunchConfigurationById(ctx context.Context, configurationId string) (config *as.LaunchConfiguration, has int, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := as.NewDescribeLaunchConfigurationsRequest()
	request.LaunchConfigurationIds = []*string{&configurationId}
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseAsClient().DescribeLaunchConfigurations(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	has = len(response.Response.LaunchConfigurationSet)
	if has < 1 {
		return
	}
	config = response.Response.LaunchConfigurationSet[0]
	return
}

func (me *AsService) DescribeLaunchConfigurationByFilter(ctx context.Context, configurationId, configurationName string) (configs []*as.LaunchConfiguration, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := as.NewDescribeLaunchConfigurationsRequest()
	filters := make([]*as.Filter, 0)
	if configurationId != "" {
		filter := &as.Filter{
			Name:   helper.String("launch-configuration-id"),
			Values: []*string{&configurationId},
		}
		filters = append(filters, filter)
	}
	if configurationName != "" {
		filter := &as.Filter{
			Name:   helper.String("launch-configuration-name"),
			Values: []*string{&configurationName},
		}
		filters = append(filters, filter)
	}
	if len(filters) > 0 {
		request.Filters = filters
	}

	offset := 0
	pageSize := 100
	configs = make([]*as.LaunchConfiguration, 0)
	for {
		request.Offset = helper.IntUint64(offset)
		request.Limit = helper.IntUint64(pageSize)
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseAsClient().DescribeLaunchConfigurations(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
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
	logId := tccommon.GetLogId(ctx)
	request := as.NewDeleteLaunchConfigurationRequest()
	request.LaunchConfigurationId = &configurationId
	ratelimit.Check(request.GetAction())
	_, err := me.client.UseAsClient().DeleteLaunchConfiguration(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	return nil
}

func (me *AsService) DescribeAutoScalingGroupById(ctx context.Context, scalingGroupId string) (scalingGroup *as.AutoScalingGroup, has int, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := as.NewDescribeAutoScalingGroupsRequest()
	request.AutoScalingGroupIds = []*string{&scalingGroupId}
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseAsClient().DescribeAutoScalingGroups(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	if response == nil || response.Response == nil || response.Response.AutoScalingGroupSet == nil {
		return
	}
	has = len(response.Response.AutoScalingGroupSet)
	if has < 1 {
		return
	}
	scalingGroup = response.Response.AutoScalingGroupSet[0]
	return
}

func (me *AsService) DescribeAutoScalingGroupByFilter(
	ctx context.Context,
	scalingGroupId, configurationId, scalingGroupName string,
	tags map[string]string,
) (scalingGroups []*as.AutoScalingGroup, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := as.NewDescribeAutoScalingGroupsRequest()
	filters := make([]*as.Filter, 0)
	if scalingGroupId != "" {
		filter := &as.Filter{
			Name:   helper.String("auto-scaling-group-id"),
			Values: []*string{&scalingGroupId},
		}
		filters = append(filters, filter)
	}
	if configurationId != "" {
		filter := &as.Filter{
			Name:   helper.String("launch-configuration-id"),
			Values: []*string{&configurationId},
		}
		filters = append(filters, filter)
	}
	if scalingGroupName != "" {
		filter := &as.Filter{
			Name:   helper.String("auto-scaling-group-name"),
			Values: []*string{&scalingGroupName},
		}
		filters = append(filters, filter)
	}
	for k, v := range tags {
		filters = append(filters, &as.Filter{
			Name:   helper.String("tag:" + k),
			Values: []*string{helper.String(v)},
		})
	}
	if len(filters) > 0 {
		request.Filters = filters
	}

	offset := 0
	pageSize := 100
	scalingGroups = make([]*as.AutoScalingGroup, 0)
	for {
		request.Offset = helper.IntUint64(offset)
		request.Limit = helper.IntUint64(pageSize)
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseAsClient().DescribeAutoScalingGroups(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
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
	logId := tccommon.GetLogId(ctx)
	request := as.NewModifyAutoScalingGroupRequest()
	request.AutoScalingGroupId = &scalingGroupId
	request.MinSize = helper.IntUint64(0)
	request.MaxSize = helper.IntUint64(0)
	request.DesiredCapacity = helper.IntUint64(0)
	ratelimit.Check(request.GetAction())
	_, err := me.client.UseAsClient().ModifyAutoScalingGroup(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	return nil
}

func (me *AsService) DeleteScalingGroup(ctx context.Context, scalingGroupId string) error {
	logId := tccommon.GetLogId(ctx)
	request := as.NewDeleteAutoScalingGroupRequest()
	request.AutoScalingGroupId = &scalingGroupId
	ratelimit.Check(request.GetAction())
	_, err := me.client.UseAsClient().DeleteAutoScalingGroup(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	return nil
}

func (me *AsService) AttachInstances(ctx context.Context, scalingGroupId string, instanceIds []string) error {
	logId := tccommon.GetLogId(ctx)
	request := as.NewAttachInstancesRequest()
	request.AutoScalingGroupId = &scalingGroupId
	request.InstanceIds = make([]*string, 0, len(instanceIds))
	for i := range instanceIds {
		request.InstanceIds = append(request.InstanceIds, &instanceIds[i])
	}
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseAsClient().AttachInstances(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	activityId := *response.Response.ActivityId

	err = resource.Retry(4*tccommon.ReadRetryTimeout, func() *resource.RetryError {
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
	logId := tccommon.GetLogId(ctx)
	request := as.NewDescribeAutoScalingActivitiesRequest()
	request.ActivityIds = []*string{&activityId}
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseAsClient().DescribeAutoScalingActivities(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	if len(response.Response.ActivitySet) < 1 {
		errRet = fmt.Errorf("activity id set is nil")
	}
	status = *response.Response.ActivitySet[0].StatusCode
	return
}

func (me *AsService) DetachInstances(ctx context.Context, scalingGroupId string, instanceIds []string) error {
	logId := tccommon.GetLogId(ctx)
	request := as.NewDetachInstancesRequest()
	request.AutoScalingGroupId = &scalingGroupId
	request.InstanceIds = make([]*string, 0, len(instanceIds))
	for i := range instanceIds {
		request.InstanceIds = append(request.InstanceIds, &instanceIds[i])
	}
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseAsClient().DetachInstances(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	activityId := *response.Response.ActivityId

	err = resource.Retry(4*tccommon.ReadRetryTimeout, func() *resource.RetryError {
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

func (me *AsService) DescribeAutoScalingAttachment(ctx context.Context, scalingGroupId string, fully bool) (instanceIds []string, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := as.NewDescribeAutoScalingInstancesRequest()
	request.Filters = []*as.Filter{
		{
			Name:   helper.String("auto-scaling-group-id"),
			Values: []*string{&scalingGroupId},
		},
	}
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseAsClient().DescribeAutoScalingInstances(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}

	instanceIds = make([]string, 0)
	for _, instance := range response.Response.AutoScalingInstanceSet {
		if *instance.CreationType == "MANUAL_ATTACHING" || fully {
			instanceIds = append(instanceIds, *instance.InstanceId)
		}
	}
	return
}

func (me *AsService) DescribeScalingPolicyById(ctx context.Context, scalingPolicyId string) (scalingPolicy *as.ScalingPolicy, has int, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := as.NewDescribeScalingPoliciesRequest()
	request.AutoScalingPolicyIds = []*string{&scalingPolicyId}
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseAsClient().DescribeScalingPolicies(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	has = len(response.Response.ScalingPolicySet)
	if has < 1 {
		return
	}
	scalingPolicy = response.Response.ScalingPolicySet[0]
	return
}

func (me *AsService) DescribeScalingPolicyByFilter(ctx context.Context, policyId, policyName, scalingGroupId string) (scalingPolicies []*as.ScalingPolicy, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := as.NewDescribeScalingPoliciesRequest()
	request.Filters = make([]*as.Filter, 0)
	if policyId != "" {
		filter := &as.Filter{
			Name:   helper.String("auto-scaling-policy-id"),
			Values: []*string{&policyId},
		}
		request.Filters = append(request.Filters, filter)
	}
	if policyName != "" {
		filter := &as.Filter{
			Name:   helper.String("scaling-policy-name"),
			Values: []*string{&policyName},
		}
		request.Filters = append(request.Filters, filter)
	}
	if scalingGroupId != "" {
		filter := &as.Filter{
			Name:   helper.String("auto-scaling-group-id"),
			Values: []*string{&scalingGroupId},
		}
		request.Filters = append(request.Filters, filter)
	}

	offset := 0
	pageSize := 100
	scalingPolicies = make([]*as.ScalingPolicy, 0)
	for {
		request.Offset = helper.IntUint64(offset)
		request.Limit = helper.IntUint64(pageSize)
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseAsClient().DescribeScalingPolicies(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
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
	logId := tccommon.GetLogId(ctx)
	request := as.NewDeleteScalingPolicyRequest()
	request.AutoScalingPolicyId = &scalingPolicyId
	ratelimit.Check(request.GetAction())
	_, err := me.client.UseAsClient().DeleteScalingPolicy(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	return nil
}

func (me *AsService) DescribeScheduledActionById(ctx context.Context, scheduledActionId string) (scheduledAction *as.ScheduledAction, has int, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := as.NewDescribeScheduledActionsRequest()
	request.ScheduledActionIds = []*string{&scheduledActionId}
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseAsClient().DescribeScheduledActions(request)
	if err != nil {
		sdkErr, ok := err.(*sdkErrors.TencentCloudSDKError)
		if ok && sdkErr.Code == AsScheduleNotFound {
			has = 0
			return
		}
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}

	has = len(response.Response.ScheduledActionSet)
	if has < 1 {
		return
	}
	scheduledAction = response.Response.ScheduledActionSet[0]
	return
}

func (me *AsService) ModifyAutoScalingGroup(ctx context.Context, request *as.ModifyAutoScalingGroupRequest) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseAsClient().ModifyAutoScalingGroup(request)

	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *AsService) DeleteScheduledAction(ctx context.Context, scheduledActonId string) error {
	logId := tccommon.GetLogId(ctx)
	request := as.NewDeleteScheduledActionRequest()
	request.ScheduledActionId = &scheduledActonId
	ratelimit.Check(request.GetAction())
	_, err := me.client.UseAsClient().DeleteScheduledAction(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	return nil
}

func (me *AsService) DescribeLifecycleHookById(ctx context.Context, lifecycleHookId string) (lifecycleHook *as.LifecycleHook, has int, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := as.NewDescribeLifecycleHooksRequest()
	request.LifecycleHookIds = []*string{&lifecycleHookId}
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseAsClient().DescribeLifecycleHooks(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	has = len(response.Response.LifecycleHookSet)
	if has < 1 {
		return
	}
	lifecycleHook = response.Response.LifecycleHookSet[0]
	return
}

func (me *AsService) DeleteLifecycleHook(ctx context.Context, lifecycleHookId string) error {
	logId := tccommon.GetLogId(ctx)
	request := as.NewDeleteLifecycleHookRequest()
	request.LifecycleHookId = &lifecycleHookId
	ratelimit.Check(request.GetAction())
	_, err := me.client.UseAsClient().DeleteLifecycleHook(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	return nil
}

func (me *AsService) DescribeNotificationById(ctx context.Context, notificationId string) (notification *as.AutoScalingNotification, has int, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := as.NewDescribeNotificationConfigurationsRequest()
	request.AutoScalingNotificationIds = []*string{&notificationId}
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseAsClient().DescribeNotificationConfigurations(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	has = len(response.Response.AutoScalingNotificationSet)
	if has < 1 {
		return
	}
	notification = response.Response.AutoScalingNotificationSet[0]
	return
}

func (me *AsService) DeleteNotification(ctx context.Context, notificationId string) error {
	logId := tccommon.GetLogId(ctx)
	request := as.NewDeleteNotificationConfigurationRequest()
	request.AutoScalingNotificationId = &notificationId
	ratelimit.Check(request.GetAction())
	_, err := me.client.UseAsClient().DeleteNotificationConfiguration(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	return nil
}

func (me *AsService) DescribeAsInstancesByFilter(ctx context.Context, param map[string]interface{}) (instances []*as.Instance, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = as.NewDescribeAutoScalingInstancesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "instance_ids" {
			request.InstanceIds = helper.Strings(v.([]string))
		}
		if k == "filters" {
			request.Filters = v.([]*as.Filter)
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
		response, err := me.client.UseAsClient().DescribeAutoScalingInstances(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.AutoScalingInstanceSet) < 1 {
			break
		}
		instances = append(instances, response.Response.AutoScalingInstanceSet...)
		if len(response.Response.AutoScalingInstanceSet) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func flattenDataDiskMappings(list []*as.DataDisk) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(list))
	for _, v := range list {
		disk := map[string]interface{}{
			"disk_size": *v.DiskSize,
		}
		if v.DiskType != nil {
			disk["disk_type"] = *v.DiskType
		}
		if v.SnapshotId != nil {
			disk["snapshot_id"] = *v.SnapshotId
		}
		if v.DeleteWithInstance != nil {
			disk["delete_with_instance"] = *v.DeleteWithInstance
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

func (me *AsService) DescribeAsAdvices(ctx context.Context, param map[string]interface{}) (advices []*as.AutoScalingAdvice, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = as.NewDescribeAutoScalingAdvicesRequest()
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

	response, err := me.client.UseAsClient().DescribeAutoScalingAdvices(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	advices = response.Response.AutoScalingAdviceSet

	return
}

func (me *AsService) DescribeAsLimits(ctx context.Context) (limits *as.DescribeAccountLimitsResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = as.NewDescribeAccountLimitsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseAsClient().DescribeAccountLimits(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	limits = response.Response

	return
}

func (me *AsService) DescribeAsLastActivity(ctx context.Context, param map[string]interface{}) (lastActivity []*as.Activity, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = as.NewDescribeAutoScalingGroupLastActivitiesRequest()
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

	response, err := me.client.UseAsClient().DescribeAutoScalingGroupLastActivities(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	lastActivity = response.Response.ActivitySet

	return
}

func (me *AsService) DescribeAsLoadBalancerById(ctx context.Context, autoScalingGroupId string) (loadBalancer *as.AutoScalingGroup, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := as.NewDescribeAutoScalingGroupsRequest()
	response := as.NewDescribeAutoScalingGroupsResponse()
	request.AutoScalingGroupIds = []*string{&autoScalingGroupId}
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseAsClient().DescribeAutoScalingGroups(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	if len(response.Response.AutoScalingGroupSet) < 1 {
		return
	}

	loadBalancer = response.Response.AutoScalingGroupSet[0]
	return
}

func (me *AsService) DeleteAsLoadBalancerById(ctx context.Context, autoScalingGroupId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := as.NewModifyLoadBalancersRequest()
	request.AutoScalingGroupId = &autoScalingGroupId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseAsClient().ModifyLoadBalancers(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}
