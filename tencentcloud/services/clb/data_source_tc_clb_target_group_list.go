package clb

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudClbTargetGroupList() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudClbTargetGroupListRead,
		Schema: map[string]*schema.Schema{
			"target_group_ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Target group ID array.",
			},

			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Filter array, which is exclusive of TargetGroupIds. Valid values: TargetGroupVpcId and TargetGroupName. Target group ID will be used first.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Filter name.",
						},
						"values": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "Filter value array.",
						},
					},
				},
			},

			"target_group_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Information set of displayed target groups.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"target_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Target group ID.",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "vpcid of target group.",
						},
						"target_group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Target group name.",
						},
						"port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Default port of target group. Note: this field may return null, indicating that no valid values can be obtained.",
						},
						"created_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Target group creation time.",
						},
						"updated_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Target group modification time.",
						},
						"associated_rule": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Array of associated rules. Note: this field may return null, indicating that no valid values can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"load_balancer_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ID of associated CLB instance.",
									},
									"listener_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ID of associated listener.",
									},
									"location_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ID of associated forwarding rule. Note: this field may return null, indicating that no valid values can be obtained.",
									},
									"protocol": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Protocol type of associated listener, such as HTTP or TCP.",
									},
									"port": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Port of associated listener.",
									},
									"domain": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Domain name of associated forwarding rule. Note: this field may return null, indicating that no valid values can be obtained.",
									},
									"url": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "URL of associated forwarding rule. Note: this field may return null, indicating that no valid values can be obtained.",
									},
									"load_balancer_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "CLB instance name.",
									},
									"listener_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Listener name.",
									},
								},
							},
						},
					},
				},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudClbTargetGroupListRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_clb_target_group_list.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("target_group_ids"); ok {
		targetGroupIdsSet := v.(*schema.Set).List()
		paramMap["TargetGroupIds"] = helper.InterfacesStringsPoint(targetGroupIdsSet)
	}

	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*clb.Filter, 0, len(filtersSet))

		for _, item := range filtersSet {
			filter := clb.Filter{}
			filterMap := item.(map[string]interface{})

			if v, ok := filterMap["name"]; ok {
				filter.Name = helper.String(v.(string))
			}
			if v, ok := filterMap["values"]; ok {
				valuesSet := v.(*schema.Set).List()
				filter.Values = helper.InterfacesStringsPoint(valuesSet)
			}
			tmpSet = append(tmpSet, &filter)
		}
		paramMap["Filters"] = tmpSet
	}

	service := ClbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var targetGroupSet []*clb.TargetGroupInfo

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeClbTargetGroupListByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		targetGroupSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(targetGroupSet))
	tmpList := make([]map[string]interface{}, 0, len(targetGroupSet))

	if targetGroupSet != nil {
		for _, targetGroupInfo := range targetGroupSet {
			targetGroupInfoMap := map[string]interface{}{}

			if targetGroupInfo.TargetGroupId != nil {
				targetGroupInfoMap["target_group_id"] = targetGroupInfo.TargetGroupId
			}

			if targetGroupInfo.VpcId != nil {
				targetGroupInfoMap["vpc_id"] = targetGroupInfo.VpcId
			}

			if targetGroupInfo.TargetGroupName != nil {
				targetGroupInfoMap["target_group_name"] = targetGroupInfo.TargetGroupName
			}

			if targetGroupInfo.Port != nil {
				targetGroupInfoMap["port"] = targetGroupInfo.Port
			}

			if targetGroupInfo.CreatedTime != nil {
				targetGroupInfoMap["created_time"] = targetGroupInfo.CreatedTime
			}

			if targetGroupInfo.UpdatedTime != nil {
				targetGroupInfoMap["updated_time"] = targetGroupInfo.UpdatedTime
			}

			if targetGroupInfo.AssociatedRule != nil {
				associatedRuleList := []interface{}{}
				for _, associatedRule := range targetGroupInfo.AssociatedRule {
					associatedRuleMap := map[string]interface{}{}

					if associatedRule.LoadBalancerId != nil {
						associatedRuleMap["load_balancer_id"] = associatedRule.LoadBalancerId
					}

					if associatedRule.ListenerId != nil {
						associatedRuleMap["listener_id"] = associatedRule.ListenerId
					}

					if associatedRule.LocationId != nil {
						associatedRuleMap["location_id"] = associatedRule.LocationId
					}

					if associatedRule.Protocol != nil {
						associatedRuleMap["protocol"] = associatedRule.Protocol
					}

					if associatedRule.Port != nil {
						associatedRuleMap["port"] = associatedRule.Port
					}

					if associatedRule.Domain != nil {
						associatedRuleMap["domain"] = associatedRule.Domain
					}

					if associatedRule.Url != nil {
						associatedRuleMap["url"] = associatedRule.Url
					}

					if associatedRule.LoadBalancerName != nil {
						associatedRuleMap["load_balancer_name"] = associatedRule.LoadBalancerName
					}

					if associatedRule.ListenerName != nil {
						associatedRuleMap["listener_name"] = associatedRule.ListenerName
					}

					associatedRuleList = append(associatedRuleList, associatedRuleMap)
				}

				targetGroupInfoMap["associated_rule"] = []interface{}{associatedRuleList}
			}

			ids = append(ids, *targetGroupInfo.TargetGroupId)
			tmpList = append(tmpList, targetGroupInfoMap)
		}

		_ = d.Set("target_group_set", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
