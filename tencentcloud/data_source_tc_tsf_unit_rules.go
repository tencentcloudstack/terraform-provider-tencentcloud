/*
Use this data source to query detailed information of tsf unit_rules

Example Usage

```hcl
data "tencentcloud_tsf_unit_rules" "unit_rules" {
  gateway_instance_id = "gw-ins-lvdypq5k"
  status = "disabled"
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	tsf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tsf/v20180326"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudTsfUnitRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTsfUnitRulesRead,
		Schema: map[string]*schema.Schema{
			"gateway_instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "gateway instance id.",
			},

			"status": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Enabled state, disabled: unpublished, enabled: published.",
			},

			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Pagination list information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "total number of records.",
						},
						"content": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "record entity list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "rule name.",
									},
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "rule ID.",
									},
									"gateway_instance_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Gateway Entity ID.",
									},
									"description": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Rule description.",
									},
									"status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Use status: enabled/disabled.",
									},
									"unit_rule_item_list": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "list of rule items.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"relationship": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Logical relationship: AND/OR.",
												},
												"dest_namespace_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Destination Namespace ID.",
												},
												"dest_namespace_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "destination namespace name.",
												},
												"name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "rule item name.",
												},
												"id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "rule item ID.",
												},
												"unit_rule_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Unitization rule ID.",
												},
												"priority": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Rule order, the smaller the higher the priority: the default is 0.",
												},
												"description": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Rule description.",
												},
												"unit_rule_tag_list": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "List of rule labels.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"tag_type": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Tag Type: U(User Tag).",
															},
															"tag_field": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "tag name.",
															},
															"tag_operator": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Operator: IN/NOT_IN/EQUAL/NOT_EQUAL/REGEX.",
															},
															"tag_value": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "tag value.",
															},
															"unit_rule_item_id": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Unitization rule item ID.",
															},
															"id": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "rule ID.",
															},
														},
													},
												},
											},
										},
									},
									"created_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "created time.",
									},
									"updated_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Updated time.",
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

func dataSourceTencentCloudTsfUnitRulesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_tsf_unit_rules.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("gateway_instance_id"); ok {
		paramMap["GatewayInstanceId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("status"); ok {
		paramMap["Status"] = helper.String(v.(string))
	}

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}

	var unitRule *tsf.TsfPageUnitRuleV2
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTsfUnitRulesByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		unitRule = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(unitRule.Content))
	tmpList := make([]map[string]interface{}, 0, len(unitRule.Content))
	if unitRule != nil {
		tsfPageUnitRuleMap := map[string]interface{}{}
		if unitRule.TotalCount != nil {
			tsfPageUnitRuleMap["total_count"] = unitRule.TotalCount
		}

		if unitRule.Content != nil {
			contentList := []interface{}{}
			for _, content := range unitRule.Content {
				contentMap := map[string]interface{}{}

				if content.Name != nil {
					contentMap["name"] = content.Name
				}

				if content.Id != nil {
					contentMap["id"] = content.Id
				}

				if content.GatewayInstanceId != nil {
					contentMap["gateway_instance_id"] = content.GatewayInstanceId
				}

				if content.Description != nil {
					contentMap["description"] = content.Description
				}

				if content.Status != nil {
					contentMap["status"] = content.Status
				}

				if content.UnitRuleItemList != nil {
					unitRuleItemListList := []interface{}{}
					for _, unitRuleItemList := range content.UnitRuleItemList {
						unitRuleItemListMap := map[string]interface{}{}

						if unitRuleItemList.Relationship != nil {
							unitRuleItemListMap["relationship"] = unitRuleItemList.Relationship
						}

						if unitRuleItemList.DestNamespaceId != nil {
							unitRuleItemListMap["dest_namespace_id"] = unitRuleItemList.DestNamespaceId
						}

						if unitRuleItemList.DestNamespaceName != nil {
							unitRuleItemListMap["dest_namespace_name"] = unitRuleItemList.DestNamespaceName
						}

						if unitRuleItemList.Name != nil {
							unitRuleItemListMap["name"] = unitRuleItemList.Name
						}

						if unitRuleItemList.Id != nil {
							unitRuleItemListMap["id"] = unitRuleItemList.Id
						}

						if unitRuleItemList.UnitRuleId != nil {
							unitRuleItemListMap["unit_rule_id"] = unitRuleItemList.UnitRuleId
						}

						if unitRuleItemList.Priority != nil {
							unitRuleItemListMap["priority"] = unitRuleItemList.Priority
						}

						if unitRuleItemList.Description != nil {
							unitRuleItemListMap["description"] = unitRuleItemList.Description
						}

						if unitRuleItemList.UnitRuleTagList != nil {
							unitRuleTagListList := []interface{}{}
							for _, unitRuleTagList := range unitRuleItemList.UnitRuleTagList {
								unitRuleTagListMap := map[string]interface{}{}

								if unitRuleTagList.TagType != nil {
									unitRuleTagListMap["tag_type"] = unitRuleTagList.TagType
								}

								if unitRuleTagList.TagField != nil {
									unitRuleTagListMap["tag_field"] = unitRuleTagList.TagField
								}

								if unitRuleTagList.TagOperator != nil {
									unitRuleTagListMap["tag_operator"] = unitRuleTagList.TagOperator
								}

								if unitRuleTagList.TagValue != nil {
									unitRuleTagListMap["tag_value"] = unitRuleTagList.TagValue
								}

								if unitRuleTagList.UnitRuleItemId != nil {
									unitRuleTagListMap["unit_rule_item_id"] = unitRuleTagList.UnitRuleItemId
								}

								if unitRuleTagList.Id != nil {
									unitRuleTagListMap["id"] = unitRuleTagList.Id
								}

								unitRuleTagListList = append(unitRuleTagListList, unitRuleTagListMap)
							}

							unitRuleItemListMap["unit_rule_tag_list"] = unitRuleTagListList
						}

						unitRuleItemListList = append(unitRuleItemListList, unitRuleItemListMap)
					}

					contentMap["unit_rule_item_list"] = unitRuleItemListList
				}

				if content.CreatedTime != nil {
					contentMap["created_time"] = content.CreatedTime
				}

				if content.UpdatedTime != nil {
					contentMap["updated_time"] = content.UpdatedTime
				}

				contentList = append(contentList, contentMap)
				ids = append(ids, *content.GatewayInstanceId)
			}

			tsfPageUnitRuleMap["content"] = contentList
		}

		_ = d.Set("result", []interface{}{tsfPageUnitRuleMap})
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
