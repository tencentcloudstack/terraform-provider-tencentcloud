package as

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudAsScalingPolicies() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudAsScalingPolicyRead,

		Schema: map[string]*schema.Schema{
			"scaling_policy_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Scaling policy ID.",
			},
			"scaling_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Scaling group ID.",
			},
			"policy_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Scaling policy name.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
			"scaling_policy_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of scaling policy. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"scaling_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Scaling policy ID.",
						},
						"policy_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Scaling policy name.",
						},
						"adjustment_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Adjustment type of the scaling rule.",
						},
						"adjustment_value": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Adjustment value of the scaling rule.",
						},
						"comparison_operator": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Comparison operator.",
						},
						"metric_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of an indicator.",
						},
						"threshold": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Alarm threshold.",
						},
						"period": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Time period in second.",
						},
						"continuous_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Retry times.",
						},
						"statistic": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Statistic types.",
						},
						"cooldown": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Cool down time of the scaling rule.",
						},
						"notification_user_group_ids": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Users need to be notified when an alarm is triggered.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudAsScalingPolicyRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_as_scaling_policies.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	asService := AsService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
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
		return err
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
			"notification_user_group_ids": helper.StringsInterfaces(scalingPolicy.NotificationUserGroupIds),
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
		if err = tccommon.WriteToFile(output.(string), scalingPolicyList); err != nil {
			return err
		}
	}

	return nil
}
