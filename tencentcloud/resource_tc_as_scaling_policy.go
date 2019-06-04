package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	as "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/as/v20180419"
)

func resourceTencentCloudAsScalingPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudAsScalingPolicyCreate,
		Read:   resourceTencentCloudAsScalingPolicyRead,
		Update: resourceTencentCloudAsScalingPolicyUpdate,
		Delete: resourceTencentCloudAsScalingPolicyDelete,

		Schema: map[string]*schema.Schema{
			"scaling_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"policy_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"adjustment_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAllowedStringValue(SCALING_GROUP_ADJUSTMENT_TYPE),
			},
			"adjustment_value": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"comparison_operator": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAllowedStringValue(SCALING_GROUP_COMPARISON_OPERATOR),
			},
			"metric_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAllowedStringValue(SCALING_GROUP_METRIC_NAME),
			},
			"threshold": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"period": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validateAllowedIntValue([]int{60, 300}),
			},
			"continuous_time": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validateIntegerInRange(1, 10),
			},
			"statistic": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      SCALING_GROUP_STATISTIC_AVERAGE,
				ValidateFunc: validateAllowedStringValue(SCALING_GROUP_STATISTIC),
			},
			"cooldown": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  300,
			},
			"notification_user_group_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceTencentCloudAsScalingPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	logId := GetLogId(nil)

	request := as.NewCreateScalingPolicyRequest()
	request.AutoScalingGroupId = stringToPointer(d.Get("scaling_group_id").(string))
	request.ScalingPolicyName = stringToPointer(d.Get("policy_name").(string))
	request.AdjustmentType = stringToPointer(d.Get("adjustment_type").(string))
	adjustMentValue := int64(d.Get("adjustment_value").(int))
	request.AdjustmentValue = &adjustMentValue
	request.MetricAlarm = &as.MetricAlarm{}
	request.MetricAlarm.ComparisonOperator = stringToPointer(d.Get("comparison_operator").(string))
	request.MetricAlarm.MetricName = stringToPointer(d.Get("metric_name").(string))
	request.MetricAlarm.Threshold = intToPointer(d.Get("threshold").(int))
	request.MetricAlarm.Period = intToPointer(d.Get("period").(int))
	request.MetricAlarm.ContinuousTime = intToPointer(d.Get("continuous_time").(int))

	if v, ok := d.GetOk("statistic"); ok {
		request.MetricAlarm.Statistic = stringToPointer(v.(string))
	}
	if v, ok := d.GetOk("cooldown"); ok {
		request.Cooldown = intToPointer(v.(int))
	}
	if v, ok := d.GetOk("notification_user_group_ids"); ok {
		notificationUserGroupIds := v.([]interface{})
		request.NotificationUserGroupIds = make([]*string, 0, len(notificationUserGroupIds))
		for _, value := range notificationUserGroupIds {
			request.NotificationUserGroupIds = append(request.NotificationUserGroupIds, stringToPointer(value.(string)))
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
	logId := GetLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	scalingPolicyId := d.Id()
	asService := AsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	scalingPolicy, err := asService.DescribeScalingPolicyById(ctx, scalingPolicyId)
	if err != nil {
		return err
	}

	d.Set("scaling_group_id", *scalingPolicy.AutoScalingGroupId)
	d.Set("policy_name", *scalingPolicy.ScalingPolicyName)
	d.Set("adjustment_type", *scalingPolicy.AdjustmentType)
	d.Set("adjustment_value", *scalingPolicy.AdjustmentValue)
	d.Set("comparison_operator", *scalingPolicy.MetricAlarm.ComparisonOperator)
	d.Set("metric_name", *scalingPolicy.MetricAlarm.MetricName)
	d.Set("threshold", *scalingPolicy.MetricAlarm.Threshold)
	d.Set("period", *scalingPolicy.MetricAlarm.Period)
	d.Set("continuous_time", *scalingPolicy.MetricAlarm.ContinuousTime)
	d.Set("statistic", *scalingPolicy.MetricAlarm.Statistic)
	d.Set("cooldown", *scalingPolicy.Cooldown)
	if scalingPolicy.NotificationUserGroupIds != nil {
		d.Set("notification_user_group_ids", flattenStringList(scalingPolicy.NotificationUserGroupIds))
	}

	return nil
}

func resourceTencentCloudAsScalingPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	logId := GetLogId(nil)

	request := as.NewModifyScalingPolicyRequest()
	scalingPolicyId := d.Id()
	request.AutoScalingPolicyId = &scalingPolicyId
	if d.HasChange("policy_name") {
		request.ScalingPolicyName = stringToPointer(d.Get("policy_name").(string))
	}
	if d.HasChange("adjustment_type") {
		request.AdjustmentType = stringToPointer(d.Get("adjustment_type").(string))
	}
	if d.HasChange("adjustment_value") {
		adjustmentValue := int64(d.Get("adjustment_value").(int))
		request.AdjustmentValue = &adjustmentValue
	}
	request.MetricAlarm = &as.MetricAlarm{}
	if d.HasChange("comparison_operator") {
		request.MetricAlarm.ComparisonOperator = stringToPointer(d.Get("comparison_operator").(string))
	}
	if d.HasChange("metric_name") {
		request.MetricAlarm.MetricName = stringToPointer(d.Get("metric_name").(string))
	}
	if d.HasChange("threshold") {
		request.MetricAlarm.Threshold = intToPointer(d.Get("threshold").(int))
	}
	if d.HasChange("period") {
		request.MetricAlarm.Period = intToPointer(d.Get("period").(int))
	}
	if d.HasChange("continuous_time") {
		request.MetricAlarm.ContinuousTime = intToPointer(d.Get("continuous_time").(int))
	}
	if d.HasChange("statistic") {
		request.MetricAlarm.Statistic = stringToPointer(d.Get("statistic").(string))
	}
	if d.HasChange("cooldown") {
		request.Cooldown = intToPointer(d.Get("cooldown").(int))
	}
	if d.HasChange("notification_user_group_ids") {
		notificationUserGroupIds := d.Get("notification_user_group_ids").([]interface{})
		request.NotificationUserGroupIds = make([]*string, 0, len(notificationUserGroupIds))
		for _, value := range notificationUserGroupIds {
			request.NotificationUserGroupIds = append(request.NotificationUserGroupIds, stringToPointer(value.(string)))
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
	logId := GetLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

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
