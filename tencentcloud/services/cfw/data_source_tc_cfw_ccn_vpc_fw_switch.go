package cfw

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cfwv20190904 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cfw/v20190904"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCfwCcnVpcFwSwitch() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCfwCcnVpcFwSwitchRead,
		Schema: map[string]*schema.Schema{
			"ccn_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "CCN ID.",
			},

			"interconnect_pairs": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Interconnect pair configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_a": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Group A.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"instance_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance ID.",
									},
									"instance_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance type such as VPC or DIRECTCONNECT.",
									},
									"instance_region": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Region where the instance is located.",
									},
									"access_cidr_mode": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Network segment mode for accessing firewall: 0-no access, 1-access all network segments associated with the instance, 2-access user-defined network segments.",
									},
									"access_cidr_list": {
										Type:        schema.TypeSet,
										Computed:    true,
										Description: "List of network segments for accessing firewall.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"group_b": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Group B.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"instance_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance ID.",
									},
									"instance_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance type such as VPC or DIRECTCONNECT.",
									},
									"instance_region": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Region where the instance is located.",
									},
									"access_cidr_mode": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Network segment mode for accessing firewall: 0-no access, 1-access all network segments associated with the instance, 2-access user-defined network segments.",
									},
									"access_cidr_list": {
										Type:        schema.TypeSet,
										Computed:    true,
										Description: "List of network segments for accessing firewall.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"interconnect_mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Interconnect mode: \"CrossConnect\": cross interconnect (each instance in group A interconnects with each instance in group B), \"FullMesh\": full mesh (group A content is identical to group B, equivalent to pairwise interconnection within the group).",
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

func dataSourceTencentCloudCfwCcnVpcFwSwitchRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_cfw_ccn_vpc_fw_switch.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = CfwService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		ccnId   string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("ccn_id"); ok {
		paramMap["CcnId"] = helper.String(v.(string))
		ccnId = v.(string)
	}

	var respData []*cfwv20190904.InterconnectPair
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCfwCcnVpcFwSwitchByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		respData = result
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	interconnectPairsList := make([]map[string]interface{}, 0, len(respData))
	if respData != nil {
		for _, interconnectPairs := range respData {
			interconnectPairsMap := map[string]interface{}{}
			groupAList := make([]map[string]interface{}, 0, len(interconnectPairs.GroupA))
			if interconnectPairs.GroupA != nil {
				for _, groupA := range interconnectPairs.GroupA {
					groupAMap := map[string]interface{}{}
					if groupA.InstanceId != nil {
						groupAMap["instance_id"] = groupA.InstanceId
					}

					if groupA.InstanceType != nil {
						groupAMap["instance_type"] = groupA.InstanceType
					}

					if groupA.InstanceRegion != nil {
						groupAMap["instance_region"] = groupA.InstanceRegion
					}

					if groupA.AccessCidrMode != nil {
						groupAMap["access_cidr_mode"] = groupA.AccessCidrMode
					}

					if groupA.AccessCidrList != nil {
						groupAMap["access_cidr_list"] = groupA.AccessCidrList
					}

					groupAList = append(groupAList, groupAMap)
				}

				interconnectPairsMap["group_a"] = groupAList
			}

			groupBList := make([]map[string]interface{}, 0, len(interconnectPairs.GroupB))
			if interconnectPairs.GroupB != nil {
				for _, groupB := range interconnectPairs.GroupB {
					groupBMap := map[string]interface{}{}
					if groupB.InstanceId != nil {
						groupBMap["instance_id"] = groupB.InstanceId
					}

					if groupB.InstanceType != nil {
						groupBMap["instance_type"] = groupB.InstanceType
					}

					if groupB.InstanceRegion != nil {
						groupBMap["instance_region"] = groupB.InstanceRegion
					}

					if groupB.AccessCidrMode != nil {
						groupBMap["access_cidr_mode"] = groupB.AccessCidrMode
					}

					if groupB.AccessCidrList != nil {
						groupBMap["access_cidr_list"] = groupB.AccessCidrList
					}

					groupBList = append(groupBList, groupBMap)
				}

				interconnectPairsMap["group_b"] = groupBList
			}

			if interconnectPairs.InterconnectMode != nil {
				interconnectPairsMap["interconnect_mode"] = interconnectPairs.InterconnectMode
			}

			interconnectPairsList = append(interconnectPairsList, interconnectPairsMap)
		}

		_ = d.Set("interconnect_pairs", interconnectPairsList)
	}

	d.SetId(ccnId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), interconnectPairsList); e != nil {
			return e
		}
	}

	return nil
}
