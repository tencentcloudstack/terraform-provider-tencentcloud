/*
Use this data source to query detailed information of teo ruleEngineSettings

Example Usage

```hcl
data "tencentcloud_teo_rule_engine_settings" "ruleEngineSettings" {
}
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
)

func dataSourceTencentCloudTeoRuleEngineSettings() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTeoRuleEngineSettingsRead,
		Schema: map[string]*schema.Schema{
			"actions": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Detail info of actions which can be used in rule engine.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Action name.",
						},
						"properties": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Action properties.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Property name.",
									},
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Property value type. Valid values:- `CHOICE`: enum type, must select one of the value in `ChoicesValue`.- `TOGGLE`: switch type, must select one of the value in `ChoicesValue`.- `OBJECT`: object type, the `ChoiceProperties` list all properties of the object.- `CUSTOM_NUM`: integer type.- `CUSTOM_STRING`: string type.",
									},
									"choices_value": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "The choices which can be used. This list may be empty.",
									},
									"min": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Min integer value can bse used when property type is `CUSTOM_NUM`. When `Min` and `Max` both are 0, this field is meaningless.",
									},
									"max": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Max integer value can bse used when property type is `CUSTOM_NUM`. When `Min` and `Max` both are 0, this field is meaningless.",
									},
									"is_multiple": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Whether this property is allowed to set multiple values.",
									},
									"is_allow_empty": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Whether this property is allowed to set empty.",
									},
									"choice_properties": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Associative properties of this property, they are all required. Note: This field may return null, indicating that no valid value can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Property name.",
												},
												"type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Property value type. Valid values:- `CHOICE`: enum type, must select one of the value in `ChoicesValue`.- `TOGGLE`: switch type, must select one of the value in `ChoicesValue`.- `CUSTOM_NUM`: integer type.- `CUSTOM_STRING`: string type.",
												},
												"choices_value": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Computed:    true,
													Description: "The choices which can bse used. This list may be empty.",
												},
												"min": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Min integer value can bse used when property type is `CUSTOM_NUM`. When `Min` and `Max` both are 0, this field is meaningless.",
												},
												"max": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Max integer value can bse used when property type is `CUSTOM_NUM`. When `Min` and `Max` both are 0, this field is meaningless.",
												},
												"is_multiple": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Whether this property is allowed to set multiple values.",
												},
												"is_allow_empty": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Whether this property is allowed to set empty.",
												},
												"extra_parameter": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Special parameter. Note: This field may return null, indicating that no valid value can be obtained.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"id": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Parameter name. Valid values:- `Action`: this extra parameter is required when modify HTTP header, this action should be a `RewriteAction`.- `StatusCode`: this extra parameter is required when modify HTTP status code, this action should be a `CodeAction`.- `NULL`: this action should be a `NormalAction`.",
															},
															"type": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Parameter value type. Valid values:- `CHOICE`: select one value from `Choices`.- `CUSTOM_NUM`: integer value.- `CUSTOM_STRING`: string value.",
															},
															"choices": {
																Type: schema.TypeSet,
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
																Computed:    true,
																Description: "A list of choices which can be used when `Type` is `CHOICE`.",
															},
														},
													},
												},
											},
										},
									},
									"extra_parameter": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Special parameter. Note: This field may return null, indicating that no valid value can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Parameter name. Valid values:- `Action`: this extra parameter is required when modify HTTP header, this action should be a `RewriteAction`.- `StatusCode`: this extra parameter is required when modify HTTP status code, this action should be a `CodeAction`.- `NULL`: this action should be a `NormalAction`.",
												},
												"type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Parameter value type. Valid values:- `CHOICE`: select one value from `Choices`.- `CUSTOM_NUM`: integer value.- `CUSTOM_STRING`: string value.",
												},
												"choices": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Computed:    true,
													Description: "A list of choices which can be used when `Type` is `CHOICE`.",
												},
											},
										},
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

func dataSourceTencentCloudTeoRuleEngineSettingsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_teo_rule_engine_settings.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	teoService := TeoService{client: meta.(*TencentCloudClient).apiV3Conn}

	var rules []*teo.RulesSettingAction
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		results, e := teoService.DescribeTeoRuleEngineSettingsByFilter(ctx)
		if e != nil {
			return retryError(e)
		}
		rules = results
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read Teo rules failed, reason:%+v", logId, err)
		return err
	}

	ruleList := []interface{}{}
	if rules != nil {
		for _, rule := range rules {
			ruleMap := map[string]interface{}{}
			if rule.Action != nil {
				ruleMap["action"] = rule.Action
			}
			if rule.Properties != nil {
				propertiesList := []interface{}{}
				for _, properties := range rule.Properties {
					propertiesMap := map[string]interface{}{}
					if properties.Name != nil {
						propertiesMap["name"] = properties.Name
					}
					if properties.Type != nil {
						propertiesMap["type"] = properties.Type
					}
					if properties.ChoicesValue != nil {
						propertiesMap["choices_value"] = properties.ChoicesValue
					}
					if properties.Min != nil {
						propertiesMap["min"] = properties.Min
					}
					if properties.Max != nil {
						propertiesMap["max"] = properties.Max
					}
					if properties.IsMultiple != nil {
						propertiesMap["is_multiple"] = properties.IsMultiple
					}
					if properties.IsAllowEmpty != nil {
						propertiesMap["is_allow_empty"] = properties.IsAllowEmpty
					}
					if properties.ChoiceProperties != nil {
						choicePropertiesList := []interface{}{}
						for _, choiceProperties := range properties.ChoiceProperties {
							choicePropertiesMap := map[string]interface{}{}
							if choiceProperties.Name != nil {
								choicePropertiesMap["name"] = choiceProperties.Name
							}
							if choiceProperties.Type != nil {
								choicePropertiesMap["type"] = choiceProperties.Type
							}
							if choiceProperties.ChoicesValue != nil {
								choicePropertiesMap["choices_value"] = choiceProperties.ChoicesValue
							}
							if choiceProperties.Min != nil {
								choicePropertiesMap["min"] = choiceProperties.Min
							}
							if choiceProperties.Max != nil {
								choicePropertiesMap["max"] = choiceProperties.Max
							}
							if choiceProperties.IsMultiple != nil {
								choicePropertiesMap["is_multiple"] = choiceProperties.IsMultiple
							}
							if choiceProperties.IsAllowEmpty != nil {
								choicePropertiesMap["is_allow_empty"] = choiceProperties.IsAllowEmpty
							}
							if choiceProperties.ExtraParameter != nil {
								extraParameterMap := map[string]interface{}{}
								if choiceProperties.ExtraParameter.Id != nil {
									extraParameterMap["id"] = choiceProperties.ExtraParameter.Id
								}
								if choiceProperties.ExtraParameter.Type != nil {
									extraParameterMap["type"] = choiceProperties.ExtraParameter.Type
								}
								if choiceProperties.ExtraParameter.Choices != nil {
									extraParameterMap["choices"] = choiceProperties.ExtraParameter.Choices
								}

								choicePropertiesMap["extra_parameter"] = []interface{}{extraParameterMap}
							}

							choicePropertiesList = append(choicePropertiesList, choicePropertiesMap)
						}
						propertiesMap["choice_properties"] = choicePropertiesList
					}
					if properties.ExtraParameter != nil {
						extraParameterMap := map[string]interface{}{}
						if properties.ExtraParameter.Id != nil {
							extraParameterMap["id"] = properties.ExtraParameter.Id
						}
						if properties.ExtraParameter.Type != nil {
							extraParameterMap["type"] = properties.ExtraParameter.Type
						}
						if properties.ExtraParameter.Choices != nil {
							extraParameterMap["choices"] = properties.ExtraParameter.Choices
						}

						propertiesMap["extra_parameter"] = []interface{}{extraParameterMap}
					}

					propertiesList = append(propertiesList, propertiesMap)
				}
				ruleMap["properties"] = propertiesList
			}

			ruleList = append(ruleList, ruleMap)
		}
		_ = d.Set("actions", ruleList)
	}

	d.SetId("rule_engine_settings")

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), ruleList); e != nil {
			return e
		}
	}
	return nil
}
