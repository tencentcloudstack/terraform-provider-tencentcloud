/*
Use this data source to query detailed information of tsf unit_rules

Example Usage

```hcl
data "tencentcloud_tsf_unit_rules" "unit_rules" {
  gateway_instance_id = "gw-ins-lvdypq5k"
  search_word = "test"
  status = ""
  }
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
				Description: "Gateway instance id.",
			},

			"search_word": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Fuzzy search based on rule name or comment content.",
			},

			"status": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Enabled status, disabled: unpublished, enabled: published.",
			},

			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Pagination list information.Note: This field may return null, indicating that no valid values can be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Record count.",
						},
						"content": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Instance list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Rule name.",
									},
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Rule Id .Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"gateway_instance_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Gateway Instance Id.Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"description": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Description.Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Enabled status .Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"unit_rule_item_list": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Rule Id list .Note: This field may return null, indicating that no valid values can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"relationship": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Logic relationship between rule：AND/OR.",
												},
												"dest_namespace_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Destination namespaceId.",
												},
												"dest_namespace_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Destination namespace name.",
												},
												"name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Rule item name.",
												},
												"id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Rule Id .Note: This field may return null, indicating that no valid values can be obtained.",
												},
												"unit_rule_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Unit rule Id list .Note: This field may return null, indicating that no valid values can be obtained.",
												},
												"priority": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Rule order, the smaller the priority, the higher the priority: default is 0.Note: This field may return null, indicating that no valid values can be obtained.",
												},
												"description": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Description.Note: This field may return null, indicating that no valid values can be obtained.",
												},
												"unit_rule_tag_list": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Unit rule tag list.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"tag_type": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Tag type : U(user tag).",
															},
															"tag_field": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Tag name.",
															},
															"tag_operator": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Operator for tag:IN/NOT_IN/EQUAL/NOT_EQUAL/REGEX.",
															},
															"tag_value": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Tag value.",
															},
															"unit_rule_item_id": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Unit rule id.Note: This field may return null, indicating that no valid values can be obtained.",
															},
															"id": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Rule id.Note: This field may return null, indicating that no valid values can be obtained.",
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
										Description: "CreatedTime.Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"updated_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "UpdatedTime.Note: This field may return null, indicating that no valid values can be obtained.",
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

	if v, ok := d.GetOk("search_word"); ok {
		paramMap["SearchWord"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("status"); ok {
		paramMap["Status"] = helper.String(v.(string))
	}

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}

	var result []*tsf.TsfPageUnitRuleV2

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTsfUnitRulesByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		result = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(result))
	if result != nil {
		tsfPageUnitRuleMap := map[string]interface{}{}

		if result.TotalCount != nil {
			tsfPageUnitRuleMap["total_count"] = result.TotalCount
		}

		if result.Content != nil {
			contentList := []interface{}{}
			for _, content := range result.Content {
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

							unitRuleItemListMap["unit_rule_tag_list"] = []interface{}{unitRuleTagListList}
						}

						unitRuleItemListList = append(unitRuleItemListList, unitRuleItemListMap)
					}

					contentMap["unit_rule_item_list"] = []interface{}{unitRuleItemListList}
				}

				if content.CreatedTime != nil {
					contentMap["created_time"] = content.CreatedTime
				}

				if content.UpdatedTime != nil {
					contentMap["updated_time"] = content.UpdatedTime
				}

				contentList = append(contentList, contentMap)
			}

			tsfPageUnitRuleMap["content"] = []interface{}{contentList}
		}

		ids = append(ids, *result.GatewayInstanceId)
		_ = d.Set("result", tsfPageUnitRuleMap)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tsfPageUnitRuleV2Map); e != nil {
			return e
		}
	}
	return nil
}
