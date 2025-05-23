package privatedns

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	privatednsv20201028 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/privatedns/v20201028"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudPrivateDnsForwardRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudPrivateDnsForwardRulesRead,
		Schema: map[string]*schema.Schema{
			"filters": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Filter parameters.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Parameter name.",
						},
						"values": {
							Type:        schema.TypeSet,
							Required:    true,
							Description: "Array of parameter values.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},

			"forward_rule_set": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Private domain list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Private domain name.",
						},
						"rule_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Forwarding rule name.",
						},
						"rule_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Rule ID.",
						},
						"rule_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Forwarding rule type. DOWN: From cloud to off-cloud; UP: From off-cloud to cloud.",
						},
						"created_at": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Creation time.",
						},
						"updated_at": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Update time.",
						},
						"end_point_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Endpoint name.",
						},
						"end_point_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Endpoint ID.",
						},
						"forward_address": {
							Type:        schema.TypeSet,
							Required:    true,
							Description: "Forwarding address.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"vpc_set": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "List of VPCs bound to the private domain.\nNote: This field may return null, indicating that no valid values can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"uniq_vpc_id": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "VpcId: vpc-xadsafsdasd.",
									},
									"region": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "VPC region: ap-guangzhou, ap-shanghai.",
									},
								},
							},
						},
						"zone_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "ID of the bound private domain.",
						},
						"tags": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Tag.\nNote: This field may return null, indicating that no valid values can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"tag_key": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Tag key.",
									},
									"tag_value": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Tag value.",
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

func dataSourceTencentCloudPrivateDnsForwardRulesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_private_dns_forward_rules.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = PrivatednsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*privatednsv20201028.Filter, 0, len(filtersSet))
		for _, item := range filtersSet {
			filtersMap := item.(map[string]interface{})
			filter := privatednsv20201028.Filter{}
			if v, ok := filtersMap["name"]; ok {
				filter.Name = helper.String(v.(string))
			}

			if v, ok := filtersMap["values"]; ok {
				valuesSet := v.(*schema.Set).List()
				for i := range valuesSet {
					values := valuesSet[i].(string)
					filter.Values = append(filter.Values, helper.String(values))
				}
			}

			tmpSet = append(tmpSet, &filter)
		}

		paramMap["Filters"] = tmpSet
	}

	var respData []*privatednsv20201028.ForwardRule
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribePrivateDnsForwardRulesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e, PRIVATEDNS_CUSTOM_RETRY_SDK_ERROR...)
		}

		respData = result
		return nil
	})

	if err != nil {
		return err
	}

	var ids []string
	forwardRuleSetList := make([]map[string]interface{}, 0, len(respData))
	if respData != nil {
		for _, forwardRuleSet := range respData {
			forwardRuleSetMap := map[string]interface{}{}
			if forwardRuleSet.Domain != nil {
				forwardRuleSetMap["domain"] = forwardRuleSet.Domain
			}

			if forwardRuleSet.RuleName != nil {
				forwardRuleSetMap["rule_name"] = forwardRuleSet.RuleName
			}

			if forwardRuleSet.RuleId != nil {
				forwardRuleSetMap["rule_id"] = forwardRuleSet.RuleId
				ids = append(ids, *forwardRuleSet.RuleId)
			}

			if forwardRuleSet.RuleType != nil {
				forwardRuleSetMap["rule_type"] = forwardRuleSet.RuleType
			}

			if forwardRuleSet.CreatedAt != nil {
				forwardRuleSetMap["created_at"] = forwardRuleSet.CreatedAt
			}

			if forwardRuleSet.UpdatedAt != nil {
				forwardRuleSetMap["updated_at"] = forwardRuleSet.UpdatedAt
			}

			if forwardRuleSet.EndPointName != nil {
				forwardRuleSetMap["end_point_name"] = forwardRuleSet.EndPointName
			}

			if forwardRuleSet.EndPointId != nil {
				forwardRuleSetMap["end_point_id"] = forwardRuleSet.EndPointId
			}

			if forwardRuleSet.ForwardAddress != nil {
				forwardRuleSetMap["forward_address"] = forwardRuleSet.ForwardAddress
			}

			vpcSetList := make([]map[string]interface{}, 0, len(forwardRuleSet.VpcSet))
			if forwardRuleSet.VpcSet != nil {
				for _, vpcSet := range forwardRuleSet.VpcSet {
					vpcSetMap := map[string]interface{}{}

					if vpcSet.UniqVpcId != nil {
						vpcSetMap["uniq_vpc_id"] = vpcSet.UniqVpcId
					}

					if vpcSet.Region != nil {
						vpcSetMap["region"] = vpcSet.Region
					}

					vpcSetList = append(vpcSetList, vpcSetMap)
				}

				forwardRuleSetMap["vpc_set"] = vpcSetList
			}

			if forwardRuleSet.ZoneId != nil {
				forwardRuleSetMap["zone_id"] = forwardRuleSet.ZoneId
			}

			tagsList := make([]map[string]interface{}, 0, len(forwardRuleSet.Tags))
			if forwardRuleSet.Tags != nil {
				for _, tags := range forwardRuleSet.Tags {
					tagsMap := map[string]interface{}{}
					if tags.TagKey != nil {
						tagsMap["tag_key"] = tags.TagKey
					}

					if tags.TagValue != nil {
						tagsMap["tag_value"] = tags.TagValue
					}

					tagsList = append(tagsList, tagsMap)
				}

				forwardRuleSetMap["tags"] = tagsList
			}

			forwardRuleSetList = append(forwardRuleSetList, forwardRuleSetMap)
		}

		_ = d.Set("forward_rule_set", forwardRuleSetList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), forwardRuleSetList); e != nil {
			return e
		}
	}

	return nil
}
