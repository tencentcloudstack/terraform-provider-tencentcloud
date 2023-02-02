/*
Provides a resource to create a tsf unit_rule

Example Usage

```hcl
resource "tencentcloud_tsf_unit_rule" "unit_rule" {
  gateway_instance_id = ""
  name = ""
  description = ""
  unit_rule_item_list {
		relationship = ""
		dest_namespace_id = ""
		dest_namespace_name = ""
		name = ""
		rule_id = ""
		unit_rule_id = ""
		priority =
		description = ""
		unit_rule_tag_list {
			tag_type = ""
			tag_field = ""
			tag_operator = ""
			tag_value = ""
			unit_rule_item_id = ""
			rule_id = ""
		}

  }
}
```

Import

tsf unit_rule can be imported using the id, e.g.

```
terraform import tencentcloud_tsf_unit_rule.unit_rule unit_rule_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	tsf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tsf/v20180326"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudTsfUnitRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTsfUnitRuleCreate,
		Read:   resourceTencentCloudTsfUnitRuleRead,
		Update: resourceTencentCloudTsfUnitRuleUpdate,
		Delete: resourceTencentCloudTsfUnitRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"gateway_instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "gateway entity ID.",
			},

			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "rule name.",
			},

			"rule_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "rule ID.",
			},

			"status": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "usage status: enabled/disabled.",
			},

			"description": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "rule description.",
			},

			"unit_rule_item_list": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "list of rule items.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"relationship": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "logical relationship: AND/OR.",
						},
						"dest_namespace_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "destination namespace ID.",
						},
						"dest_namespace_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "destination namespace name.",
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "rule item name.",
						},
						"rule_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "rule item ID.",
						},
						"unit_rule_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Unitization rule ID.",
						},
						"priority": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "rule order, the smaller the higher the priority: the default is 0.",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "rule description.",
						},
						"unit_rule_tag_list": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "list of rule labels.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"tag_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Tag Type: U(User Tag).",
									},
									"tag_field": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "label name.",
									},
									"tag_operator": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Operator: IN/NOT_IN/EQUAL/NOT_EQUAL/REGEX.",
									},
									"tag_value": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "tag value.",
									},
									"unit_rule_item_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Unitization rule item ID.",
									},
									"rule_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "rule ID.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudTsfUnitRuleCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_unit_rule.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request = tsf.NewCreateUnitRuleRequest()
		// response = tsf.NewCreateUnitRuleResponse()
		ruleId string
	)
	if v, ok := d.GetOk("gateway_instance_id"); ok {
		request.GatewayInstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOk("unit_rule_item_list"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			unitRuleItem := tsf.UnitRuleItem{}
			if v, ok := dMap["relationship"]; ok {
				unitRuleItem.Relationship = helper.String(v.(string))
			}
			if v, ok := dMap["dest_namespace_id"]; ok {
				unitRuleItem.DestNamespaceId = helper.String(v.(string))
			}
			if v, ok := dMap["dest_namespace_name"]; ok {
				unitRuleItem.DestNamespaceName = helper.String(v.(string))
			}
			if v, ok := dMap["name"]; ok {
				unitRuleItem.Name = helper.String(v.(string))
			}
			if v, ok := dMap["rule_id"]; ok {
				unitRuleItem.Id = helper.String(v.(string))
			}
			if v, ok := dMap["unit_rule_id"]; ok {
				unitRuleItem.UnitRuleId = helper.String(v.(string))
			}
			if v, ok := dMap["priority"]; ok {
				unitRuleItem.Priority = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["description"]; ok {
				unitRuleItem.Description = helper.String(v.(string))
			}
			if v, ok := dMap["unit_rule_tag_list"]; ok {
				for _, item := range v.([]interface{}) {
					unitRuleTagListMap := item.(map[string]interface{})
					unitRuleTag := tsf.UnitRuleTag{}
					if v, ok := unitRuleTagListMap["tag_type"]; ok {
						unitRuleTag.TagType = helper.String(v.(string))
					}
					if v, ok := unitRuleTagListMap["tag_field"]; ok {
						unitRuleTag.TagField = helper.String(v.(string))
					}
					if v, ok := unitRuleTagListMap["tag_operator"]; ok {
						unitRuleTag.TagOperator = helper.String(v.(string))
					}
					if v, ok := unitRuleTagListMap["tag_value"]; ok {
						unitRuleTag.TagValue = helper.String(v.(string))
					}
					if v, ok := unitRuleTagListMap["unit_rule_item_id"]; ok {
						unitRuleTag.UnitRuleItemId = helper.String(v.(string))
					}
					if v, ok := unitRuleTagListMap["rule_id"]; ok {
						unitRuleTag.Id = helper.String(v.(string))
					}
					unitRuleItem.UnitRuleTagList = append(unitRuleItem.UnitRuleTagList, &unitRuleTag)
				}
			}
			request.UnitRuleItemList = append(request.UnitRuleItemList, &unitRuleItem)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTsfClient().CreateUnitRule(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		// response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tsf unitRule failed, reason:%+v", logId, err)
		return err
	}

	// ruleId = *response.Response.RuleId
	d.SetId(ruleId)

	return resourceTencentCloudTsfUnitRuleRead(d, meta)
}

func resourceTencentCloudTsfUnitRuleRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_unit_rule.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}

	id := d.Id()

	unitRule, err := service.DescribeTsfUnitRuleById(ctx, id)
	if err != nil {
		return err
	}

	if unitRule == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TsfUnitRule` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if unitRule.GatewayInstanceId != nil {
		_ = d.Set("gateway_instance_id", unitRule.GatewayInstanceId)
	}

	if unitRule.Name != nil {
		_ = d.Set("name", unitRule.Name)
	}

	if unitRule.Id != nil {
		_ = d.Set("rule_id", unitRule.Id)
	}

	if unitRule.Status != nil {
		_ = d.Set("status", unitRule.Status)
	}

	if unitRule.Description != nil {
		_ = d.Set("description", unitRule.Description)
	}

	if unitRule.UnitRuleItemList != nil {
		unitRuleItemListList := []interface{}{}
		for _, unitRuleItemList := range unitRule.UnitRuleItemList {
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
				unitRuleItemListMap["rule_id"] = unitRuleItemList.Id
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
						unitRuleTagListMap["rule_id"] = unitRuleTagList.Id
					}

					unitRuleTagListList = append(unitRuleTagListList, unitRuleTagListMap)
				}

				unitRuleItemListMap["unit_rule_tag_list"] = []interface{}{unitRuleTagListList}
			}

			unitRuleItemListList = append(unitRuleItemListList, unitRuleItemListMap)
		}

		_ = d.Set("unit_rule_item_list", unitRuleItemListList)

	}

	return nil
}

func resourceTencentCloudTsfUnitRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_unit_rule.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := tsf.NewUpdateUnitRuleRequest()

	id := d.Id()

	request.Id = &id

	immutableArgs := []string{"gateway_instance_id", "name", "rule_id", "status", "description", "unit_rule_item_list"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("name") {
		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}
	}

	if d.HasChange("description") {
		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}
	}

	if d.HasChange("unit_rule_item_list") {
		if v, ok := d.GetOk("unit_rule_item_list"); ok {
			for _, item := range v.([]interface{}) {
				dMap := item.(map[string]interface{})
				unitRuleItem := tsf.UnitRuleItem{}
				if v, ok := dMap["relationship"]; ok {
					unitRuleItem.Relationship = helper.String(v.(string))
				}
				if v, ok := dMap["dest_namespace_id"]; ok {
					unitRuleItem.DestNamespaceId = helper.String(v.(string))
				}
				if v, ok := dMap["dest_namespace_name"]; ok {
					unitRuleItem.DestNamespaceName = helper.String(v.(string))
				}
				if v, ok := dMap["name"]; ok {
					unitRuleItem.Name = helper.String(v.(string))
				}
				if v, ok := dMap["rule_id"]; ok {
					unitRuleItem.Id = helper.String(v.(string))
				}
				if v, ok := dMap["unit_rule_id"]; ok {
					unitRuleItem.UnitRuleId = helper.String(v.(string))
				}
				if v, ok := dMap["priority"]; ok {
					unitRuleItem.Priority = helper.IntInt64(v.(int))
				}
				if v, ok := dMap["description"]; ok {
					unitRuleItem.Description = helper.String(v.(string))
				}
				if v, ok := dMap["unit_rule_tag_list"]; ok {
					for _, item := range v.([]interface{}) {
						unitRuleTagListMap := item.(map[string]interface{})
						unitRuleTag := tsf.UnitRuleTag{}
						if v, ok := unitRuleTagListMap["tag_type"]; ok {
							unitRuleTag.TagType = helper.String(v.(string))
						}
						if v, ok := unitRuleTagListMap["tag_field"]; ok {
							unitRuleTag.TagField = helper.String(v.(string))
						}
						if v, ok := unitRuleTagListMap["tag_operator"]; ok {
							unitRuleTag.TagOperator = helper.String(v.(string))
						}
						if v, ok := unitRuleTagListMap["tag_value"]; ok {
							unitRuleTag.TagValue = helper.String(v.(string))
						}
						if v, ok := unitRuleTagListMap["unit_rule_item_id"]; ok {
							unitRuleTag.UnitRuleItemId = helper.String(v.(string))
						}
						if v, ok := unitRuleTagListMap["rule_id"]; ok {
							unitRuleTag.Id = helper.String(v.(string))
						}
						unitRuleItem.UnitRuleTagList = append(unitRuleItem.UnitRuleTagList, &unitRuleTag)
					}
				}
				request.UnitRuleItemList = append(request.UnitRuleItemList, &unitRuleItem)
			}
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTsfClient().UpdateUnitRule(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update tsf unitRule failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTsfUnitRuleRead(d, meta)
}

func resourceTencentCloudTsfUnitRuleDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_unit_rule.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}
	id := d.Id()

	if err := service.DeleteTsfUnitRuleById(ctx, id); err != nil {
		return err
	}

	return nil
}
