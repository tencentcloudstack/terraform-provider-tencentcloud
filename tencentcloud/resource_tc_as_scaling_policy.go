/*
Provides a resource for an AS (Auto scaling) policy.

Example Usage

```hcl
resource "tencentcloud_as_scaling_policy" "scaling_policy" {
  scaling_group_id    = "asg-n32ymck2"
  policy_name         = "tf-as-scaling-policy"
  adjustment_type     = "EXACT_CAPACITY"
  adjustment_value    = 0
  comparison_operator = "GREATER_THAN"
  metric_name         = "CPU_UTILIZATION"
  threshold           = 80
  period              = 300
  continuous_time     = 10
  statistic           = "AVERAGE"
  cooldown            = 360
}
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	as "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/as/v20180419"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudAsScalingPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudAsScalingPolicyCreate,
		Read:   resourceTencentCloudAsScalingPolicyRead,
		Update: resourceTencentCloudAsScalingPolicyUpdate,
		Delete: resourceTencentCloudAsScalingPolicyDelete,

		Schema: map[string]*schema.Schema{
			"scaling_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of a scaling group.",
			},
			"policy_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of a policy used to define a reaction when an alarm is triggered.",
			},
			"adjustment_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAllowedStringValue(SCALING_GROUP_ADJUSTMENT_TYPE),
				Description:  "Specifies whether the adjustment is an absolute number or a percentage of the current capacity. Valid values: `CHANGE_IN_CAPACITY`, `EXACT_CAPACITY` and `PERCENT_CHANGE_IN_CAPACITY`.",
			},
			"adjustment_value": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Define the number of instances by which to scale.For `CHANGE_IN_CAPACITY` type or PERCENT_CHANGE_IN_CAPACITY, a positive increment adds to the current capacity and a negative value removes from the current capacity. For `EXACT_CAPACITY` type, it defines an absolute number of the existing Auto Scaling group size.",
			},
			"comparison_operator": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAllowedStringValue(SCALING_GROUP_COMPARISON_OPERATOR),
				Description:  "Comparison operator. Valid values: `GREATER_THAN`, `GREATER_THAN_OR_EQUAL_TO`, `LESS_THAN`, `LESS_THAN_OR_EQUAL_TO`, `EQUAL_TO` and `NOT_EQUAL_TO`.",
			},
			"metric_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAllowedStringValue(SCALING_GROUP_METRIC_NAME),
				Description:  "Name of an indicator. Valid values: `CPU_UTILIZATION`, `MEM_UTILIZATION`, `LAN_TRAFFIC_OUT`, `LAN_TRAFFIC_IN`, `WAN_TRAFFIC_OUT` and `WAN_TRAFFIC_IN`.",
			},
			"threshold": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Alarm threshold.",
			},
			"period": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validateAllowedIntValue([]int{60, 300}),
				Description:  "Time period in second. Valid values: `60` and `300`.",
			},
			"continuous_time": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validateIntegerInRange(1, 10),
				Description:  "Retry times. Valid value ranges: (1~10).",
			},
			"statistic": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      SCALING_GROUP_STATISTIC_AVERAGE,
				ValidateFunc: validateAllowedStringValue(SCALING_GROUP_STATISTIC),
				Description:  "Statistic types. Valid values: `AVERAGE`, `MAXIMUM` and `MINIMUM`. Default is `AVERAGE`.",
			},
			"cooldown": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     300,
				Description: "Cooldwon time in second. Default is `30`0.",
			},
			"notification_user_group_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "An ID group of users to be notified when an alarm is triggered.",
			},
		},
	}
}

func resourceTencentCloudAsScalingPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_as_scaling_policy.create")()

	logId := getLogId(contextNil)

	request := as.NewCreateScalingPolicyRequest()
	request.AutoScalingGroupId = helper.String(d.Get("scaling_group_id").(string))
	request.ScalingPolicyName = helper.String(d.Get("policy_name").(string))
	request.AdjustmentType = helper.String(d.Get("adjustment_type").(string))
	adjustMentValue := int64(d.Get("adjustment_value").(int))
	request.AdjustmentValue = &adjustMentValue
	request.MetricAlarm = &as.MetricAlarm{}
	request.MetricAlarm.ComparisonOperator = helper.String(d.Get("comparison_operator").(string))
	request.MetricAlarm.MetricName = helper.String(d.Get("metric_name").(string))
	request.MetricAlarm.Threshold = helper.IntUint64(d.Get("threshold").(int))
	request.MetricAlarm.Period = helper.IntUint64(d.Get("period").(int))
	request.MetricAlarm.ContinuousTime = helper.IntUint64(d.Get("continuous_time").(int))

	if v, ok := d.GetOk("statistic"); ok {
		request.MetricAlarm.Statistic = helper.String(v.(string))
	}
	if v, ok := d.GetOk("cooldown"); ok {
		request.Cooldown = helper.IntUint64(v.(int))
	}
	if v, ok := d.GetOk("notification_user_group_ids"); ok {
		notificationUserGroupIds := v.([]interface{})
		request.NotificationUserGroupIds = make([]*string, 0, len(notificationUserGroupIds))
		for _, value := range notificationUserGroupIds {
			request.NotificationUserGroupIds = append(request.NotificationUserGroupIds, helper.String(value.(string)))
		}
	}

	response, err := meta.(*TencentCloudClient).apiV3Conn.UseAsClient().CreateScalingPolicy(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response.AutoScalingPolicyId == nil {
		return fmt.Errorf("scaling policy id is nil")
	}
	d.SetId(*response.Response.AutoScalingPolicyId)

	return resourceTencentCloudAsScalingPolicyRead(d, meta)
}

func resourceTencentCloudAsScalingPolicyRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_as_scaling_policy.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	scalingPolicyId := d.Id()
	asService := AsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		scalingPolicy, has, e := asService.DescribeScalingPolicyById(ctx, scalingPolicyId)
		if e != nil {
			return retryError(e)
		}
		if has == 0 {
			d.SetId("")
			return nil
		}
		_ = d.Set("scaling_group_id", *scalingPolicy.AutoScalingGroupId)
		_ = d.Set("policy_name", *scalingPolicy.ScalingPolicyName)
		_ = d.Set("adjustment_type", *scalingPolicy.AdjustmentType)
		_ = d.Set("adjustment_value", *scalingPolicy.AdjustmentValue)
		_ = d.Set("comparison_operator", *scalingPolicy.MetricAlarm.ComparisonOperator)
		_ = d.Set("metric_name", *scalingPolicy.MetricAlarm.MetricName)
		_ = d.Set("threshold", *scalingPolicy.MetricAlarm.Threshold)
		_ = d.Set("period", *scalingPolicy.MetricAlarm.Period)
		_ = d.Set("continuous_time", *scalingPolicy.MetricAlarm.ContinuousTime)
		_ = d.Set("statistic", *scalingPolicy.MetricAlarm.Statistic)
		_ = d.Set("cooldown", *scalingPolicy.Cooldown)
		if scalingPolicy.NotificationUserGroupIds != nil {
			_ = d.Set("notification_user_group_ids", helper.StringsInterfaces(scalingPolicy.NotificationUserGroupIds))
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
func resourceTencentCloudAsScalingPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_as_scaling_policy.update")()

	logId := getLogId(contextNil)

	request := as.NewModifyScalingPolicyRequest()
	scalingPolicyId := d.Id()
	request.AutoScalingPolicyId = &scalingPolicyId
	if d.HasChange("policy_name") {
		request.ScalingPolicyName = helper.String(d.Get("policy_name").(string))
	}
	if d.HasChange("adjustment_type") {
		request.AdjustmentType = helper.String(d.Get("adjustment_type").(string))
	}
	if d.HasChange("adjustment_value") {
		adjustmentValue := int64(d.Get("adjustment_value").(int))
		request.AdjustmentValue = &adjustmentValue
	}
	request.MetricAlarm = &as.MetricAlarm{}

	//these two parameter must pass together
	request.MetricAlarm.ComparisonOperator = helper.String(d.Get("comparison_operator").(string))
	request.MetricAlarm.Threshold = helper.IntUint64(d.Get("threshold").(int))
	request.MetricAlarm.MetricName = helper.String(d.Get("metric_name").(string))
	request.MetricAlarm.Period = helper.IntUint64(d.Get("period").(int))
	request.MetricAlarm.ContinuousTime = helper.IntUint64(d.Get("continuous_time").(int))
	request.MetricAlarm.Statistic = helper.String(d.Get("statistic").(string))

	if d.HasChange("cooldown") {
		request.Cooldown = helper.IntUint64(d.Get("cooldown").(int))
	}
	if d.HasChange("notification_user_group_ids") {
		notificationUserGroupIds := d.Get("notification_user_group_ids").([]interface{})
		request.NotificationUserGroupIds = make([]*string, 0, len(notificationUserGroupIds))
		for _, value := range notificationUserGroupIds {
			request.NotificationUserGroupIds = append(request.NotificationUserGroupIds, helper.String(value.(string)))
		}
	}

	response, err := meta.(*TencentCloudClient).apiV3Conn.UseAsClient().ModifyScalingPolicy(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return nil
}

func resourceTencentCloudAsScalingPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_as_scaling_policy.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	scalingPolicyId := d.Id()
	asService := AsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	err := asService.DeleteScalingPolicy(ctx, scalingPolicyId)
	if err != nil {
		return err
	}

	return nil
}
