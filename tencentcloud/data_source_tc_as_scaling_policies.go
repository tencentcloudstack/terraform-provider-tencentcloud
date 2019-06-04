package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceTencentCloudAsScalingPolicies() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudAsScalingPolicyRead,

		Schema: map[string]*schema.Schema{
			"scaling_policy_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"scaling_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"policy_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"result_output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"scaling_policy_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"scaling_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policy_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"adjustment_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"adjustment_value": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"comparison_operator": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"metric_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"threshold": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"period": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"continuous_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"statistic": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cooldown": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"notification_user_group_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudAsScalingPolicyRead(d *schema.ResourceData, meta interface{}) error {
	logId := GetLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	asService := AsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	scalingPolicyId := ""
	scalingGroupId := ""
	policyName := ""
	if v, ok := d.GetOk("scaling_policy_id"); ok {
		scalingPolicyId = v.(string)
	}
	if v, ok := d.GetOk("scaling_group_id"); ok {
		scalingGroupId = v.(string)
	}
	if v, ok := d.GetOk("policy_name"); ok {
		policyName = v.(string)
	}

	scalingPolicies, err := asService.DescribeScalingPolicyByFilter(ctx, scalingPolicyId, policyName, scalingGroupId)
	if err != nil {
		return nil
	}

	scalingPolicyList := make([]map[string]interface{}, 0, len(scalingPolicies))
	for _, scalingPolicy := range scalingPolicies {
		mapping := map[string]interface{}{
			"scaling_group_id":            *scalingPolicy.AutoScalingGroupId,
			"policy_name":                 *scalingPolicy.ScalingPolicyName,
			"adjustment_type":             *scalingPolicy.AdjustmentType,
			"adjustment_value":            *scalingPolicy.AdjustmentValue,
			"comparison_operator":         *scalingPolicy.MetricAlarm.ComparisonOperator,
			"metric_name":                 *scalingPolicy.MetricAlarm.MetricName,
			"threshold":                   *scalingPolicy.MetricAlarm.Threshold,
			"period":                      *scalingPolicy.MetricAlarm.Period,
			"continuous_time":             *scalingPolicy.MetricAlarm.ContinuousTime,
			"statistic":                   *scalingPolicy.MetricAlarm.Statistic,
			"cooldown":                    *scalingPolicy.Cooldown,
			"notification_user_group_ids": flattenStringList(scalingPolicy.NotificationUserGroupIds),
		}
		scalingPolicyList = append(scalingPolicyList, mapping)
	}
	d.SetId("ScalingPolicyList" + scalingGroupId + scalingGroupId + policyName)
	err = d.Set("scaling_policy_list", scalingPolicyList)
	if err != nil {
		log.Printf("[CRITAL]%s provider set configuration list fail, reason:%s\n ", logId, err.Error())
		return err
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		writeToFile(output.(string), scalingPolicyList)
	}

	return nil
}
