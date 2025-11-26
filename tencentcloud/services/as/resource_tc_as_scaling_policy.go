package as

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	as "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/as/v20180419"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudAsScalingPolicy() *schema.Resource {
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
				Optional:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(SCALING_GROUP_ADJUSTMENT_TYPE),
				Description:  "Specifies whether the adjustment is an absolute number or a percentage of the current capacity. Valid values: `CHANGE_IN_CAPACITY`, `EXACT_CAPACITY` and `PERCENT_CHANGE_IN_CAPACITY`.",
			},
			"adjustment_value": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Define the number of instances by which to scale.For `CHANGE_IN_CAPACITY` type or PERCENT_CHANGE_IN_CAPACITY, a positive increment adds to the current capacity and a negative value removes from the current capacity. For `EXACT_CAPACITY` type, it defines an absolute number of the existing Auto Scaling group size.",
			},
			"comparison_operator": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(SCALING_GROUP_COMPARISON_OPERATOR),
				Description:  "Comparison operator. Valid values: `GREATER_THAN`, `GREATER_THAN_OR_EQUAL_TO`, `LESS_THAN`, `LESS_THAN_OR_EQUAL_TO`, `EQUAL_TO` and `NOT_EQUAL_TO`.",
			},
			"metric_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(SCALING_GROUP_METRIC_NAME),
				Description:  "Name of an indicator. Valid values: `CPU_UTILIZATION`, `MEM_UTILIZATION`, `LAN_TRAFFIC_OUT`, `LAN_TRAFFIC_IN`, `WAN_TRAFFIC_OUT` and `WAN_TRAFFIC_IN`.",
			},
			"threshold": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Alarm threshold.",
			},
			"period": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: tccommon.ValidateAllowedIntValue([]int{60, 300}),
				Description:  "Time period in second. Valid values: `60` and `300`.",
			},
			"continuous_time": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: tccommon.ValidateIntegerInRange(1, 10),
				Description:  "Retry times. Valid value ranges: (1~10).",
			},
			"statistic": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(SCALING_GROUP_STATISTIC),
				Description:  "Statistic types. Valid values: `AVERAGE`, `MAXIMUM` and `MINIMUM`. Default is `AVERAGE`.",
			},
			"cooldown": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Cooldwon time in second. Default is `300`.",
			},
			"notification_user_group_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "An ID group of users to be notified when an alarm is triggered.",
			},
			"policy_type": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "Alarm triggering policy type, the default type is SIMPLE. Value range: SIMPLE: Simple policy; TARGET_TRACKING: Target tracking policy.",
			},
			"predefined_metric_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Predefined monitoring items, applicable only to target tracking policies, and required in target tracking policy scenarios. Value range: ASG_AVG_CPU_UTILIZATION: Average CPU utilization; ASG_AVG_LAN_TRAFFIC_OUT: Average intranet outbound bandwidth; ASG_AVG_LAN_TRAFFIC_IN: Average intranet inbound bandwidth; ASG_AVG_WAN_TRAFFIC_OUT: Average internet outbound bandwidth; ASG_AVG_WAN_TRAFFIC_IN: Average internet inbound bandwidth.",
			},
			"target_value": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Target value, applicable only to target tracking strategies, and required in target tracking strategy scenarios. ASG_AVG_CPU_UTILIZATION: [1, 100), Unit: %; ASG_AVG_LAN_TRAFFIC_OUT: >0, Unit: Mbps; ASG_AVG_LAN_TRAFFIC_IN: >0, Unit: Mbps; ASG_AVG_WAN_TRAFFIC_OUT: >0, Unit: Mbps; ASG_AVG_WAN_TRAFFIC_IN: >0, Unit: Mbps.",
			},
			"estimated_instance_warmup": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Instance warm-up time, in seconds, applicable only to target tracking strategies. Value range is 0-3600, with a default warm-up time of 300 seconds.",
			},
			"disable_scale_in": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to disable scaling down applies only to the target tracking strategy; the default value is false. Value range: true: The target tracking strategy only triggers scaling up; false: The target tracking strategy triggers both scaling up and scaling down.",
			},
		},
	}
}

func resourceTencentCloudAsScalingPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_as_scaling_policy.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := as.NewCreateScalingPolicyRequest()
	request.AutoScalingGroupId = helper.String(d.Get("scaling_group_id").(string))
	request.ScalingPolicyName = helper.String(d.Get("policy_name").(string))
	if v, ok := d.GetOk("adjustment_type"); ok {
		request.AdjustmentType = helper.String(v.(string))
	}
	if v, ok := d.GetOkExists("adjustment_value"); ok {
		request.AdjustmentValue = helper.IntInt64(v.(int))
	}
	metricAlarm := &as.MetricAlarm{}
	var hasMetricAlarm bool
	if v, ok := d.GetOk("comparison_operator"); ok {
		metricAlarm.ComparisonOperator = helper.String(v.(string))
		hasMetricAlarm = true
	}
	if v, ok := d.GetOk("metric_name"); ok {
		metricAlarm.MetricName = helper.String(v.(string))
		hasMetricAlarm = true
	}
	if v, ok := d.GetOkExists("threshold"); ok {
		metricAlarm.Threshold = helper.IntUint64(v.(int))
		hasMetricAlarm = true
	}
	if v, ok := d.GetOkExists("period"); ok {
		metricAlarm.Period = helper.IntUint64(v.(int))
		hasMetricAlarm = true
	}
	if v, ok := d.GetOkExists("continuous_time"); ok {
		metricAlarm.ContinuousTime = helper.IntUint64(v.(int))
		hasMetricAlarm = true
	}
	if v, ok := d.GetOk("statistic"); ok {
		metricAlarm.Statistic = helper.String(v.(string))
		hasMetricAlarm = true
	}
	if hasMetricAlarm {
		request.MetricAlarm = metricAlarm
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

	if v, ok := d.GetOk("policy_type"); ok {
		request.ScalingPolicyType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("predefined_metric_type"); ok {
		request.PredefinedMetricType = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("target_value"); ok {
		request.TargetValue = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("estimated_instance_warmup"); ok {
		request.EstimatedInstanceWarmup = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("disable_scale_in"); ok {
		request.DisableScaleIn = helper.Bool(v.(bool))
	}

	response, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseAsClient().CreateScalingPolicy(request)
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
	defer tccommon.LogElapsed("resource.tencentcloud_as_scaling_policy.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	scalingPolicyId := d.Id()
	asService := AsService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		scalingPolicy, has, e := asService.DescribeScalingPolicyById(ctx, scalingPolicyId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		if has == 0 {
			d.SetId("")
			return nil
		}

		if scalingPolicy.AutoScalingGroupId != nil {
			_ = d.Set("scaling_group_id", *scalingPolicy.AutoScalingGroupId)
		}
		if scalingPolicy.ScalingPolicyName != nil {
			_ = d.Set("policy_name", *scalingPolicy.ScalingPolicyName)
		}
		if scalingPolicy.AdjustmentType != nil {
			_ = d.Set("adjustment_type", *scalingPolicy.AdjustmentType)
		}
		if scalingPolicy.AdjustmentValue != nil {
			_ = d.Set("adjustment_value", *scalingPolicy.AdjustmentValue)
		}
		if scalingPolicy.MetricAlarm != nil {
			if scalingPolicy.MetricAlarm.ComparisonOperator != nil {
				_ = d.Set("comparison_operator", *scalingPolicy.MetricAlarm.ComparisonOperator)
			}
			if scalingPolicy.MetricAlarm.MetricName != nil {
				_ = d.Set("metric_name", *scalingPolicy.MetricAlarm.MetricName)
			}
			if scalingPolicy.MetricAlarm.Threshold != nil {
				_ = d.Set("threshold", *scalingPolicy.MetricAlarm.Threshold)
			}
			if scalingPolicy.MetricAlarm.Period != nil {
				_ = d.Set("period", *scalingPolicy.MetricAlarm.Period)
			}
			if scalingPolicy.MetricAlarm.ContinuousTime != nil {
				_ = d.Set("continuous_time", *scalingPolicy.MetricAlarm.ContinuousTime)
			}
			if scalingPolicy.MetricAlarm.Statistic != nil {
				_ = d.Set("statistic", *scalingPolicy.MetricAlarm.Statistic)
			}
		}
		if scalingPolicy.Cooldown != nil {
			_ = d.Set("cooldown", *scalingPolicy.Cooldown)
		}
		if scalingPolicy.NotificationUserGroupIds != nil {
			_ = d.Set("notification_user_group_ids", helper.StringsInterfaces(scalingPolicy.NotificationUserGroupIds))
		}
		if scalingPolicy.ScalingPolicyType != nil {
			_ = d.Set("policy_type", *scalingPolicy.ScalingPolicyType)
		}
		if scalingPolicy.PredefinedMetricType != nil {
			_ = d.Set("predefined_metric_type", *scalingPolicy.PredefinedMetricType)
		}
		if scalingPolicy.TargetValue != nil {
			_ = d.Set("target_value", *scalingPolicy.TargetValue)
		}
		if scalingPolicy.EstimatedInstanceWarmup != nil {
			_ = d.Set("estimated_instance_warmup", *scalingPolicy.EstimatedInstanceWarmup)
		}
		if scalingPolicy.DisableScaleIn != nil {
			_ = d.Set("disable_scale_in", *scalingPolicy.DisableScaleIn)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
func resourceTencentCloudAsScalingPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_as_scaling_policy.update")()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := as.NewModifyScalingPolicyRequest()
	scalingPolicyId := d.Id()
	request.AutoScalingPolicyId = &scalingPolicyId
	if d.HasChange("policy_name") {
		request.ScalingPolicyName = helper.String(d.Get("policy_name").(string))
	}
	if d.HasChange("adjustment_type") {
		if v, ok := d.GetOk("adjustment_type"); ok {
			request.AdjustmentType = helper.String(v.(string))
		}
	}
	if d.HasChange("adjustment_value") {
		if v, ok := d.GetOkExists("adjustment_value"); ok {
			request.AdjustmentValue = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("comparison_operator") || d.HasChange("threshold") || d.HasChange("metric_name") || d.HasChange("period") || d.HasChange("continuous_time") || d.HasChange("statistic") {
		request.MetricAlarm = &as.MetricAlarm{}

		if v, ok := d.GetOk("comparison_operator"); ok {
			request.MetricAlarm.ComparisonOperator = helper.String(v.(string))
		}

		if v, ok := d.GetOkExists("threshold"); ok {
			request.MetricAlarm.Threshold = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOk("metric_name"); ok {
			request.MetricAlarm.MetricName = helper.String(v.(string))
		}

		if v, ok := d.GetOkExists("period"); ok {
			request.MetricAlarm.Period = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOkExists("continuous_time"); ok {
			request.MetricAlarm.ContinuousTime = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOk("statistic"); ok {
			request.MetricAlarm.Statistic = helper.String(v.(string))
		}
	}

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

	if v, ok := d.GetOk("predefined_metric_type"); ok {
		request.PredefinedMetricType = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("target_value"); ok {
		request.TargetValue = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("estimated_instance_warmup"); ok {
		request.EstimatedInstanceWarmup = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("disable_scale_in"); ok {
		request.DisableScaleIn = helper.Bool(v.(bool))
	}

	response, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseAsClient().ModifyScalingPolicy(request)
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
	defer tccommon.LogElapsed("resource.tencentcloud_as_scaling_policy.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	scalingPolicyId := d.Id()
	asService := AsService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	err := asService.DeleteScalingPolicy(ctx, scalingPolicyId)
	if err != nil {
		return err
	}

	return nil
}
